package main

import (
	"fmt"
	"os"

	"github.com/gcapizzi/spot/internal/commands"
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

			authURL, clientChannel, err := spotify.Authenticate(cmd.Context(), clientID, clientSecret)
			if err != nil {
				return err
			}
			fmt.Printf("👉 %s\n", authURL)

			client = <-clientChannel
			fmt.Println("✅ Authenticated successfully.")

			return nil
		},
	}

	var playlistName string
	parseTracksCmd := &cobra.Command{
		Use:   "parse-tracks",
		Short: "Create a playlist from a list of track titles",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseTracksCmd := commands.NewParseTracks(client)

			err := parseTracksCmd.Run(cmd.Context(), playlistName, os.Stdin)
			if err != nil {
				return err
			}

			fmt.Printf("Playlist \"%s\" created.\n", playlistName)

			return nil
		},
	}
	parseTracksCmd.Flags().StringVarP(&playlistName, "name", "n", "spot", "the playlist name")
	rootCmd.AddCommand(parseTracksCmd)

	parseAlbumsCmd := &cobra.Command{
		Use:   "parse-albums",
		Short: "Create a playlist from a list of album titles",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseAlbumsCmd := commands.NewParseAlbums(client)

			err := parseAlbumsCmd.Run(cmd.Context(), playlistName, os.Stdin)
			if err != nil {
				return err
			}

			fmt.Printf("Playlist \"%s\" created.\n", playlistName)

			return nil
		},
	}
	parseAlbumsCmd.Flags().StringVarP(&playlistName, "name", "n", "spot", "the playlist name")
	rootCmd.AddCommand(parseAlbumsCmd)

	rootCmd.Execute()
}
