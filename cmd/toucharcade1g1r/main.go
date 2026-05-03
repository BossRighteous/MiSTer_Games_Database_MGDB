package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/config"
	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/gamelist"
)

func main() {
	dirPath := config.CommandRootPath
	arcadePath := filepath.Join(dirPath, "cores_arcade")
	csvPath := filepath.Join(arcadePath, "ArcadeDatabase.csv")
	xmlPath := filepath.Join(arcadePath, "gamelist.xml")

	// setname -> true if parent_title is empty (i.e. this is a parent ROM)
	parentSetnames, err := loadParentSetnames(csvPath)
	if err != nil {
		fmt.Printf("Unable to load ArcadeDatabase.csv: %v\n", err)
		return
	}

	data, err := os.ReadFile(xmlPath)
	if err != nil {
		fmt.Printf("Unable to read %s: %v\n", xmlPath, err)
		return
	}

	gl := gamelist.ParseGamelist(data)

	outDir := filepath.Join(dirPath, "cores_arcade_1g1r")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Printf("Unable to create directory %s: %v\n", outDir, err)
		return
	}

	for _, game := range gl.Games {
		if game.ID == "" || game.ID == "0" || game.Path == "" {
			continue
		}
		if !parentSetnames[setname(game.Path)] {
			continue
		}
		destPath := filepath.Join(outDir, game.Path)
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			fmt.Printf("Unable to create directory %s: %v\n", destDir, err)
			continue
		}
		fmt.Printf("Creating 1G1R file %s\n", destPath)
		fo, err := os.Create(destPath)
		if err != nil {
			fmt.Printf("Unable to create file %s: %v\n", destPath, err)
			continue
		}
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}
}

// loadParentSetnames reads ArcadeDatabase.csv and returns a set of setnames
// whose parent_title column is empty (indicating they are the parent ROM).
func loadParentSetnames(csvPath string) (map[string]bool, error) {
	f, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	// read header to find column indices
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}
	setnameIdx, parentIdx := -1, -1
	for i, col := range header {
		switch col {
		case "setname":
			setnameIdx = i
		case "parent_title":
			parentIdx = i
		}
	}
	if setnameIdx < 0 || parentIdx < 0 {
		return nil, fmt.Errorf("required columns not found in CSV header")
	}

	parents := make(map[string]bool)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("CSV read error: %v\n", err)
			continue
		}
		if len(record) <= setnameIdx || record[setnameIdx] == "" {
			continue
		}
		if record[parentIdx] == "" {
			parents[record[setnameIdx]] = true
		}
	}
	return parents, nil
}

// setname extracts the ROM name from a gamelist path like "./1941.zip" -> "1941"
func setname(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

