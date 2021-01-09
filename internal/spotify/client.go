package spotify

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/gcapizzi/spot/internal/playlist"
	"github.com/zmb3/spotify"
)

type Client struct {
	spotifyClient *spotify.Client
}

func Authenticate(clientID, clientSecret string) (string, chan Client, error) {
	state, err := generateRandomState()
	if err != nil {
		return "", nil, err
	}

	auth := spotify.NewAuthenticator(
		"http://localhost:8080",
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistModifyPrivate,
	)
	auth.SetAuthInfo(clientID, clientSecret)

	clientChannel := make(chan Client)
	closeChannel := make(chan bool)

	server := http.Server{Addr: ":8080", Handler: http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		token, err := auth.Token(state, request)
		if err != nil {
			http.Error(response, "Authentication failed", http.StatusForbidden)
			log.Fatal(err)
		}

		spotifyClient := auth.NewClient(token)
		client := Client{spotifyClient: &spotifyClient}
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

func (c Client) FindTrack(query string) (playlist.Track, error) {
	results, err := c.spotifyClient.Search(query, spotify.SearchTypeTrack)
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

func (c Client) CreatePlaylist(name string) (playlist.Playlist, error) {
	user, err := c.spotifyClient.CurrentUser()
	if err != nil {
		return playlist.Playlist{}, err
	}

	spotifyPlaylist, err := c.spotifyClient.CreatePlaylistForUser(user.ID, name, "", false)
	if err != nil {
		return playlist.Playlist{}, err
	}

	return playlist.Playlist{ID: string(spotifyPlaylist.ID), Name: name}, nil
}

func (c Client) AddTrackToPlaylist(playlist playlist.Playlist, track playlist.Track) error {
	_, err := c.spotifyClient.AddTracksToPlaylist(
		spotify.ID(playlist.ID),
		spotify.ID(track.ID),
	)

	return err
}
