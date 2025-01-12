package gamelist

import (
	"encoding/xml"
)

type Gamelist struct {
	XMLName  xml.Name `xml:"gameList"`
	Provider Provider `xml:"provider"`
	Games    []Game   `xml:"game"`
}

type Provider struct {
	//XMLName  xml.Name `xml:"provider"`
	System           string `xml:"System"`
	Software         string `xml:"software"`
	Database         string `xml:"database"`
	Web              string `xml:"web"`
	MisterCollection string `xml:"misterguicollection"`
	MisterSystemIds  string `xml:"misterguisystemids"`
}

type Game struct {
	//XMLName     xml.Name `xml:"game"`
	ID          string `xml:"id,attr"`
	Source      string `xml:"source,attr"`
	Path        string `xml:"path"`
	Name        string `xml:"name"`
	Desc        string `xml:"desc"`
	Rating      string `xml:"rating"`
	ReleaseDate string `xml:"releasedate"`
	Developer   string `xml:"developer"`
	Publisher   string `xml:"publisher"`
	Genre       string `xml:"genre"`
	Players     string `xml:"players"`
	Image       string `xml:"image"`
	Thumbnail   string `xml:"thumbnail"`
	GenreID     string `xml:"genreid"`
}

func ParseGamelist(data []byte) *Gamelist {
	gamelist := &Gamelist{}
	err := xml.Unmarshal(data, gamelist)
	if err != nil {
		panic(err)
	}
	return gamelist
}
