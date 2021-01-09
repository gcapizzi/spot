package command

import (
	"bufio"
	"os"

	"github.com/gcapizzi/spot/internal/spotify"
)

type CreatePlaylistCommand struct {
	client spotify.Client
}

func NewCreatePlaylistCommand(client spotify.Client) CreatePlaylistCommand {
	return CreatePlaylistCommand{client: client}
}

func (c CreatePlaylistCommand) Run() error {
	playlist, err := c.client.CreatePlaylist("spot")
	if err != nil {
		return err
	}

	tracklistScanner := bufio.NewScanner(os.Stdin)
	for tracklistScanner.Scan() {
		trackQuery := tracklistScanner.Text()
		track, err := c.client.FindTrack(trackQuery)
		if err == nil {
			c.client.AddTrackToPlaylist(playlist, track)
		}
	}

	return nil
}
