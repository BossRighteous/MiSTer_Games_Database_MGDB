package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/gamelist"
)

func main() {
	cliArgs := os.Args
	if len(cliArgs) < 2 {
		fmt.Println("Usage: auditscrapes <cores|cores_1g1r>")
		return
	}

	rootDir := cliArgs[1]
	rootPath := filepath.Join("/mnt/c/Users/bossr/Code/MiSTer_Games_Data_Utils", rootDir)

	systemEntries, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Printf("Unable to read directory %s: %v\n", rootPath, err)
		return
	}

	for _, systemEntry := range systemEntries {
		if !systemEntry.IsDir() {
			continue
		}

		systemName := systemEntry.Name()
		systemPath := filepath.Join(rootPath, systemName)

		auditSystem(systemName, systemPath)
	}
}

func auditSystem(systemName, systemPath string) {
	fileMap := map[string]interface{}{}

	entries, err := os.ReadDir(systemPath)
	if err != nil {
		fmt.Printf("[%s] Unable to read directory: %v\n", systemName, err)
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() && name == "media" {
			continue
		}
		if entry.IsDir() {
			continue
		}
		if name == "gamelist.xml" {
			continue
		}
		ext := filepath.Ext(name)
		if ext == ".rdb" || ext == ".mgdb" {
			continue
		}
		fileMap[name] = struct{}{}
	}

	xmlPath := filepath.Join(systemPath, "gamelist.xml")
	data, err := os.ReadFile(xmlPath)
	if err != nil {
		fmt.Printf("\n=== %s ===\n", systemName)
		fmt.Println("WARNING: no gamelist.xml found")
		return
	}

	gl := gamelist.ParseGamelist(data)

	found := 0
	for _, game := range gl.Games {
		if game.Path == "" {
			continue
		}
		filename := strings.TrimPrefix(game.Path, "./")
		if _, exists := fileMap[filename]; exists {
			delete(fileMap, filename)
			found++
		}
	}

	notFound := len(fileMap)

	fmt.Printf("\n=== %s ===\n", systemName)
	if notFound > 0 {
		fmt.Println("Not in gamelist:")
		for filename := range fileMap {
			fmt.Println(" ", filename)
		}
	}
	fmt.Printf("Found: %d | Not in gamelist: %d\n", found, notFound)
}
