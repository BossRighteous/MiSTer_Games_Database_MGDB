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

var TestData []byte = []byte(`
<gameList>
  <provider>
    <System>Nintendo 64</System>
    <software>Skraper</software>
    <database>ScreenScraper.fr</database>
    <web>http://www.screenscraper.fr</web>
  </provider>
  <game id="5426" source="ScreenScraper.fr">
    <path>./007 - The World Is Not Enough (USA) (v2) (Beta).v64</path>
    <name>007 : The World is Not Enough</name>
    <desc>The World Is Not Enough is a first-person shooter video game developed by Eurocom and based on the 1999 James Bond film of the same name. It was published by Electronic Arts and released for the Nintendo 64 on October 17, 2000, shortly before the release of its PlayStation counterpart. The game features a single-player campaign in which players assume the role of MI6 agent James Bond as he fights to stop a terrorist from triggering a nuclear meltdown in the waters of Istanbul. It includes a split-screen multiplayer mode where up to four players can compete in different types of deathmatch and objective-based games.</desc>
    <rating>0.65</rating>
    <releasedate>20001101T000000</releasedate>
    <developer>Eurocom Developments</developer>
    <publisher>Electronic Arts</publisher>
    <genre>Shooter / 1st person-Shooter-Action</genre>
    <players>1-4</players>
    <image>./media/screenshot/007 - The World Is Not Enough (USA) (v21) (Beta).png</image>
    <thumbnail>./media/screenshottitle/007 - The World Is Not Enough (USA) (v21) (Beta).png</thumbnail>
    <genreid>259</genreid>
  </game>
</gameList>
`)
