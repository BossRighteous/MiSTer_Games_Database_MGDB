# update_all Game Media Support

## Format and Storage
- Unified Media folders by type in predicatable install directories
  - `docs/<system>/<media-type>/<title-slug>.<ext>`
    - boxart2d, screenshot, metadata?
    - Where is docs relative to install path, is it customizable?
- Maps 'known rom sets' (RetroArch RDB) to Zaparoo Title Slugs
  - This mapping is purely filename based
  - Slugs roll up common names under parent deduplication
  - Not 1:1 with ScreenScraper GameIDs or any other source
  - Multiple Titles may exist for a conceptual game by virtue of naming convention, but slug attempts to remedy this
- Simple binary provided that given a path or filename will return a slug of confirmed computation

## Systems Map Data
  - MiSTer {core} directory expectation
  - Retroarch DB URL
  - Skraper compatible system directory
  - Zaparoo SystemID

## Scraper File Generation
- Loop systems map
- Download RDB file and pass through libretrodb_tool for NDJSON output
- Loop NDJSON and touch an empty file for each filename into Skraper compatible directory

## Scraper Processing
- Skraper UI used with folowing profile settings:
  - (Systems should map automatically in wizard)
  - Recallbox gamelist.xml
  - Image Boxart2D
  - Image Screenshot
  - unlinked, undeduped?
  - Omit not-found?

## Scraper Result Processing
- Loop Systems Map
  - Open Skraper System Dir/ gamelist.xml
  - Loop all records
    - Skip empty ID (no scrape data found)
      - `processed/<mister-core>/missing/{filename}.ext`??
    - Skip already mapped Zaparoo Titles
  - pass relative path to Zaparoo Title Slugify method
  - copy linked or relative image data from `./media/<type>/<filname>.png` to `processed/<mister-core>/<type>/<zaparoo-slug>.png`
    - Attempt to JPEG convert anything? or just pass raw?
  - copy unmarshalled meta data from XML nodes to `processed/<mister-core>/metadata/<zaparoo-slug>.json`
  - Add title slug to processed map to skip retries

## Zaparoo DB merger
- A scraper `Merge update_all Data` will check these known directories for reference data. Tags will be relationally added to DB. A Title property token like `$UA` can be used to indicate that the record has local data, and API retrievals should merge it from disk. Unlike tags, properties are meant to be returned only in this lookup scope.

## Sharing
This gets us a declarative FS basis for Zaparoo to lazy index by expected FS path resolution. It also provides a knowable resoluton path (via slug transform binary and published slug computations) for other 1st party scripts to use the media. Zaparoo will add a scraper for 

Archive.org appears to be a backbone for this sort of sharing. Github releases is another option.

## Open Questions
- All or nothing?
  - How do users indicate their media preferences for subset data
- Updates and maintenance
  - Is this collection stored and updated on the network store? How is the DB managed from a 'complete' set.


