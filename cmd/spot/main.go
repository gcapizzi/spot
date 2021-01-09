package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gcapizzi/spot/internal/playlist"
	"github.com/gcapizzi/spot/internal/spotify"
)

func main() {
	clientID := os.Getenv("SPOT_CLIENT_ID")
	clientSecret := os.Getenv("SPOT_CLIENT_SECRET")

	authURL, clientChannel, err := spotify.Authenticate(clientID, clientSecret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", authURL)

	client := <-clientChannel

	playlistParser := playlist.NewParser(client)

	err = playlistParser.CreatePlaylistFromText("spot", os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}
