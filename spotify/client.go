package spotify

import (
	"crypto/rand"
	"log"
	"math/big"
	"net/http"

	"github.com/zmb3/spotify"
)

type Client interface {
	CurrentUserId() (string, error)
}

type Zmb3Client struct {
	spotifyClient *spotify.Client
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

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		token, err := auth.Token(state, request)
		if err != nil {
			http.Error(response, "Authentication failed", http.StatusForbidden)
			log.Fatal(err)
		}

		spotifyClient := auth.NewClient(token)
		client := Zmb3Client{spotifyClient: &spotifyClient}
		clientChannel <- client
	})
	go http.ListenAndServe(":8080", nil)

	authUrl := auth.AuthURL(state)

	return authUrl, clientChannel
}

func generateRandomState() (string, error) {
	randomInt, err := rand.Int(rand.Reader, big.NewInt(1000000000))
	state := randomInt.String()
	return state, err
}
