package rdb

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BossRighteous/MiSTer_Games_Data_Utils/pkg/utils"
)

type RdbJsonROM struct {
	Serial           string `json:"serial"`
	MD5              string `json:"md5"`
	SHA1             string `json:"sha1"`
	CRC              string `json:"crc"`
	Size             int    `json:"size"`
	RomName          string `json:"rom_name"`
	Region           string `json:"region"`
	Description      string `json:"description"`
	Name             string `json:"name"`
	Publisher        string `json:"publisher"`
	Developer        string `json:"developer"`
	ReleaseYear      int    `json:"releaseyear"`
	ReleaseMonth     int    `json:"releasemonth"`
	Users            int    `json:"users"`
	Genre            string `json:"genre"`
	Franchise        string `json:"franchise"`
	RDBID            int    `json:"rbd_id"`
	ParentExternalID string `json:"parent_external_id"`
	MGDBGameID       int    `json:"mgdb_game_id"`
}

func LoadNDJSON(corePath string) ([]RdbJsonROM, error) {
	rdbJsonPath := filepath.Join(corePath, "rdb.ndjson")
	defaultRoms := make([]RdbJsonROM, 0)
	fmt.Printf("Opening %s\n", rdbJsonPath)
	rdbFile, err := os.Open(rdbJsonPath)
	if err != nil {
		fmt.Printf("Error openingT %s\n", rdbJsonPath)
		return defaultRoms, err
	}
	fmt.Printf("Trying ReadAll from Local %s\n", rdbJsonPath)
	rdbBytes, err := io.ReadAll(rdbFile)
	if err != nil && err != io.EOF {
		fmt.Printf("Unable parse local bytes %s, skipping core\n", rdbJsonPath)
		return defaultRoms, err
	}
	return ParseNDJSON(rdbBytes)
}

func ParseNDJSON(jsonStream []byte) ([]RdbJsonROM, error) {
	// escape / char
	stringContent := strings.ReplaceAll(string(jsonStream[:]), "\\", "/")

	dec := json.NewDecoder(strings.NewReader(stringContent))
	rdbID := 1
	roms := make([]RdbJsonROM, 0)
	for {
		var jsonRom RdbJsonROM
		if err := dec.Decode(&jsonRom); err == io.EOF {
			break
		} else if err != nil {
			return roms, err
		}
		jsonRom.RDBID = rdbID
		rdbID++
		//fmt.Printf("%#v\n", jsonRom)
		roms = append(roms, jsonRom)
	}
	return roms, nil
}

// Reindex dedupe on slug
// replace map if full romname is shorter (favors USA)
// Target only one per permutation
func MapDupeROMSlugs(roms []RdbJsonROM) map[string]string {
	dupeMap := make(map[string]string) // slug:romName
	for _, rom := range roms {
		fileExt := filepath.Ext(rom.RomName)
		filename, _ := utils.CutSuffix(rom.RomName, fileExt)
		slug := utils.SlugifyString(filename)
		if existing, ok := dupeMap[slug]; ok {
			if len(rom.RomName) < len(existing) {
				dupeMap[slug] = rom.RomName
			}
		} else {
			dupeMap[slug] = rom.RomName
		}
	}
	return dupeMap
}
