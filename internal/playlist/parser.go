package playlist

import (
	"bufio"
	"io"
)

type Client interface {
	FindTrack(string) (Track, error)
	CreatePlaylist(string) (Playlist, error)
	AddTrackToPlaylist(Playlist, Track) error
}

type Playlist struct {
	ID   string
	Name string
}

type Track struct {
	ID      string
	Title   string
	Artists []string
	Album   string
}

type Parser struct {
	client Client
}

func NewParser(client Client) Parser {
	return Parser{client: client}
}

func (p Parser) CreatePlaylistFromText(name string, reader io.Reader) error {
	playlist, err := p.client.CreatePlaylist(name)
	if err != nil {
		return err
	}

	tracklistScanner := bufio.NewScanner(reader)
	for tracklistScanner.Scan() {
		trackQuery := tracklistScanner.Text()
		track, err := p.client.FindTrack(trackQuery)
		if err == nil {
			p.client.AddTrackToPlaylist(playlist, track)
		}
	}

	return nil
}
