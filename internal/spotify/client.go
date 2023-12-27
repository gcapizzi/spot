package spotify

import (
	"context"
	"fmt"

	"github.com/gcapizzi/spot/internal/playlist"
	"github.com/zmb3/spotify/v2"
)

type Client struct {
	spotifyClient *spotify.Client
}

func (c Client) FindTrack(ctx context.Context, query string) (playlist.Track, error) {
	results, err := c.spotifyClient.Search(ctx, query, spotify.SearchTypeTrack)
	if err != nil {
		return playlist.Track{}, err
	}
	if len(results.Tracks.Tracks) == 0 {
		return playlist.Track{}, fmt.Errorf("no track found for query %s", query)
	}

	firstResult := results.Tracks.Tracks[0]
	return playlist.Track{
		ID:      string(firstResult.ID),
		Title:   firstResult.Name,
		Artists: []string{firstResult.Artists[0].Name},
		Album:   firstResult.Album.Name,
	}, nil
}

func (c Client) CreatePlaylist(ctx context.Context, name string) (playlist.Playlist, error) {
	user, err := c.spotifyClient.CurrentUser(ctx)
	if err != nil {
		return playlist.Playlist{}, err
	}

	spotifyPlaylist, err := c.spotifyClient.CreatePlaylistForUser(ctx, user.ID, name, "", false, false)
	if err != nil {
		return playlist.Playlist{}, err
	}

	return playlist.Playlist{ID: string(spotifyPlaylist.ID), Name: name}, nil
}

func (c Client) AddTrackToPlaylist(ctx context.Context, playlist playlist.Playlist, track playlist.Track) error {
	_, err := c.spotifyClient.AddTracksToPlaylist(
		ctx,
		spotify.ID(playlist.ID),
		spotify.ID(track.ID),
	)

	return err
}
