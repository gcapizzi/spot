package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"

	spot "github.com/gcapizzi/spot/internal"
)

type ParseTracks struct {
	client Client
}

func NewParseTracks(client Client) ParseTracks {
	return ParseTracks{client: client}
}

func (c ParseTracks) Run(ctx context.Context, name string, reader io.Reader) error {
	tracklistScanner := bufio.NewScanner(reader)

	var tracks []spot.Track
	for tracklistScanner.Scan() {
		trackQuery := tracklistScanner.Text()
		track, err := c.client.FindTrack(ctx, trackQuery)
		if err == nil {
			tracks = append(tracks, track)
		}
	}

	if len(tracks) == 0 {
		return fmt.Errorf(`no tracks found, playlist "%s" not created`, name)
	}

	playlist, err := c.client.CreatePlaylist(ctx, name)
	if err != nil {
		return err
	}

	return c.client.AddTracksToPlaylist(ctx, playlist, tracks)
}
