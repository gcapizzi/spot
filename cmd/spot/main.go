package main

import (
	"fmt"
	"os"

	"github.com/gcapizzi/spot/internal/playlist"
	"github.com/gcapizzi/spot/internal/spotify"
	"github.com/spf13/cobra"
)

func main() {
	var client spotify.Client

	var rootCmd = &cobra.Command{
		Use:   "spot",
		Short: "Spot is the Spotify swiss army knife.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			clientID := os.Getenv("SPOT_CLIENT_ID")
			clientSecret := os.Getenv("SPOT_CLIENT_SECRET")

			authURL, clientChannel, err := spotify.Authenticate(clientID, clientSecret)
			if err != nil {
				return err
			}
			fmt.Printf("Navigate to the following URL in your browser:\n\n%s\n\n", authURL)

			client = <-clientChannel
			fmt.Println("Authenticated successfully.")

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			playlistParser := playlist.NewParser(client)

			err := playlistParser.CreatePlaylistFromText("spot", os.Stdin)
			if err != nil {
				return err
			}

			fmt.Println(`Playlist "spot" created.`)

			return nil
		},
	}

	rootCmd.Execute()
}
