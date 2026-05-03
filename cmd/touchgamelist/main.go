package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/config"
	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/gamelist"
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
			parseGamelistAndTouch(dataConfig)
		}
		return
	}

	// Else try single
	dataConfig, ok := config.DataConfigs[configKey]
	if !ok {
		fmt.Println("Invalid DataConfig key")
		return
	}
	parseGamelistAndTouch(dataConfig)
}

func parseGamelistAndTouch(dataConfig config.DataConfig) {
	dirPath := config.CommandRootPath
	coresPath := filepath.Join(dirPath, "cores")
	corePath := filepath.Join(coresPath, dataConfig.ScrapeFolder)
	xmlPath := filepath.Join(corePath, "gamelist.xml")

	data, err := os.ReadFile(xmlPath)
	if err != nil {
		fmt.Printf("Unable to read %s: %v\n", xmlPath, err)
		return
	}

	gl := gamelist.ParseGamelist(data)

	for _, game := range gl.Games {
		if game.Path == "" {
			continue
		}
		gamePath := filepath.Join(corePath, game.Path)
		fmt.Printf("Creating Game file %s\n", gamePath)
		fo, err := os.Create(gamePath)
		if err != nil {
			fmt.Printf("Unable to write file %s\n", gamePath)
			continue
		}
		if err := fo.Close(); err != nil {
			panic(err)
		}
		fmt.Printf("Closing Game file %s\n", gamePath)
	}
}
