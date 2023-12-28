package playlist

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

type Client interface {
	FindTrack(context.Context, string) (Track, error)
	FindAlbum(context.Context, string) (Album, error)
	CreatePlaylist(context.Context, string) (Playlist, error)
	AddTracksToPlaylist(context.Context, Playlist, []Track) error
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

type Album struct {
	ID     string
	Tracks []Track
}

type Parser struct {
	client Client
}

func NewParser(client Client) Parser {
	return Parser{client: client}
}

func (p Parser) CreatePlaylistFromTrackList(ctx context.Context, name string, reader io.Reader) error {
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

	return p.client.AddTracksToPlaylist(ctx, playlist, tracks)
}

func (p Parser) CreatePlaylistFromAlbumList(ctx context.Context, name string, reader io.Reader) error {
	albumlistScanner := bufio.NewScanner(reader)

	var albums []Album
	for albumlistScanner.Scan() {
		albumQuery := albumlistScanner.Text()
		album, err := p.client.FindAlbum(ctx, albumQuery)
		if err == nil {
			albums = append(albums, album)
		}
	}

	if len(albums) == 0 {
		return fmt.Errorf(`no albums found, playlist "%s" not created`, name)
	}

	playlist, err := p.client.CreatePlaylist(ctx, name)
	if err != nil {
		return err
	}

	for _, album := range albums {
		err := p.client.AddTracksToPlaylist(ctx, playlist, album.Tracks)
		if err != nil {
			return err
		}
	}

	return nil
}
