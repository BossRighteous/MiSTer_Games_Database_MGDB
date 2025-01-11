package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/config"
	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/gamelist"
	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/mgdb"
	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/rdb"
	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/sqlite"
)

func main() {

	cliArgs := os.Args
	fmt.Println(cliArgs)
	if len(cliArgs) < 2 {
		fmt.Println("No DataConfig key argument provided")
		return
	}
	configKey := cliArgs[1]

	// keyword to process all in sequence
	if configKey == "all" {
		for _, dataConfig := range config.DataConfigs {
			buildMGDB(dataConfig)
		}
		return
	}

	// Else try single
	dataConfig, ok := config.DataConfigs[configKey]
	if !ok {
		fmt.Println("Invalid DataConfig key")
		return
	}
	buildMGDB(dataConfig)
}

func makeRDBMap(corePath string) map[string]int {
	rdbMap := make(map[string]int)
	rdbPath := filepath.Join(corePath, "libretro.rdb")
	fmt.Printf("Opening %s\n", rdbPath)
	rdbFile, err := os.Open(rdbPath)
	rdbBytes := make([]byte, 0)
	if err != nil {
		fmt.Println("Unableto Open RDB")
		return rdbMap
	}
	bytes, err := io.ReadAll(rdbFile)
	if err != nil && err != io.EOF {
		fmt.Println("Unable parse RDB")
		return rdbMap
	}
	rdbBytes = append(rdbBytes, bytes...)
	if len(rdbBytes) == 0 {
		fmt.Println("RDB empty")
		return rdbMap
	}
	games := rdb.Parse(rdbBytes)
	for _, game := range games {
		fmt.Println(game.ROMName, game.CRC32)
		rdbMap[game.ROMName] = int(game.CRC32)
	}
	return rdbMap
}

