/*

Series:

	╔══════Series List═════╗
	║ ┌───────────────────┐║
	║ │     Term List     │║
	║ ├───────────────────┤║
	║ │    Series Data    │║
	║ ├───────────────────┤║
	║ │      Trailer      │║
	║ └───────────────────┘║
	╚══════════════════════╝

	╔══════════Term List═══════════╗
	║ ┌──────────────────────────┐ ║
	║ │   Term Count <uint32>    │ ║
	║ └──────────────────────────┘ ║
	║ ┏━━━━━━━━━━━━━━━━━━━━━━━━━━┓ ║
	║ ┃ ┌──────────────────────┐ ┃ ║
	║ ┃ │  len(Term) <varint>  │ ┃ ║
	║ ┃ ├──────────────────────┤ ┃ ║
	║ ┃ │    Term <byte...>    │ ┃ ║
	║ ┃ └──────────────────────┘ ┃ ║
	║ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━┛ ║
	║ ┏━━━━━━━━━━━━━━━━━━━━━━━━━━┓ ║
	║ ┃ ┌──────────────────────┐ ┃ ║
	║ ┃ │  len(Term) <varint>  │ ┃ ║
	║ ┃ ├──────────────────────┤ ┃ ║
	║ ┃ │    Term <byte...>    │ ┃ ║
	║ ┃ └──────────────────────┘ ┃ ║
	║ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━┛ ║
	╚══════════════════════════════╝

	╔═════════Series Data══════════╗
	║ ┌──────────────────────────┐ ║
	║ │  Series Count <uint32>   │ ║
	║ └──────────────────────────┘ ║
	║ ┏━━━━━━━━━━━━━━━━━━━━━━━━━━┓ ║
	║ ┃ ┌──────────────────────┐ ┃ ║
	║ ┃ │     Flag <uint8>     │ ┃ ║
	║ ┃ ├──────────────────────┤ ┃ ║
	║ ┃ │ len(Series) <varint> │ ┃ ║
	║ ┃ ├──────────────────────┤ ┃ ║
	║ ┃ │   Series <byte...>   │ ┃ ║
	║ ┃ └──────────────────────┘ ┃ ║
	║ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━┛ ║
	║             ...              ║
	╚══════════════════════════════╝

	╔════════════Trailer══════════════╗
	║ ┌─────────────────────────────┐ ║
	║ │  Term List Offset <uint32>  │ ║
	║ ├─────────────────────────────┤ ║
	║ │  Series Data Pos <uint32>   │ ║
	║ └─────────────────────────────┘ ║
	╚═════════════════════════════════╝


Tag Sets:

	╔════════Tag Set═════════╗
	║┌──────────────────────┐║
	║│   Tag Values Block   │║
	║├──────────────────────┤║
	║│         ...          │║
	║├──────────────────────┤║
	║│    Tag Keys Block    │║
	║├──────────────────────┤║
	║│       Trailer        │║
	║└──────────────────────┘║
	╚════════════════════════╝

	╔═══════Tag Values Block═══════╗
	║                              ║
	║ ┏━━━━━━━━Value List━━━━━━━━┓ ║
	║ ┃                          ┃ ║
	║ ┃┏━━━━━━━━━Value━━━━━━━━━━┓┃ ║
	║ ┃┃┌──────────────────────┐┃┃ ║
	║ ┃┃│     Flag <uint8>     │┃┃ ║
	║ ┃┃├──────────────────────┤┃┃ ║
	║ ┃┃│ len(Value) <varint>  │┃┃ ║
	║ ┃┃├──────────────────────┤┃┃ ║
	║ ┃┃│   Value <byte...>    │┃┃ ║
	║ ┃┃├──────────────────────┤┃┃ ║
	║ ┃┃│ len(Series) <varint> │┃┃ ║
	║ ┃┃├──────────────────────┤┃┃ ║
	║ ┃┃│SeriesIDs <uint32...> │┃┃ ║
	║ ┃┃└──────────────────────┘┃┃ ║
	║ ┃┗━━━━━━━━━━━━━━━━━━━━━━━━┛┃ ║
	║ ┃           ...            ┃ ║
	║ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━┛ ║
	║ ┏━━━━━━━━Hash Index━━━━━━━━┓ ║
	║ ┃ ┌──────────────────────┐ ┃ ║
	║ ┃ │ len(Values) <uint32> │ ┃ ║
	║ ┃ ├──────────────────────┤ ┃ ║
	║ ┃ │Value Offset <uint64> │ ┃ ║
	║ ┃ ├──────────────────────┤ ┃ ║
	║ ┃ │         ...          │ ┃ ║
	║ ┃ └──────────────────────┘ ┃ ║
	║ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━┛ ║
	╚══════════════════════════════╝

	╔════════Tag Key Block═════════╗
	║                              ║
	║ ┏━━━━━━━━━Key List━━━━━━━━━┓ ║
	║ ┃                          ┃ ║
	║ ┃┏━━━━━━━━━━Key━━━━━━━━━━━┓┃ ║
	║ ┃┃┌──────────────────────┐┃┃ ║
	║ ┃┃│     Flag <uint8>     │┃┃ ║
	║ ┃┃├──────────────────────┤┃┃ ║
	║ ┃┃│Value Offset <uint64> │┃┃ ║
	║ ┃┃├──────────────────────┤┃┃ ║
	║ ┃┃│  len(Key) <varint>   │┃┃ ║
	║ ┃┃├──────────────────────┤┃┃ ║
	║ ┃┃│    Key <byte...>     │┃┃ ║
	║ ┃┃└──────────────────────┘┃┃ ║
	║ ┃┗━━━━━━━━━━━━━━━━━━━━━━━━┛┃ ║
	║ ┃           ...            ┃ ║
	║ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━┛ ║
	║ ┏━━━━━━━━Hash Index━━━━━━━━┓ ║
	║ ┃ ┌──────────────────────┐ ┃ ║
	║ ┃ │  len(Keys) <uint32>  │ ┃ ║
	║ ┃ ├──────────────────────┤ ┃ ║
	║ ┃ │ Key Offset <uint64>  │ ┃ ║
	║ ┃ ├──────────────────────┤ ┃ ║
	║ ┃ │         ...          │ ┃ ║
	║ ┃ └──────────────────────┘ ┃ ║
	║ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━┛ ║
	╚══════════════════════════════╝

	╔════════════Trailer══════════════╗
	║ ┌─────────────────────────────┐ ║
	║ │  Hash Index Offset <uint64> │ ║
	║ ├─────────────────────────────┤ ║
	║ │    Tag Set Size <uint64>    │ ║
	║ ├─────────────────────────────┤ ║
	║ │   Tag Set Version <uint16>  │ ║
	║ └─────────────────────────────┘ ║
	╚═════════════════════════════════╝


Measurements

	╔══════════Measurements Block═══════════╗
	║                                       ║
	║ ┏━━━━━━━━━Measurement List━━━━━━━━━━┓ ║
	║ ┃                                   ┃ ║
	║ ┃┏━━━━━━━━━━Measurement━━━━━━━━━━━┓ ┃ ║
	║ ┃┃┌─────────────────────────────┐ ┃ ┃ ║
	║ ┃┃│        Flag <uint8>         │ ┃ ┃ ║
	║ ┃┃├─────────────────────────────┤ ┃ ┃ ║
	║ ┃┃│  Tag Block Offset <uint64>  │ ┃ ┃ ║
	║ ┃┃├─────────────────────────────┤ ┃ ┃ ║
	║ ┃┃│     len(Name) <varint>      │ ┃ ┃ ║
	║ ┃┃├─────────────────────────────┤ ┃ ┃ ║
	║ ┃┃│       Name <byte...>        │ ┃ ┃ ║
	║ ┃┃├─────────────────────────────┤ ┃ ┃ ║
	║ ┃┃│    len(Series) <uint32>     │ ┃ ┃ ║
	║ ┃┃├─────────────────────────────┤ ┃ ┃ ║
	║ ┃┃│    SeriesIDs <uint32...>    │ ┃ ┃ ║
	║ ┃┃└─────────────────────────────┘ ┃ ┃ ║
	║ ┃┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ ┃ ║
	║ ┃                ...                ┃ ║
	║ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ ║
	║ ┏━━━━━━━━━━━━Hash Index━━━━━━━━━━━━━┓ ║
	║ ┃ ┌───────────────────────────────┐ ┃ ║
	║ ┃ │  len(Measurements) <uint32>   │ ┃ ║
	║ ┃ ├───────────────────────────────┤ ┃ ║
	║ ┃ │  Measurement Offset <uint64>  │ ┃ ║
	║ ┃ ├───────────────────────────────┤ ┃ ║
	║ ┃ │              ...              │ ┃ ║
	║ ┃ └───────────────────────────────┘ ┃ ║
	║ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ ║
	║ ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ ║
	║ ┃              Trailer              ┃ ║
	║ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ ║
	╚═══════════════════════════════════════╝

	╔════════════Trailer══════════════╗
	║ ┌─────────────────────────────┐ ║
	║ │  Hash Index Offset <uint64> │ ║
	║ ├─────────────────────────────┤ ║
	║ │     Block Size <uint64>     │ ║
	║ ├─────────────────────────────┤ ║
	║ │    Block Version <uint16>   │ ║
	║ └─────────────────────────────┘ ║
	╚═════════════════════════════════╝

*/
package tsi1
