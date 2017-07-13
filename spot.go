package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gcapizzi/spot/spotify"
)

func main() {
	clientId := os.Getenv("SPOT_CLIENT_ID")
	clientSecret := os.Getenv("SPOT_CLIENT_SECRET")

	authUrl, clientChannel := spotify.Authenticate(clientId, clientSecret)

	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", authUrl)

	client := <-clientChannel

	userId, err := client.CurrentUserId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", userId)
}
