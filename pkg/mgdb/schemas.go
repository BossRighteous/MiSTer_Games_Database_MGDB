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

type SlugRom struct {
	Slug               string
	GameID             int
	SupportedSystemIds string
}

type RomCrc struct {
	CRC32 string
	Slug  string
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
