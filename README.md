# MiSTer_Games_Data_Utils
MGDB media database Generation Utilities for [MiSTer_Games_GUI](https://github.com/BossRighteous/MiSTer_Games_GUI) to allow GUI media browsing of your [MiSTer FPGA](https://github.com/MiSTer-devel/Wiki_MiSTer/wiki) library

MGDBs attempt to create single portable databases of rich media that would otherwise require scraping.

Instead of writing a multi-request utility to fetch media based on a local games collection; a single SQLite3 media database is used and compared to your local MiSTer Games media.

This strategy trades requests and file counts for initially padded databases.

The working theory is:
- Many MiSTer users will have NoIntro/Tosec/Redump compliant collections
- These collections are mostly stable over time
- Retroarch DB RDBs cover NoIntro/Tosec/Redump dumps for the most part
- By `touching` 'every' known filename we can create a dummy collection of all ROMs without needing to possess said ROMs
- By slugifying filenames and dropping (..) [..] content we can map and dedupe by slug match
- By only touching/scraping one known file per slug we can connect multiple files with the same slug to the same game record
- By scraping these files based on filename we can compile the 'entire collection' of media, rom, and game references
- This portable collection is easier to share, specifically via GitHub releases as larger single binaries
- SQLite is fast for the use case, and allows lower memory consumption vs loading large collections and XML into memory for a GUI. The ARM chip on the MiSTer is low resource.

## Beta release notice

As the project is in beta, MGDB binaries and schemas are subject to change. Please avoid downloading the full collections not if later replacement is a concern from a write/storage perspective.

## License and attributions
The code of the repository is MIT license.
- RDB data and utils courtesy of [Libretro](https://github.com/libretro/libretro-database) contributors
- MiSTer MGL/MRA system mappings courtesy of [Wizzomafizzo](https://github.com/wizzomafizzo/mrext/blob/main/pkg/games/systems.go)

The resulting MGDB database files are Create Commons Attribution-NonCommercial-ShareAlike 4.0 International license.
- RDB data and utils courtesy of [Libretro](https://github.com/libretro/libretro-database) contributors
- Descriptions images and media courtesy [Screenscraper.fr](https://screenscraper.fr/) contributors under Create Commons Attribution-NonCommercial-ShareAlike 4.0 International license.
- Scraping services provided by [Skraper](https://www.skraper.net/)

## Dev Usage

Script to create directories, download RDB files from libretro github, export RBD to NDJSON format
```
go run ./cmd/setuprdb/main.go {SystemID || 'all'}
```

Script to parse NDJSON files, map to unique slugs, and create a single known empty rom file per slug for scraping.
```
go run ./cmd/touchndjson/main.go {SystemID || 'all'}
```

Manual Step:
Run Skraper or equivalent on each core directory to compile 'complete meta' set in gamelist.xml

Script scan gamelist.xml, RDB info, and related images into a relational SQLite3 DB (MGDB)
```
go run ./cmd/buildmgdb/main.go {SystemID || 'all'}
```