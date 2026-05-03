package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/config"
)

const xmlHeader = `<?xml version="1.0" encoding="utf-8" standalone="yes"?>` + "\n"

type Gamelist struct {
	XMLName  xml.Name `xml:"gameList"`
	Provider Provider `xml:"provider"`
	Games    []Game   `xml:"game"`
}

type Provider struct {
	System   string `xml:"System"`
	Software string `xml:"software"`
	Database string `xml:"database"`
	Web      string `xml:"web"`
}

type Game struct {
	ID          string `xml:"id,attr"`
	Source      string `xml:"source,attr"`
	Path        string `xml:"path"`
	Name        string `xml:"name"`
	Desc        string `xml:"desc"`
	Rating      string `xml:"rating,omitempty"`
	ReleaseDate string `xml:"releasedate,omitempty"`
	Developer   string `xml:"developer,omitempty"`
	Publisher   string `xml:"publisher,omitempty"`
	Genre       string `xml:"genre,omitempty"`
	Players     string `xml:"players,omitempty"`
	Image       string `xml:"image,omitempty"`
	Thumbnail   string `xml:"thumbnail,omitempty"`
	GenreID     string `xml:"genreid,omitempty"`
}

type arcadeEntry struct {
	setname      string
	name         string
	manufacturer string
}

func main() {
	dirPath := config.CommandRootPath
	arcadePath := filepath.Join(dirPath, "cores_arcade")

	csvPath := filepath.Join(arcadePath, "ArcadeDatabase.csv")
	xmlPath := filepath.Join(arcadePath, "gamelist.xml")

	entries, err := loadCSV(csvPath)
	if err != nil {
		fmt.Printf("Unable to load CSV %s: %v\n", csvPath, err)
		return
	}

	data, err := os.ReadFile(xmlPath)
	if err != nil {
		fmt.Printf("Unable to read %s: %v\n", xmlPath, err)
		return
	}

	gl := &Gamelist{}
	if err := xml.Unmarshal(data, gl); err != nil {
		fmt.Printf("Unable to parse gamelist.xml: %v\n", err)
		return
	}

	var merged []Game
	droppedEmptyID := 0
	droppedNoCSV := 0
	remapped := 0

	for _, game := range gl.Games {
		if game.ID == "" {
			droppedEmptyID++
			continue
		}
		setname := extractSetname(game.Path)
		fmt.Println("path setname", setname)
		entry, ok := entries[setname]
		if !ok || entry.setname == "" {
			droppedNoCSV++
			continue
		}
		game.Path = "./" + entry.name + ".mra"
		if entry.manufacturer != "" {
			game.Developer = entry.manufacturer
			game.Publisher = entry.manufacturer
		}
		remapped++
		merged = append(merged, game)
	}

	gl.Games = merged

	out, err := xml.MarshalIndent(gl, "", "  ")
	if err != nil {
		fmt.Printf("Unable to marshal XML: %v\n", err)
		return
	}

	if err := os.WriteFile(xmlPath, append([]byte(xmlHeader), out...), 0644); err != nil {
		fmt.Printf("Unable to write %s: %v\n", xmlPath, err)
		return
	}

	fmt.Printf("Remapped: %d, dropped (empty id): %d, dropped (no CSV match): %d\n",
		remapped, droppedEmptyID, droppedNoCSV)
}

func extractSetname(path string) string {
	base := filepath.Base(path)
	if !strings.HasSuffix(base, ".zip") {
		return ""
	}
	return strings.TrimSuffix(base, ".zip")
}

func loadCSV(csvPath string) (map[string]arcadeEntry, error) {
	f, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}

	entries := make(map[string]arcadeEntry)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if len(record) < 12 || record[0] == "" {
			continue
		}
		entries[record[0]] = arcadeEntry{
			setname:      record[0],
			name:         record[1],
			manufacturer: record[11],
		}
	}
	return entries, nil
}
