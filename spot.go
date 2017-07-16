package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gcapizzi/spot/spotify"
)

func main() {
	clientId := os.Getenv("SPOT_CLIENT_ID")
	clientSecret := os.Getenv("SPOT_CLIENT_SECRET")

	authUrl, clientChannel, err := spotify.Authenticate(clientId, clientSecret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", authUrl)
	client := <-clientChannel

	userId, err := client.CurrentUserId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged in as ", userId)

	playlist, err := client.CreatePlaylist("spot")
	if err != nil {
		log.Fatal(err)
	}

	tracklistScanner := bufio.NewScanner(os.Stdin)
	for tracklistScanner.Scan() {
		trackQuery := tracklistScanner.Text()
		track, err := client.FindTrack(trackQuery)
		if err == nil {
			client.AddTrackToPlaylist(playlist, track)
		}
	}
}