func buildMGDB(dataConfig config.DataConfig) {
	coresPath := "./cores/"
	coreDir := dataConfig.ScrapeFolder
	corePath := filepath.Join(coresPath, coreDir)

	systemIds := make([]string, len(dataConfig.Systems))
	for i, system := range dataConfig.Systems {
		systemIds[i] = system.Id
	}

	primarySystem := dataConfig.Systems[0]
	mgdbFilename := fmt.Sprintf("%v (%v) (%v)", primarySystem.Name, primarySystem.Category, primarySystem.ReleaseDate)
	mgdbFilename = strings.ReplaceAll(mgdbFilename, "/", "-")

	dbInfo := mgdb.MGDBInfo{
		CollectionName:     mgdbFilename,
		GamesFolder:        coreDir,
		SupportedSystemIds: strings.Join(systemIds, ","),
		BuildDate:          time.Now().Format("2006-01-02"),
		MGDBVersion:        "1.0",
		Description:        "Compiled for MiSTer_Games_GUI by @BossRighteous.\nMedia courtesy https://screenscraper.fr/ contributors and sources made available under Create Commons Attribution-NonCommercial-ShareAlike 4.0 International.\nROM data courtesy Libretro under Creative Commons Attribution-ShareAlike 4.0 International.",
	}

	// Open gameslist.xml
	gamelistFile, err := os.Open(filepath.Join(corePath, "gamelist.xml"))
	if err != nil {
		fmt.Println("Unable to open gamelist.xml file")
		panic(err)
	}
	gamelistBytes, err := io.ReadAll(gamelistFile)
	if err != nil && err != io.EOF {
		fmt.Println("Unable to read gamelist.xml file")
		panic(err)
	}

	// Parse gamelist into usable structs
	gamelist := gamelist.ParseGamelist(gamelistBytes)
	// Allocate table maps by pk
	// Consider reindexing game ID

	reindexedGames := []mgdb.Game{
		{
			GameID:      0,
			Name:        "~Unknown",
			Description: "Loose ROMs not matched by MGDB data",
		},
	}
	gameMap := make(map[int]int)
	gameMap[0] = 0
	romMap := make(map[string]mgdb.GamelistRom)
	romMap[""] = mgdb.GamelistRom{}
	reindexedGenres := []mgdb.Genre{
		{GenreID: 0, Name: "~Unknown"},
	}
	genreMap := make(map[int]int)
	genreMap[0] = 0

	// Mapping for binary blobs
	screenshotMap := make(map[int]string) // [gameId]imagePath
	screenshotMap[0] = ""
	titleScreenMap := make(map[int]string) //[gameId]imagePath
	titleScreenMap[0] = ""
	blobHashMap := make(map[string]bool) // [hash]exists
	rdbMap := makeRDBMap(corePath)

	reAscii := regexp.MustCompile("[[:^ascii:]]")

	// Reinflate map of RDB filenames CRC32s since we have them

	// Reorganize into table maps by game.ID
	for _, game := range gamelist.Games {
		fmt.Printf("%+v\n", game)

		glGameID, err := strconv.Atoi(game.ID)
		if err != nil || glGameID == 0 {
			fmt.Println("Unable to parse Game.ID as int, skipping")
			continue
		}

		glGenreID := 0
		glGenreID, _ = strconv.Atoi(game.GenreID)

		genreID := 0
		if foundGenreID, ok := genreMap[glGenreID]; !ok {
			genreID = len(reindexedGenres)
			reindexedGenres = append(reindexedGenres, mgdb.Genre{
				GenreID: genreID,
				Name:    game.Genre,
			})
			genreMap[glGenreID] = genreID
		} else {
			genreID = foundGenreID
		}

		gameID := 0
		if foundGameID, ok := gameMap[glGameID]; !ok {
			gameID = len(reindexedGames)

			// Reformat date string
			fmtReleaseDate := ""
			if len(game.ReleaseDate) >= 8 {
				fmtReleaseDate = fmt.Sprintf("%v-%v-%v", game.ReleaseDate[0:4], game.ReleaseDate[4:6], game.ReleaseDate[6:8])
			}

			// Drop Non-ASCII
			fmtDescription := reAscii.ReplaceAllLiteralString(game.Desc, "")

			reindexedGames = append(reindexedGames, mgdb.Game{
				GameID:      gameID,
				Name:        game.Name,
				IsIndexed:   0,
				GenreId:     genreID,
				Description: fmtDescription,
				Rating:      game.Rating,
				ReleaseDate: fmtReleaseDate,
				Developer:   game.Developer,
				Publisher:   game.Publisher,
			})
			gameMap[glGameID] = gameID
		} else {
			gameID = foundGameID
		}

		// Polyfil
		HasSuffix := func(s, suffix string) bool {
			return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
		}
		CutSuffix := func(s, suffix string) (before string, found bool) {
			if !HasSuffix(s, suffix) {
				return s, false
			}
			return s[:len(s)-len(suffix)], true
		}

		// configure romname as filename without extension
		fileBase := filepath.Base(game.Path)
		fileExt := filepath.Ext(game.Path)
		filename, _ := CutSuffix(fileBase, fileExt)
		if _, ok := romMap[filename]; !ok {
			crc := 0
			if rdbCrc, ok := rdbMap[fileBase]; ok {
				crc = rdbCrc
			}

			romMap[filename] = mgdb.GamelistRom{
				FileName:           filename,
				GameID:             gameID,
				SupportedSystemIds: "",
				CRC32:              crc,
			}
		}
		// Add slug for fuzzy match
		slug := mgdb.SlugifyString(filename)
		if _, ok := romMap[slug]; !ok {
			romMap[slug] = mgdb.GamelistRom{
				FileName:           slug,
				GameID:             gameID,
				SupportedSystemIds: "",
			}
		}

		// For initial Map, save full path, will read bytes and decompose later
		if _, ok := screenshotMap[gameID]; !ok && game.Image != "" {
			screenshotMap[gameID] = game.Image
		}

		// For initial Map, save full path, will read bytes and decompose later
		if _, ok := titleScreenMap[gameID]; !ok && game.Thumbnail != "" {
			titleScreenMap[gameID] = game.Thumbnail
		}
	}

	dbPath := filepath.Join(corePath, mgdbFilename+".mgdb")
	db, err := sqlite.CreateMGDB(dbPath)
	if err != nil {
		fmt.Println("Unable to allocate DB at ", dbPath)
		panic(err)
	}

	sqlite.InsertMGDBInfo(db, dbInfo)
	sqlite.BulkInsertGames(db, reindexedGames)
	sqlite.BulkInsertGenres(db, reindexedGenres)
	sqlite.BulkInsertGamelistRoms(db, romMap)
	sqlite.BulkInsertImageMap(db, "Screenshot", screenshotMap, blobHashMap, corePath)
	sqlite.BulkInsertImageMap(db, "TitleScreen", titleScreenMap, blobHashMap, corePath)
	fmt.Println("MGDB Built Successfully")
	db.Close()
}
