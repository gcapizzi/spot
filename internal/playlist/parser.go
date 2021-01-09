package playlist

import (
	"bufio"
	"fmt"
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
	tracklistScanner := bufio.NewScanner(reader)

	var tracks []Track
	for tracklistScanner.Scan() {
		trackQuery := tracklistScanner.Text()
		track, err := p.client.FindTrack(trackQuery)
		if err == nil {
			tracks = append(tracks, track)
		}
	}

	if len(tracks) == 0 {
		return fmt.Errorf(`no tracks found, playlist "%s" not created`, name)
	}

	playlist, err := p.client.CreatePlaylist(name)
	if err != nil {
		return err
	}

	for _, track := range tracks {
		err = p.client.AddTrackToPlaylist(playlist, track)
		if err != nil {
			return err
		}
	}

	return nil
}
