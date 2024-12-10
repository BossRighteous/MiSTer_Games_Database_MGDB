package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/config"
	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/rdb"
)

func main() {
	fmt.Println("ran")

	dirPath := "/mnt/c/Users/bossr/Code/MiSTer_Games_Data_Utils"

	// Make cores dir if not exist
	coresPath := filepath.Join(dirPath, "cores")
	err := os.MkdirAll(coresPath, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Unable to create path %v\n", coresPath))
	}

	fmt.Println("cores folder ready")

	for coreLabel, config := range config.DataConfigs {
		if config.RdbName == "" {
			fmt.Printf("no RDB name for %v Skipping\n", coreLabel)
			continue
		}
		fmt.Printf("Starting core %s\n", coreLabel)

		// make core dir if not exists
		corePath := filepath.Join(coresPath, config.ScrapeFolder)
		err := os.MkdirAll(corePath, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("Unable to create path %v\n", corePath))
		}
		fmt.Printf("Core Path Created %s\n", corePath)

		// if rdb file exists, open
		rdbPath := filepath.Join(corePath, "libretro.rdb")
		fmt.Printf("Opening %s\n", rdbPath)
		rdbFile, err := os.Open(rdbPath)
		rdbBytes := make([]byte, 0)
		if err != nil {
			fmt.Printf("Error openingT %s\n", rdbPath)
			// else fetch remote rdb file, save to filename in cores/{core}/{name}
			fmtFile := strings.Replace(url.QueryEscape(config.RdbName), "+", "%20", -1)
			url := fmt.Sprintf("%s%s", rdb.RootRdbUrl, fmtFile)
			fmt.Printf("Trying GET %s\n", url)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Unable to GET url %s, skipping core\n", url)
				continue
			}
			fmt.Printf("GET StatusCode %v\n", resp.StatusCode)

			fmt.Printf("Trying ReadAll from GET %s\n", config.RdbName)
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Printf("Unable to Read RDB respone from GET url %s, skipping core\n", url)
				continue
			}
			rdbBytes = append(rdbBytes, body...)

			fmt.Printf("Saving RDB Bytes from GET %s\n", config.RdbName)
			fo, err := os.Create(rdbPath)
			if err != nil {
				fmt.Printf("Unable to write file %s\n", rdbPath)
				continue
			}
			fo.Write(rdbBytes)
			if err := fo.Close(); err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("Trying ReadAll from Local %s\n", config.RdbName)
			bytes, err := io.ReadAll(rdbFile)
			if err != nil && err != io.EOF {
				fmt.Printf("Unable parse local bytes %s, skipping core\n", config.RdbName)
				continue
			}
			rdbBytes = append(rdbBytes, bytes...)
		}
		fmt.Printf("rdbBytes appended %s\n", config.RdbName)

		if len(rdbBytes) == 0 {
			fmt.Printf("RDB bytes empty %s, skipping core\n", config.RdbName)
		}
		fmt.Printf("Parsing RDB File %s\n", config.RdbName)
		games := rdb.Parse(rdbBytes)
		fmt.Printf("%v Games read from %s\n", len(games), config.RdbName)

		// TODO: Refactor to only touch unique filename roots without extensions
		// Scraping only requires one instance and the RDB indexing can occur by similar match
		for _, game := range games {
			// Touch file, start empty
			gamePath := filepath.Join(corePath, game.ROMName)
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
}
