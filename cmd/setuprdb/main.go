package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
			processConfig(dataConfig)
		}
		return
	}

	// Else try single
	dataConfig, ok := config.DataConfigs[configKey]
	if !ok {
		fmt.Println("Invalid DataConfig key")
		return
	}
	processConfig(dataConfig)

}

func processConfig(dataConfig config.DataConfig) {
	rdbPath := fetchRDB(dataConfig)
	if rdbPath != "" {
		makeNDJSON(dataConfig, rdbPath)
	}
}

func makeNDJSON(dataConfig config.DataConfig, rdbPath string) {
	dirPath := config.CommandRootPath
	toolPath := "./libretrodb_tool"

	corePath := filepath.Join(dirPath, "cores", dataConfig.ScrapeFolder)

	ndjsonPath := filepath.Join(corePath, "rdb.ndjson")
	// open the out file for writing
	outfile, err := os.Create(ndjsonPath)
	if err != nil {
		fmt.Println("Cannot create NDJSON", dataConfig.ScrapeFolder, ndjsonPath)
		return
	}
	defer outfile.Close()

	cmd := exec.Command(toolPath, rdbPath, "list")
	cmd.Stdout = outfile

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running libretro_tool", dataConfig.ScrapeFolder)
		fmt.Println(err)
		return
	}
}

func fetchRDB(dataConfig config.DataConfig) string {
	coreLabel := dataConfig.ScrapeFolder
	dirPath := "/mnt/c/Users/bossr/Code/MiSTer_Games_Data_Utils"

	// Make cores dir if not exist
	coresPath := filepath.Join(dirPath, "cores")
	err := os.MkdirAll(coresPath, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Unable to create path %v\n", coresPath))
	}
	fmt.Println("cores folder ready")

	if dataConfig.RdbName == "" {
		fmt.Printf("no RDB name for %v Skipping\n", coreLabel)
		return ""
	}
	fmt.Printf("Starting RDB Fetch %s\n", coreLabel)

	// make core dir if not exists
	corePath := filepath.Join(coresPath, dataConfig.ScrapeFolder)
	coreMkErr := os.MkdirAll(corePath, os.ModePerm)
	if coreMkErr != nil {
		panic(fmt.Sprintf("Unable to create path %v\n", corePath))
	}
	fmt.Printf("Core Path Created %s\n", corePath)

	// Make roms dir if not exists
	romsPath := filepath.Join(corePath, "roms")
	romPathErr := os.MkdirAll(romsPath, os.ModePerm)
	if romPathErr != nil {
		panic(fmt.Sprintf("Unable to create roms path %v\n", romsPath))
	}
	fmt.Println("cores/roms folder ready")

	// if rdb doesn't exist, download
	rdbPath := filepath.Join(corePath, "libretro.rdb")
	if _, err := os.Stat(rdbPath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Error openingT %s\n", rdbPath)
		// else fetch remote rdb file, save to filename in cores/{core}/{name}
		fmtFile := strings.Replace(url.QueryEscape(dataConfig.RdbName), "+", "%20", -1)
		url := fmt.Sprintf("%s%s", rdb.RootRdbUrl, fmtFile)
		fmt.Printf("Trying GET %s\n", url)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Unable to GET url %s, skipping core\n", url)
			return ""
		}
		fmt.Printf("GET StatusCode %v\n", resp.StatusCode)

		fmt.Printf("Trying ReadAll from GET %s\n", dataConfig.RdbName)
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Unable to Read RDB respone from GET url %s, skipping core\n", url)
			return ""
		}

		fmt.Printf("Saving RDB Bytes from GET %s\n", dataConfig.RdbName)
		fo, err := os.Create(rdbPath)
		if err != nil {
			fmt.Printf("Unable to write file %s\n", rdbPath)
			return ""
		}
		fo.Write(body)
		if err := fo.Close(); err != nil {
			panic(err)
		}
		fmt.Println("Saved RDB", dataConfig.RdbName)
	} else {
		fmt.Println("Existing RDB, Skipping", dataConfig.RdbName)
	}
	return rdbPath
}
