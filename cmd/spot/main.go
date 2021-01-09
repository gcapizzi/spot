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
	fmt.Printf("Navigate to the following URL in your browser:\n\n%s\n\n", authURL)

	client := <-clientChannel
	fmt.Println("Authenticated successfully.")

	playlistParser := playlist.NewParser(client)

	err = playlistParser.CreatePlaylistFromText("spot", os.Stdin)
	if err != nil {
		log.Fatalf("Error: %s.", err)
	}

	fmt.Println(`Playlist "spot" created.`)
}
