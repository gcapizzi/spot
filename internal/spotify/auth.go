package spotify

import (
	"context"
	"crypto/rand"
	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"log"
	"math/big"
	"net/http"
)

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
