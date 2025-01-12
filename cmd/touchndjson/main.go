package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/config"
	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/rdb"
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
			parseNDJSONAndTouch(dataConfig)
		}
		return
	}

	// Else try single
	dataConfig, ok := config.DataConfigs[configKey]
	if !ok {
		fmt.Println("Invalid DataConfig key")
		return
	}
	parseNDJSONAndTouch(dataConfig)
}

func parseNDJSONAndTouch(dataConfig config.DataConfig) {
	dirPath := config.CommandRootPath
	coresPath := filepath.Join(dirPath, "cores")
	corePath := filepath.Join(coresPath, dataConfig.ScrapeFolder)

	roms, err := loadNDJSON(corePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	dupeMap := rdb.MapDupeROMSlugs(roms)

	for _, romName := range dupeMap {
		// Touch file, start empty
		romsPath := filepath.Join(corePath, "roms")
		romPath := filepath.Join(romsPath, romName)
		fmt.Printf("Creating Game file %s\n", romPath)
		fo, err := os.Create(romPath)
		if err != nil {
			fmt.Printf("Unable to write file %s\n", romPath)
			continue
		}
		if err := fo.Close(); err != nil {
			panic(err)
		}
		fmt.Printf("Closing Game file %s\n", romPath)
	}
}

func loadNDJSON(corePath string) ([]rdb.RdbJsonROM, error) {
	return rdb.LoadNDJSON(corePath)
}
