package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gcapizzi/spot/command"
	"github.com/gcapizzi/spot/spotify"
)

func main() {
	client := authenticate()

	userId, err := client.CurrentUserId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged in as ", userId)

	createPlaylistCommand := command.NewCreatePlaylistCommand(client)

	err = createPlaylistCommand.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func authenticate() spotify.Client {
	clientId := os.Getenv("SPOT_CLIENT_ID")
	clientSecret := os.Getenv("SPOT_CLIENT_SECRET")

	authUrl, clientChannel, err := spotify.Authenticate(clientId, clientSecret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", authUrl)
	client := <-clientChannel

	return client
}
