package spotify

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/zmb3/spotify"
)

type Client interface {
	CurrentUserId() (string, error)
	FindTrack(string) (Track, error)
	CreatePlaylist(string) (Playlist, error)
	AddTrackToPlaylist(Playlist, Track) error
}

type Zmb3Client struct {
	spotifyClient *spotify.Client
}

type Track struct {
	Title     string
	Artists   []string
	Album     string
	spotifyId string
}

func (client Zmb3Client) CurrentUserId() (string, error) {
	currentUser, err := client.spotifyClient.CurrentUser()
	if err != nil {
		return "", err
	}

	return currentUser.ID, nil
}

func Authenticate(clientId, clientSecret string) (string, chan Client, error) {
	state, err := generateRandomState()
	if err != nil {
		return "", nil, err
	}

	auth := spotify.NewAuthenticator(
		"http://localhost:8080",
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistModifyPrivate,
	)
	auth.SetAuthInfo(clientId, clientSecret)

	clientChannel := make(chan Client)
	closeChannel := make(chan bool)

	server := http.Server{Addr: ":8080", Handler: http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		token, err := auth.Token(state, request)
		if err != nil {
			http.Error(response, "Authentication failed", http.StatusForbidden)
			log.Fatal(err)
		}

		spotifyClient := auth.NewClient(token)
		client := Zmb3Client{spotifyClient: &spotifyClient}
		clientChannel <- client
		closeChannel <- true
	})}

	go server.ListenAndServe()

	go func() {
		<-closeChannel
		server.Close()
	}()

	authUrl := auth.AuthURL(state)

	return authUrl, clientChannel, nil
}

func generateRandomState() (string, error) {
	randomInt, err := rand.Int(rand.Reader, big.NewInt(1000000000))
	state := randomInt.String()
	return state, err
}

func (client Zmb3Client) FindTrack(query string) (Track, error) {
	results, err := client.spotifyClient.Search(query, spotify.SearchTypeTrack)
	if err != nil {
		return Track{}, err
	}
	if len(results.Tracks.Tracks) == 0 {
		return Track{}, fmt.Errorf("No track found for query %s", query)
	}

	firstResult := results.Tracks.Tracks[0]
	track := Track{
		Title:     firstResult.Name,
		Artists:   []string{firstResult.Artists[0].Name},
		Album:     firstResult.Album.Name,
		spotifyId: string(firstResult.ID),
	}
	return track, nil
}

type Playlist struct {
	Name      string
	spotifyId string
}

func (client Zmb3Client) CreatePlaylist(name string) (Playlist, error) {
	userId, err := client.CurrentUserId()
	if err != nil {
		return Playlist{}, err
	}

	spotifyPlaylist, err := client.spotifyClient.CreatePlaylistForUser(userId, name, "", false)
	if err != nil {
		return Playlist{}, err
	}

	return Playlist{Name: name, spotifyId: string(spotifyPlaylist.ID)}, nil
}

func (client Zmb3Client) AddTrackToPlaylist(playlist Playlist, track Track) error {
	_, err := client.spotifyClient.AddTracksToPlaylist(
		spotify.ID(playlist.spotifyId),
		spotify.ID(track.spotifyId),
	)

	return err
}
