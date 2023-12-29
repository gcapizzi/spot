package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"

	spot "github.com/gcapizzi/spot/internal"
)

type ParseAlbums struct {
	client Client
}

func NewParseAlbums(client Client) ParseAlbums {
	return ParseAlbums{client: client}
}

func (c ParseAlbums) Run(ctx context.Context, name string, reader io.Reader) error {
	albumlistScanner := bufio.NewScanner(reader)

	var albums []spot.Album
	for albumlistScanner.Scan() {
		albumQuery := albumlistScanner.Text()
		album, err := c.client.FindAlbum(ctx, albumQuery)
		if err == nil {
			albums = append(albums, album)
		}
	}

	if len(albums) == 0 {
		return fmt.Errorf(`no albums found, playlist "%s" not created`, name)
	}

	playlist, err := c.client.CreatePlaylist(ctx, name)
	if err != nil {
		return err
	}

	for _, album := range albums {
		err := c.client.AddTracksToPlaylist(ctx, playlist, album.Tracks)
		if err != nil {
			return err
		}
	}

	return nil
}
