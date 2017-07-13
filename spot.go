package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/zmb3/spotify"
)

func main() {
	clientId := os.Getenv("SPOT_CLIENT_ID")
	clientSecret := os.Getenv("SPOT_CLIENT_SECRET")

	authUrl, clientChannel := authenticate(clientId, clientSecret)

	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", authUrl)

	client := <-clientChannel

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
}

func authenticate(clientId, clientSecret string) (string, chan *spotify.Client) {
	state, _ := generateRandomState()

	auth := spotify.NewAuthenticator("http://localhost:8080", spotify.ScopeUserReadPrivate)
	auth.SetAuthInfo(clientId, clientSecret)

	clientChannel := make(chan *spotify.Client)

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		token, err := auth.Token(state, request)
		if err != nil {
			http.Error(response, "Authentication failed", http.StatusForbidden)
			log.Fatal(err)
		}

		client := auth.NewClient(token)
		clientChannel <- &client
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
