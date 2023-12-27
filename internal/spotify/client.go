package spotify

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/gcapizzi/spot/internal/playlist"
	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
)

type Client struct {
	spotifyClient *spotify.Client
}

func Authenticate(ctx context.Context, clientID, clientSecret string) (string, chan Client, error) {
	state, err := generateRandomState()
	if err != nil {
		return "", nil, err
	}

	auth := spotifyauth.New(spotifyauth.WithRedirectURL("http://localhost:8080"), spotifyauth.WithScopes(
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopePlaylistModifyPrivate,
	), spotifyauth.WithClientID(clientID), spotifyauth.WithClientSecret(clientSecret))

	clientChannel := make(chan Client)
	closeChannel := make(chan bool)

	server := http.Server{Addr: ":8080", Handler: http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		token, err := auth.Token(ctx, state, request)
		if err != nil {
			http.Error(response, "Authentication failed", http.StatusForbidden)
			log.Fatal(err)
		}

		spotifyClient := spotify.New(auth.Client(ctx, token))
		client := Client{spotifyClient: spotifyClient}
		clientChannel <- client
		closeChannel <- true
	})}

	go server.ListenAndServe()

	go func() {
		<-closeChannel
		server.Close()
	}()

	authURL := auth.AuthURL(state)

	return authURL, clientChannel, nil
}

func generateRandomState() (string, error) {
	randomInt, err := rand.Int(rand.Reader, big.NewInt(1000000000))
	state := randomInt.String()
	return state, err
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
