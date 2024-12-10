package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/config"
	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/gamelist"
	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/mgdb"
	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/sqlite"
)

func main() {

	cliArgs := os.Args
	fmt.Println(cliArgs)
	coresPath := "./cores/"
	if len(cliArgs) < 2 {
		fmt.Println("No DataConfig key argument provided")
		return
	}
	configKey := cliArgs[1]
	dataConfig, ok := config.DataConfigs[configKey]
	if !ok {
		fmt.Println("Invalid DataConfig key")
		return
	}
	coreDir := dataConfig.ScrapeFolder
	corePath := filepath.Join(coresPath, coreDir)

	systemIds := make([]string, len(dataConfig.Systems))
	for i, system := range dataConfig.Systems {
		systemIds[i] = system.Id
	}

	dbInfo := mgdb.MGDBInfo{
		CollectionName:     dataConfig.Systems[0].Name,
		GamesFolder:        coreDir,
		SupportedSystemIds: strings.Join(systemIds, ","),
		BuildDate:          time.Now().Format("2006-01-02"),
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
	reindexedGames := []mgdb.Game{{
		GameID:      0,
		Name:        "Unknown ROMs",
		Description: "Unknown ROMs",
	}}
	gameMap := make(map[int]int)
	gameMap[0] = 0
	romMap := make(map[string]mgdb.GamelistRom)
	romMap[""] = mgdb.GamelistRom{}
	reindexedGenres := make([]mgdb.Genre, 1)
	genreMap := make(map[int]int)
	genreMap[0] = 0
	screenshotMap := make(map[int]mgdb.Screenshot)
	screenshotMap[0] = mgdb.Screenshot{}
	titleScreenMap := make(map[int]mgdb.TitleScreen)
	titleScreenMap[0] = mgdb.TitleScreen{}

	// TODO: Reinflate map of RDB filenames for system?

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
			reindexedGames = append(reindexedGames, mgdb.Game{
				GameID:      gameID,
				Name:        game.Name,
				IsIndexed:   0,
				GenreId:     genreID,
				Description: game.Desc,
				Rating:      game.Rating,
				ReleaseDate: game.ReleaseDate,
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
			romMap[filename] = mgdb.GamelistRom{
				FileName: filename,
				GameID:   gameID,
			}
		}
		// Add slug for fuzzy match
		slug := mgdb.SlugifyString(filename)
		if _, ok := romMap[slug]; !ok {
			romMap[slug] = mgdb.GamelistRom{
				FileName: slug,
				GameID:   gameID,
			}
		}

		// For initial Map, save full path, will read bytes and decompose later
		if _, ok := screenshotMap[gameID]; !ok && game.Image != "" {
			screenshotMap[gameID] = mgdb.Screenshot{
				GameID:   gameID,
				FilePath: game.Image,
			}
		}

		// For initial Map, save full path, will read bytes and decompose later
		if _, ok := titleScreenMap[gameID]; !ok && game.Thumbnail != "" {
			titleScreenMap[gameID] = mgdb.TitleScreen{
				GameID:   gameID,
				FilePath: game.Thumbnail,
			}
		}
	}

	dbPath := filepath.Join(corePath, coreDir+".mgdb")
	db, err := sqlite.CreateMGDB(dbPath)
	if err != nil {
		fmt.Println("Unable to allocate DB at ", dbPath)
		panic(err)
	}

	sqlite.InsertMGDBInfo(db, dbInfo)
	sqlite.BulkInsertGames(db, reindexedGames)
	sqlite.BulkInsertGenres(db, reindexedGenres)
	sqlite.BulkInsertGamelistRoms(db, romMap)
	sqlite.BulkInsertScreenshots(db, screenshotMap, corePath)
	sqlite.BulkInsertTitleScreens(db, titleScreenMap, corePath)
	fmt.Println("MGDB Built Successfully")
	db.Close()
}
