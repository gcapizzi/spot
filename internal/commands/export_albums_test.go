package commands_test

import (
	"context"
	"strings"
	"testing"

	spot "github.com/gcapizzi/spot/internal"
	"github.com/gcapizzi/spot/internal/commands"

	. "github.com/onsi/gomega"
)

func TestExportAlbums(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("when everything works", func(t *testing.T) {
		client := NewFakeClient()
		client.Albums = map[string]spot.Album{
			"A": {
				ID:  "a",
				URL: "http://open.spotify.com/album/a",
			},
			"B": {
				ID:  "b",
				URL: "http://open.spotify.com/album/b",
			},
		}
		cmd := commands.NewExportAlbums(client)

		var output strings.Builder
		cmd.Run(context.Background(), &output)

		g.Expect(output.String()).To(ContainSubstring("http://open.spotify.com/album/a"))
		g.Expect(output.String()).To(ContainSubstring("http://open.spotify.com/album/b"))
	})
}
