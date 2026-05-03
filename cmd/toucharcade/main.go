package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/config"
)

const csvURL = "https://raw.githubusercontent.com/MiSTer-devel/ArcadeDatabase_MiSTer/main/ArcadeDatabase.csv"

func main() {
	dirPath := config.CommandRootPath
	arcadePath := filepath.Join(dirPath, "cores_arcade")

	if err := os.MkdirAll(arcadePath, 0755); err != nil {
		fmt.Printf("Unable to create directory %s: %v\n", arcadePath, err)
		return
	}

	csvPath := filepath.Join(arcadePath, "ArcadeDatabase.csv")
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Printf("Downloading CSV to %s\n", csvPath)
		if err := downloadFile(csvURL, csvPath); err != nil {
			fmt.Printf("Unable to download CSV: %v\n", err)
			return
		}
	}

	f, err := os.Open(csvPath)
	if err != nil {
		fmt.Printf("Unable to open CSV %s: %v\n", csvPath, err)
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	// skip header
	if _, err := reader.Read(); err != nil {
		fmt.Printf("Unable to read CSV header: %v\n", err)
		return
	}

	rowCount := 0
	fileCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("CSV read error: %v\n", err)
			continue
		}
		if len(record) == 0 || record[0] == "" {
			rowCount++
			continue
		}
		rowCount++
		zipPath := filepath.Join(arcadePath, record[0]+".zip")
		fo, err := os.Create(zipPath)
		if err != nil {
			fmt.Printf("Unable to create file %s: %v\n", zipPath, err)
			continue
		}
		if err := fo.Close(); err != nil {
			fmt.Printf("Unable to close file %s: %v\n", zipPath, err)
			continue
		}
		fileCount++
	}

	fmt.Printf("CSV rows: %d, files created: %d\n", rowCount, fileCount)
}

func downloadFile(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
