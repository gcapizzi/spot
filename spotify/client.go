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
}

type Zmb3Client struct {
	spotifyClient *spotify.Client
}

type Track struct {
	Title   string
	Artists []string
	Album   string
}

func (client Zmb3Client) CurrentUserId() (string, error) {
	currentUser, err := client.spotifyClient.CurrentUser()
	if err != nil {
		return "", err
	}

	return currentUser.ID, nil
}

func Authenticate(clientId, clientSecret string) (string, chan Client) {
	state, _ := generateRandomState()

	auth := spotify.NewAuthenticator("http://localhost:8080", spotify.ScopeUserReadPrivate)
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

	return authUrl, clientChannel
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
	return Track{Title: firstResult.Name, Artists: []string{firstResult.Artists[0].Name}, Album: firstResult.Album.Name}, nil
}
