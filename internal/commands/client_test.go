package commands_test

import (
	"context"
	"fmt"
	"maps"
	"slices"

	spot "github.com/gcapizzi/spot/internal"
)

func NewFakeClient() *FakeClient {
	return &FakeClient{
		Playlists: map[string][]spot.Track{},
	}
}

type FakeClient struct {
	Tracks                map[string]spot.Track
	Albums                map[string]spot.Album
	Playlists             map[string][]spot.Track
	CreatePlaylistErr     error
	AddTrackToPlaylistErr error
}

func (c *FakeClient) FindTrack(ctx context.Context, query string) (spot.Track, error) {
	t, ok := c.Tracks[query]
	if !ok {
		return spot.Track{}, fmt.Errorf("cannot find track '%s'", query)
	}

	return t, nil
}

func (c *FakeClient) FindAlbum(ctx context.Context, query string) (spot.Album, error) {
	a, ok := c.Albums[query]
	if !ok {
		return spot.Album{}, fmt.Errorf("cannot find album '%s'", query)
	}

	return a, nil
}

func (c *FakeClient) CreatePlaylist(ctx context.Context, name string) (spot.Playlist, error) {
	if c.CreatePlaylistErr != nil {
		return spot.Playlist{}, c.CreatePlaylistErr
	}

	c.Playlists[name] = []spot.Track{}
	return spot.Playlist{ID: name, Name: name}, nil
}

func (c *FakeClient) AddTracksToPlaylist(ctx context.Context, playlist spot.Playlist, tracks []spot.Track) error {
	if c.AddTrackToPlaylistErr != nil {
		return c.AddTrackToPlaylistErr
	}

	c.Playlists[playlist.ID] = append(c.Playlists[playlist.ID], tracks...)
	return nil
}

func (c *FakeClient) SavedAlbums(ctx context.Context) ([]spot.Album, error) {
	return slices.Collect(maps.Values(c.Albums)), nil
}
