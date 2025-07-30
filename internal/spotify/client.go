package spotify

import (
	"context"
	"fmt"

	spot "github.com/gcapizzi/spot/internal"
	"github.com/zmb3/spotify/v2"
)

type Client struct {
	spotifyClient *spotify.Client
}

func (c Client) FindTrack(ctx context.Context, query string) (spot.Track, error) {
	results, err := c.spotifyClient.Search(ctx, query, spotify.SearchTypeTrack)
	if err != nil {
		return spot.Track{}, err
	}
	if len(results.Tracks.Tracks) == 0 {
		return spot.Track{}, fmt.Errorf("no track found for query %s", query)
	}

	firstResult := results.Tracks.Tracks[0]

	var artists []string
	for _, artist := range firstResult.Artists {
		artists = append(artists, artist.Name)
	}

	return spot.Track{
		ID:      string(firstResult.ID),
		Title:   firstResult.Name,
		Artists: artists,
		Album:   firstResult.Album.Name,
	}, nil
}

func (c Client) FindAlbum(ctx context.Context, query string) (spot.Album, error) {
	results, err := c.spotifyClient.Search(ctx, query, spotify.SearchTypeAlbum)
	if err != nil {
		return spot.Album{}, err
	}
	if len(results.Albums.Albums) == 0 {
		return spot.Album{}, fmt.Errorf("no album found for query %s", query)
	}

	firstResult := results.Albums.Albums[0]
	album, err := c.spotifyClient.GetAlbum(ctx, firstResult.ID)

	var artists []string
	for _, artist := range album.Artists {
		artists = append(artists, artist.Name)
	}

	var tracks []spot.Track
	for _, t := range album.Tracks.Tracks {
		tracks = append(tracks, spot.Track{
			ID:      string(t.ID),
			Title:   t.Name,
			Artists: []string{t.Artists[0].Name},
			Album:   t.Album.Name,
		})
	}

	return spot.Album{
		ID:      string(firstResult.ID),
		Title:   album.Name,
		Artists: artists,
		Tracks:  tracks,
	}, nil
}

func (c Client) CreatePlaylist(ctx context.Context, name string) (spot.Playlist, error) {
	user, err := c.spotifyClient.CurrentUser(ctx)
	if err != nil {
		return spot.Playlist{}, err
	}

	spotifyPlaylist, err := c.spotifyClient.CreatePlaylistForUser(ctx, user.ID, name, "", false, false)
	if err != nil {
		return spot.Playlist{}, err
	}

	return spot.Playlist{ID: string(spotifyPlaylist.ID), Name: name}, nil
}

func (c Client) AddTracksToPlaylist(ctx context.Context, playlist spot.Playlist, tracks []spot.Track) error {
	var trackIDs []spotify.ID
	for _, track := range tracks {
		trackIDs = append(trackIDs, spotify.ID(track.ID))
	}

	_, err := c.spotifyClient.AddTracksToPlaylist(
		ctx,
		spotify.ID(playlist.ID),
		trackIDs...,
	)

	return err
}

func (c Client) SavedAlbums(ctx context.Context) ([]spot.Album, error) {
	return nil, nil
}
