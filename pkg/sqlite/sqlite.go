package sqlite

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
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
		MGDBVersion text not null,
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
		DeveloperID integer not null,
		PublisherID integer not null,
		Players text not null,
		ExternalID text not null,
		ScreenshotHash text,
		TitleScreenHash text
	 );
	 CREATE INDEX game_name_idx ON Game (Name);
	 CREATE INDEX game_genre_idx ON Game (GenreID);
	 CREATE INDEX game_developer_idx ON Game (DeveloperID);
	 CREATE INDEX game_publisher_idx ON Game (PublisherID);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	// Slugs were derived from full deduped list
	sqlStmt = `
	drop table if exists SlugRom;
	create table SlugRom (
		Slug text primary key not null,
		GameID integer not null,
		SupportedSystemIds text not null
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	// CRC32 hex strings for possible slug mapping later
	sqlStmt = `
	drop table if exists RomCrc;
	create table RomCrc (
		CRC32 text primary key not null,
		Slug text not null
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
		GameID integer not null,
		SupportedSystemIds text not null
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
	);
	CREATE INDEX genre_name_idx ON Genre (Name);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	sqlStmt = `
	drop table if exists Developer;
	create table Developer (
		DeveloperID integer primary key not null,
		Name text not null
	);
	CREATE INDEX developer_name_idx ON Developer (Name);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	sqlStmt = `
	drop table if exists Publisher;
	create table Publisher (
		PublisherID integer primary key not null,
		Name text not null
	);
	CREATE INDEX publisher_name_idx ON Publisher (Name);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	sqlStmt = `
	drop table if exists ImageBlob;
	create table ImageBlob (
		Hash text primary key not null,
		Bytes blob not null
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	db.Exec("pragma synchronous = off")
	db.Exec("pragma journal_mode = off")
	db.Exec("pragma page_size = 4096")
	return db, nil
}

func GetMD5Hash(bytes []byte) string {
	hash := md5.Sum(bytes)
	return hex.EncodeToString(hash[:])
}

func CreateMGDB(path string) (*sql.DB, error) {
	db, err := allocDB(path)
	if err != nil {
		return nil, err
	}
	//db.Exec("VACUUM")
	return db, nil
}

func Vacuum(db *sql.DB) {
	sqlStmt := `VACUUM;`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		fmt.Println("Vacuum Failed")
	}
	fmt.Println("Vacuum Executed")
}

func InsertMGDBInfo(db *sql.DB, info mgdb.MGDBInfo) {
	stmt, err := db.Prepare(
		"insert into MGDBInfo(" +
			"CollectionName, GamesFolder, SupportedSystemIds, BuildDate, MGDBVersion, Description" +
			") values (?, ?, ?, ?, ?, ?)",
	)
	if err != nil {
		panic("InsertMGDBInfo Prepare")
	}
	_, err = stmt.Exec(
		info.CollectionName,
		info.GamesFolder,
		info.SupportedSystemIds,
		info.BuildDate,
		info.MGDBVersion,
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
				"GameID, Name, IsIndexed, GenreID, Rating, ReleaseDate, " +
				"DeveloperID, PublisherID, Players, Description, ExternalID" +
				") values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		)
		if err != nil {
			fmt.Printf("%+v\n", game)
			panic("BulkInsertGames Prepare")
		}
		_, err = stmt.Exec(
			game.GameID,
			game.Name,
			game.IsIndexed,
			game.GenreID,
			game.Rating,
			game.ReleaseDate,
			game.DeveloperID,
			game.PublisherID,
			game.Players,
			game.Description,
			game.ExternalID,
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

func BulkInsertDevelopers(db *sql.DB, developers []mgdb.Developer) {
	for _, developer := range developers {
		fmt.Println("adding Developer", developer.Name)
		stmt, err := db.Prepare(
			"insert into Developer(" +
				"DeveloperID, Name" +
				") values (?, ?)",
		)
		if err != nil {
			fmt.Printf("%+v\n", developer)
			panic("BulkInsertDevelopers Prepare")
		}
		_, err = stmt.Exec(
			developer.DeveloperID,
			developer.Name,
		)
		if err != nil {
			fmt.Printf("%+v\n", developer)
			panic("BulkInsertDevelopers Exec")
		}
	}
}

func BulkInsertPublishers(db *sql.DB, publishers []mgdb.Publisher) {
	for _, publisher := range publishers {
		fmt.Println("adding Publisher", publisher.Name)
		stmt, err := db.Prepare(
			"insert into Publisher(" +
				"PublisherID, Name" +
				") values (?, ?)",
		)
		if err != nil {
			fmt.Printf("%+v\n", publisher)
			panic("BulkInsertPublishers Prepare")
		}
		_, err = stmt.Exec(
			publisher.PublisherID,
			publisher.Name,
		)
		if err != nil {
			fmt.Printf("%+v\n", publisher)
			panic("BulkInsertPublishers Exec")
		}
	}
}

func BulkInsertSlugRoms(db *sql.DB, slugRomMap map[string]mgdb.SlugRom) {
	for _, rom := range slugRomMap {
		fmt.Println("adding SlugRom", rom.Slug)
		stmt, err := db.Prepare(
			"insert into SlugRom(" +
				"Slug, GameID, SupportedSystemIds" +
				") values (?, ?, ?)",
		)
		if err != nil {
			fmt.Printf("%+v\n", rom)
			panic("BulkInsertSlugRoms Prepare")
		}
		_, err = stmt.Exec(
			rom.Slug,
			rom.GameID,
			rom.SupportedSystemIds,
		)
		if err != nil {
			fmt.Printf("%+v\n", rom)
			panic("BulkInsertSlugRoms Exec")
		}
	}
}

func BulkInsertRomCrcs(db *sql.DB, romCrcs []mgdb.RomCrc) {
	for _, rom := range romCrcs {
		fmt.Println("adding RomCrc", rom.CRC32, rom.Slug)
		stmt, err := db.Prepare(
			"insert into RomCrc(" +
				"CRC32, Slug" +
				") values (?, ?)",
		)
		if err != nil {
			fmt.Printf("%+v\n", rom)
			panic("BulkInsertRomCrcs Prepare")
		}
		_, err = stmt.Exec(
			rom.CRC32,
			rom.Slug,
		)
		if err != nil {
			fmt.Printf("%+v\n", rom)
			panic("BulkInsertRomCrcs Exec")
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

func BulkInsertImageMap(db *sql.DB, imgType string, imageMap map[int]string, md5Map map[string]bool, basePath string) {
	for gameID, filePath := range imageMap {
		if gameID == 0 || filePath == "" {
			continue
		}
		fmt.Printf("adding image %v %v\n", imgType, filePath)
		imgPath := filepath.Join(basePath, filePath)
		blob := safeLoadFileBytes(imgPath)
		if blob == nil {
			continue
		}

		// Compare hash for dedupe
		hash := GetMD5Hash(blob)
		if _, ok := md5Map[hash]; !ok {
			// ADD IMAGE RECORD
			stmt, err := db.Prepare(
				"insert into ImageBlob(Hash, Bytes) values (?, ?)",
			)
			if err != nil {
				fmt.Printf("%v %v", hash, filePath)
				panic("BulkInsertImageMap ImageBlob Prepare")
			}
			_, err = stmt.Exec(
				hash,
				blob,
			)
			if err != nil {
				fmt.Printf("%+v\n", filePath)
				panic("BulkInsertImageMap ImageBlob Exec")
			}

			md5Map[hash] = true
		}

		// Update game with hash by type

		column := "ScreenshotHash"
		if imgType == "TitleScreen" {
			column = "TitleScreenHash"
		}

		stmt, err := db.Prepare(
			"update Game set " + column + " = ? where GameID = ?",
		)
		if err != nil {
			fmt.Printf("%v %v\n", gameID, hash)
			panic("BulkInsertImageMap Game Update Prepare")
		}
		_, err = stmt.Exec(
			hash,
			gameID,
		)
		if err != nil {
			fmt.Printf("%v %v\n", gameID, hash)
			panic("BulkInsertScreenshots Game Update Exec")
		}
	}
}
