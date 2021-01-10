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
		Use:          "spot",
		Short:        "Spot is the Spotify swiss army knife.",
		SilenceUsage: true,
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
	}

	var playlistName string
	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "Create a playlist from a list of titles",
		RunE: func(cmd *cobra.Command, args []string) error {
			playlistParser := playlist.NewParser(client)

			err := playlistParser.CreatePlaylistFromText(playlistName, os.Stdin)
			if err != nil {
				return err
			}

			fmt.Printf("Playlist \"%s\" created.\n", playlistName)

			return nil
		},
	}
	parseCmd.Flags().StringVarP(&playlistName, "name", "n", "spot", "the playlist name")
	rootCmd.AddCommand(parseCmd)

	rootCmd.Execute()
}
