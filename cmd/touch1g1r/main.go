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

	if configKey == "all" {
		for _, dataConfig := range config.DataConfigs {
			parseGamelistAnd1G1R(dataConfig)
		}
		return
	}

	dataConfig, ok := config.DataConfigs[configKey]
	if !ok {
		fmt.Println("Invalid DataConfig key")
		return
	}
	parseGamelistAnd1G1R(dataConfig)
}

func parseGamelistAnd1G1R(dataConfig config.DataConfig) {
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

	// Map ID -> shortest path
	idToPath := make(map[string]string)
	for _, game := range gl.Games {
		if game.ID == "" || game.Path == "" {
			continue
		}
		existing, ok := idToPath[game.ID]
		if !ok || len(game.Path) < len(existing) {
			idToPath[game.ID] = game.Path
		}
	}

	cores1g1rPath := filepath.Join(dirPath, "cores_1g1r")
	core1g1rPath := filepath.Join(cores1g1rPath, dataConfig.ScrapeFolder)

	if err := os.MkdirAll(core1g1rPath, 0755); err != nil {
		fmt.Printf("Unable to create directory %s: %v\n", core1g1rPath, err)
		return
	}

	for _, gamePath := range idToPath {
		destPath := filepath.Join(core1g1rPath, gamePath)
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			fmt.Printf("Unable to create directory %s: %v\n", destDir, err)
			continue
		}
		fmt.Printf("Creating 1G1R file %s\n", destPath)
		fo, err := os.Create(destPath)
		if err != nil {
			fmt.Printf("Unable to write file %s\n", destPath)
			continue
		}
		if err := fo.Close(); err != nil {
			panic(err)
		}
		fmt.Printf("Closing 1G1R file %s\n", destPath)
	}
}
