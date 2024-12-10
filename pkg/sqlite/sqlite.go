package sqlite

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/mgdb"
	_ "github.com/mattn/go-sqlite3"
)

func allocDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	sqlStmt := `	
	drop table if exists MGDBInfo;
	create table if not exists MGDBInfo (
		CollectionName text not null,
		GamesFolder text not null,
		SupportedSystemIds text not null,
		BuildDate text not null,
		Description text not null
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	sqlStmt = `
	drop table if exists Game;
	create table Game (
		GameID integer primary key not null,
		Name text not null,
		IsIndexed integer not null,
		GenreID integer not null,
		Description text not null,
		Rating text not null,
		ReleaseDate text not null,
		Developer text not null,
		Publisher text not null,
		Players text not null
	 ) WITHOUT ROWID;
	 CREATE INDEX game_name_idx ON Game (Name);
	 CREATE INDEX game_genre_idx ON Game (GenreID);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	// I don't think we'll really need this, plus zip, chd issues per ext
	/*
		sqlStmt = `
		drop table if exists RDBRom;
		create table RDBRom (
			ROMName text primary key not null,
			Name text not null,
			CRC32 integer not null,
			Size integer not null,
			Serial text not null,
			GameID integer not null
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return db, err
		}
	*/

	// FileName without ext should be good match
	// Slugs can go here too!
	sqlStmt = `
	drop table if exists GamelistRom;
	create table GamelistRom (
		FileName text primary key not null,
		GameID integer not null
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	sqlStmt = `
	drop table if exists IndexedRom;
	create table IndexedRom (
		Path text primary key not null,
		FileName text not null,
		FileExt text not null,
		GameID integer not null
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	sqlStmt = `
	drop table if exists Genre;
	create table Genre (
		GenreID integer primary key not null,
		Name text not null
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	sqlStmt = `
	drop table if exists Screenshot;
	create table Screenshot (
		GameID integer primary key not null,
		FilePath text not null,
		Bytes blob
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	sqlStmt = `
	drop table if exists TitleScreen;
	create table TitleScreen (
		GameID integer primary key not null,
		FilePath text not null,
		Bytes blob
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	db.Exec("pragma synchronous = off")
	db.Exec("pragma journal_mode = off")
	//db.Exec("pragma pagesize = 1024")
	return db, nil
}

func CreateMGDB(path string) (*sql.DB, error) {
	db, err := allocDB(path)
	if err != nil {
		return nil, err
	}

	//db.Exec("VACUUM")
	return db, nil
}

func InsertMGDBInfo(db *sql.DB, info mgdb.MGDBInfo) {
	stmt, err := db.Prepare(
		"insert into MGDBInfo(" +
			"CollectionName, GamesFolder, SupportedSystemIds, BuildDate, Description" +
			") values (?, ?, ?, ?, ?)",
	)
	if err != nil {
		panic("InsertMGDBInfo Prepare")
	}
	_, err = stmt.Exec(
		info.CollectionName,
		info.GamesFolder,
		info.SupportedSystemIds,
		info.BuildDate,
		info.Description,
	)
	if err != nil {
		panic("InsertMGDBInfo Exec")
	}
}

func BulkInsertGames(db *sql.DB, games []mgdb.Game) {
	for _, game := range games {
		fmt.Println("adding game", game.Name)
		stmt, err := db.Prepare(
			"insert into Game(" +
				"GameID, Name, IsIndexed, GenreId, Rating, ReleaseDate, " +
				"Developer, Publisher, Players, Description" +
				") values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		)
		if err != nil {
			fmt.Printf("%+v\n", game)
			panic("BulkInsertGames Prepare")
		}
		_, err = stmt.Exec(
			game.GameID,
			game.Name,
			game.IsIndexed,
			game.GenreId,
			game.Rating,
			game.ReleaseDate,
			game.Developer,
			game.Publisher,
			game.Players,
			game.Description,
		)
		if err != nil {
			fmt.Printf("%+v\n", game)
			panic("BulkInsertGames Exec")
		}
	}

}

func BulkInsertGenres(db *sql.DB, genres []mgdb.Genre) {
	for _, genre := range genres {
		fmt.Println("adding genre", genre.Name)
		stmt, err := db.Prepare(
			"insert into Genre(" +
				"GenreID, Name" +
				") values (?, ?)",
		)
		if err != nil {
			fmt.Printf("%+v\n", genre)
			panic("BulkInsertGenres Prepare")
		}
		_, err = stmt.Exec(
			genre.GenreID,
			genre.Name,
		)
		if err != nil {
			fmt.Printf("%+v\n", genre)
			panic("BulkInsertGenres Exec")
		}
	}
}

func BulkInsertGamelistRoms(db *sql.DB, gamelistRoms map[string]mgdb.GamelistRom) {
	for _, rom := range gamelistRoms {
		fmt.Println("adding rom", rom.FileName)
		stmt, err := db.Prepare(
			"insert into GamelistRom(" +
				"FileName, GameID" +
				") values (?, ?)",
		)
		if err != nil {
			fmt.Printf("%+v\n", rom)
			panic("BulkInsertGamelistRoms Prepare")
		}
		_, err = stmt.Exec(
			rom.FileName,
			rom.GameID,
		)
		if err != nil {
			fmt.Printf("%+v\n", rom)
			panic("BulkInsertGamelistRoms Exec")
		}
	}
}

func safeLoadFileBytes(path string) []byte {
	var b []byte
	imgFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Unable to open file path", path)
		return b
	}
	imageBytes, err := io.ReadAll(imgFile)
	if err != nil && err != io.EOF {
		fmt.Println("Unable to file file", path)
		return b
	}
	return imageBytes
}

func BulkInsertScreenshots(db *sql.DB, screenshots map[int]mgdb.Screenshot, basePath string) {
	for _, screen := range screenshots {
		fmt.Println("adding screenshot", screen.FilePath)
		screenPath := filepath.Join(basePath, screen.FilePath)
		blob := safeLoadFileBytes(screenPath)
		if blob != nil {
			screen.FilePath = ""
		}

		stmt, err := db.Prepare(
			"insert into Screenshot(" +
				"FilePath, GameID, Bytes" +
				") values (?, ?, ?)",
		)
		if err != nil {
			fmt.Printf("%+v\n", screen)
			panic("BulkInsertScreenshots Prepare")
		}
		_, err = stmt.Exec(
			screen.FilePath,
			screen.GameID,
			blob,
		)
		if err != nil {
			fmt.Printf("%+v\n", screen)
			panic("BulkInsertScreenshots Exec")
		}
	}
}

func BulkInsertTitleScreens(db *sql.DB, titleScreens map[int]mgdb.TitleScreen, basePath string) {
	for _, screen := range titleScreens {
		fmt.Println("adding titlescreen", screen.FilePath)
		screenPath := filepath.Join(basePath, screen.FilePath)
		blob := safeLoadFileBytes(screenPath)
		if blob != nil {
			screen.FilePath = ""
		}

		stmt, err := db.Prepare(
			"insert into TitleScreen(" +
				"FilePath, GameID, Bytes" +
				") values (?, ?, ?)",
		)
		if err != nil {
			fmt.Printf("%+v\n", screen)
			panic("BulkInsertTitleScreens Prepare")
		}
		_, err = stmt.Exec(
			screen.FilePath,
			screen.GameID,
			blob,
		)
		if err != nil {
			fmt.Printf("%+v\n", screen)
			panic("BulkInsertTitleScreens Exec")
		}
	}
}
