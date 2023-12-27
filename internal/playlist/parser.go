package playlist

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

type Client interface {
	FindTrack(context.Context, string) (Track, error)
	CreatePlaylist(context.Context, string) (Playlist, error)
	AddTrackToPlaylist(context.Context, Playlist, Track) error
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

func (p Parser) CreatePlaylistFromText(ctx context.Context, name string, reader io.Reader) error {
	tracklistScanner := bufio.NewScanner(reader)

	var tracks []Track
	for tracklistScanner.Scan() {
		trackQuery := tracklistScanner.Text()
		track, err := p.client.FindTrack(ctx, trackQuery)
		if err == nil {
			tracks = append(tracks, track)
		}
	}

	if len(tracks) == 0 {
		return fmt.Errorf(`no tracks found, playlist "%s" not created`, name)
	}

	playlist, err := p.client.CreatePlaylist(ctx, name)
	if err != nil {
		return err
	}

	for _, track := range tracks {
		err = p.client.AddTrackToPlaylist(ctx, playlist, track)
		if err != nil {
			return err
		}
	}

	return nil
}
