package subsonic

import "encoding/xml"

type GetNowPlayingResponse struct {
	XMLName       xml.Name   `xml:"subsonic-response"`
	Status        string     `xml:"status,attr"`
	Version       string     `xml:"version,attr"`
	Type          string     `xml:"type,attr"`
	ServerVersion string     `xml:"serverVersion,attr"`
	OpenSubsonic  string     `xml:"openSubsonic,attr"`
	NowPlaying    NowPlaying `xml:"nowPlaying"`
}

type NowPlaying struct {
	Entry []Entry `xml:"entry"`
}

type Entry struct {
	ID                 string        `xml:"id,attr"`
	Parent             string        `xml:"parent,attr"`
	IsDir              string        `xml:"isDir,attr"`
	Title              string        `xml:"title,attr"`
	Album              string        `xml:"album,attr"`
	Artist             string        `xml:"artist,attr"`
	Track              string        `xml:"track,attr"`
	Year               string        `xml:"year,attr"`
	Genre              string        `xml:"genre,attr"`
	CoverArt           string        `xml:"coverArt,attr"`
	Size               string        `xml:"size,attr"`
	ContentType        string        `xml:"contentType,attr"`
	Suffix             string        `xml:"suffix,attr"`
	Duration           string        `xml:"duration,attr"`
	BitRate            string        `xml:"bitRate,attr"`
	Path               string        `xml:"path,attr"`
	DiscNumber         string        `xml:"discNumber,attr"`
	Created            string        `xml:"created,attr"`
	AlbumID            string        `xml:"albumId,attr"`
	ArtistID           string        `xml:"artistId,attr"`
	Type               string        `xml:"type,attr"`
	IsVideo            string        `xml:"isVideo,attr"`
	Comment            string        `xml:"comment,attr,omitempty"`
	SortName           string        `xml:"sortName,attr"`
	MediaType          string        `xml:"mediaType,attr"`
	MusicBrainzID      string        `xml:"musicBrainzId,attr,omitempty"`
	ChannelCount       string        `xml:"channelCount,attr"`
	SamplingRate       string        `xml:"samplingRate,attr"`
	BitDepth           string        `xml:"bitDepth,attr"`
	DisplayArtist      string        `xml:"displayArtist,attr"`
	DisplayAlbumArtist string        `xml:"displayAlbumArtist,attr"`
	Username           string        `xml:"username,attr"`
	MinutesAgo         string        `xml:"minutesAgo,attr"`
	PlayerID           string        `xml:"playerId,attr"`
	PlayerName         string        `xml:"playerName,attr"`
	PlayCount          string        `xml:"playCount,attr,omitempty"`
	Genres             []genre       `xml:"genres"`
	Artists            []artist      `xml:"artists"`
	AlbumArtists       []albumArtist `xml:"albumArtists"`
}

type artist struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type albumArtist struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type genre struct {
	Name string `xml:"name,attr"`
}
