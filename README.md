# MiSTer_Games_Data_Utils
Data Generation Utilities for MiSTer_Games_GUI

WIP for MGDB file generation and replacement. Will upload stable packages into latest release as MGDB schema expectations finalized.

## Usage
I'm using screenscraper.fr to generate meta packages

Download RDB files from libretro, looping to 'touch' empty file for each game
```
go run ./cmd/touchfiles/main.go
```

Run Skraper or equivalent on each core directory to compile complete meta set in gamelist.xml

Merge data from gamelist.xml and libretro.rdb into standalone sqlite DB. \
This includes image blobs
```
go run ./cmd/buildsql/main.go
```

Download rdb file(s)
parse each file 'touching' unique filenames to create a name-ref for Skraper
Skraper config tbd. Currently cropping screenshot and title to 320x240 max but this may be a mistake