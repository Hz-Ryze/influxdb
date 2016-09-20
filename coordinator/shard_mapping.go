package coordinator

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/influxdata/influxdb/influxql"
	"github.com/influxdata/influxdb/tsdb"
)

type Source struct {
	Database        string
	RetentionPolicy string
}

type ShardInfo struct {
	*tsdb.Shard
	Measurements []string
	StartTime    time.Time
}

type ShardMapper map[Source][]*ShardInfo

func (a ShardMapper) FieldDimensions() (fields map[string]influxql.DataType, dimensions map[string]struct{}, err error) {
	fields = make(map[string]influxql.DataType)
	dimensions = make(map[string]struct{})

	for _, shards := range a {
		for _, sh := range shards {
			f, d, err := sh.FieldDimensions(sh.Measurements)
			if err != nil {
				return nil, nil, err
			}

			for k, typ := range f {
				if _, ok := fields[k]; typ != influxql.Unknown && (!ok || typ < fields[k]) {
					fields[k] = typ
				}
			}
			for k := range d {
				dimensions[k] = struct{}{}
			}
		}
	}
	return
}

func (a ShardMapper) CreateIterator(opt influxql.IteratorOptions) (influxql.Iterator, error) {
	if influxql.Sources(opt.Sources).HasSystemSource() {
		var shards []influxql.IteratorCreator
		for _, group := range a {
			for _, sh := range group {
				shards = append(shards, sh)
			}
		}
		return influxql.NewLazyIterator(shards, opt)
	}

	var mu sync.Mutex
	parallelism := runtime.GOMAXPROCS(0)
	ch := make(chan struct{}, parallelism)
	for i := 0; i < parallelism; i++ {
		ch <- struct{}{}
	}

	var wg sync.WaitGroup
	seriesMap := make(map[string]map[Source]map[string][]*SeriesInfo)
	for s, shards := range a {
		for _, sh := range shards {
			wg.Add(1)
			<-ch
			go func(source Source, sh *ShardInfo) {
				defer wg.Done()
				defer func() {
					ch <- struct{}{}
				}()

				for _, name := range sh.Measurements {
					series, err := createTagSets(sh, name, &opt)
					if series == nil || err != nil {
						return
					}

					mu.Lock()
					// Retrieve the series mapping for the measurement.
					m, ok := seriesMap[name]
					if !ok {
						m = make(map[Source]map[string][]*SeriesInfo)
						seriesMap[name] = m
					}

					// Retrieve the series mapping for the current source.
					ms, ok := m[source]
					if !ok {
						ms = make(map[string][]*SeriesInfo)
						m[source] = ms
					}

					// Insert each series into map in the appropriate location.
					for _, s := range series {
						ms[string(s.SeriesKey)] = append(ms[string(s.SeriesKey)], s)
					}
					mu.Unlock()
				}
			}(s, sh)
		}
	}
	wg.Wait()

	// Retrieve the measurements and order them.
	names := make([]string, 0, len(seriesMap))
	for name := range seriesMap {
		names = append(names, name)
	}
	sort.Strings(names)

	// Create lazy iterator for each measurement in order.
	ics := make([]influxql.IteratorCreator, len(names))
	for i, name := range names {
		sources := seriesMap[name] // map[Source]map[string][]interface{}

		// Merge the iterators from each database source.
		outer := make([]influxql.IteratorCreator, 0, len(sources))
		for _, s := range sources {
			keys := make([]string, 0, len(s))
			for k := range s {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			// Iterate through each series key.
			inner := make([]influxql.IteratorCreator, len(keys))
			for i, key := range keys {
				// For each key, sort the series info by time and create a lazy iterator.
				shards := s[key]
				sort.Sort(byTime(shards))

				shardics := make([]influxql.IteratorCreator, len(shards))
				for i, shard := range shards {
					shardics[i] = shard
				}
				inner[i] = influxql.LazyIteratorCreator(shardics)
			}
			outer = append(outer, influxql.LazyIteratorCreator(inner))
		}
		ics[i] = influxql.IteratorCreators(outer)
	}
	return influxql.NewLazyIterator(ics, opt)
}

type SeriesInfo struct {
	*tsdb.Shard
	Measurement *tsdb.Measurement
	SeriesKey   []byte
	TagSet      *influxql.TagSet
	StartTime   time.Time
}

func (si *SeriesInfo) CreateIterator(opt influxql.IteratorOptions) (influxql.Iterator, error) {
	return si.Shard.CreateSeriesIterator(si.Measurement, si.TagSet, opt)
}

func (si *SeriesInfo) String() string {
	return fmt.Sprintf("{measurement: %s, series key: %s, tag sets: %d}", si.Measurement.Name, string(si.SeriesKey), len(si.TagSet.SeriesKeys))
}

type byTime []*SeriesInfo

func (a byTime) Len() int { return len(a) }
func (a byTime) Less(i, j int) bool {
	return a[i].StartTime.Before(a[j].StartTime)
}
func (a byTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// marshalTags converts a tag mapping to a byte string. This is different from the normal
// series key because it contains empty tags too.
func marshalTags(tags map[string]string) []byte {
	// Empty maps marshal to empty bytes.
	if len(tags) == 0 {
		return nil
	}

	// Extract keys and determine final size.
	sz := (len(tags) * 2) - 1 // separators
	keys := make([]string, 0, len(tags))
	for k, v := range tags {
		keys = append(keys, k)
		sz += len(k) + len(v)
	}
	sort.Strings(keys)

	// Generate marshaled bytes.
	b := make([]byte, sz)
	buf := b
	for i, k := range keys {
		copy(buf, k)
		buf[len(k)] = '='
		buf = buf[len(k)+1:]

		v := tags[k]
		copy(buf, v)
		if i < len(keys)-1 {
			buf[len(v)] = ','
			buf = buf[len(v)+1:]
		}
	}
	return b
}

func createTagSets(si *ShardInfo, name string, opt *influxql.IteratorOptions) ([]*SeriesInfo, error) {
	mm := si.Shard.MeasurementByName(name)
	if mm == nil {
		return nil, nil
	}

	// Determine tagsets for this measurement based on dimensions and filters.
	tagSets, err := mm.TagSets(opt.Dimensions, opt.Condition)
	if err != nil {
		return nil, err
	}

	// Calculate tag sets and apply SLIMIT/SOFFSET.
	tagSets = influxql.LimitTagSets(tagSets, opt.SLimit, opt.SOffset)

	if len(tagSets) == 0 {
		return nil, nil
	}

	series := make([]*SeriesInfo, 0, len(tagSets))
	for _, t := range tagSets {
		series = append(series, &SeriesInfo{
			Shard:       si.Shard,
			Measurement: mm,
			SeriesKey:   marshalTags(t.Tags),
			TagSet:      t,
			StartTime:   si.StartTime,
		})
	}
	return series, nil
}
