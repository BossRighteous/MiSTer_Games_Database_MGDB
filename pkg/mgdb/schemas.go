package mgdb

type MGDBInfo struct {
	CollectionName     string
	GamesFolder        string
	SupportedSystemIds string
	BuildDate          string
	MGDBVersion        string
	Description        string
}

type Game struct {
	GameID          int
	Name            string
	IsIndexed       int
	GenreId         int
	Rating          string
	ReleaseDate     string
	Developer       string
	Publisher       string
	Players         string
	Description     string
	ScreenshotHash  string
	TitleScreenHash string
}

type RDBRom struct {
	ROMName string
	Name    string
	CRC32   int
	Size    int
	Serial  string
	GameID  int
}

type GamelistRom struct {
	FileName           string
	GameID             int
	SupportedSystemIds string
	CRC32              int
}

type IndexedRom struct {
	Path               string
	FileName           string
	FileExt            string
	GameID             int
	SupportedSystemIds string
}

type Genre struct {
	GenreID int
	Name    string
}

type ImageBlob struct {
	Hash  string
	Bytes []byte
}
