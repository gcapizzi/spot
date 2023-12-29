package commands_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	spot "github.com/gcapizzi/spot/internal"
	"github.com/gcapizzi/spot/internal/commands"

	. "github.com/onsi/gomega"
)

func TestParseTracks(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("when everything works", func(t *testing.T) {
		client := NewFakeClient()
		client.Tracks = map[string]spot.Track{
			"one":   {ID: "1"},
			"two":   {ID: "2"},
			"three": {ID: "3"},
		}
		cmd := commands.NewParseTracks(client)

		cmd.Run(context.Background(), "playlist", strings.NewReader("one\ntwo\nthree"))

		g.Expect(client.Playlists).To(HaveLen(1))
		g.Expect(client.Playlists["playlist"]).To(Equal([]spot.Track{
			{ID: "1"},
			{ID: "2"},
			{ID: "3"},
		}))
	})

	t.Run("when the text is empty", func(t *testing.T) {
		client := NewFakeClient()
		cmd := commands.NewParseTracks(client)

		err := cmd.Run(context.Background(), "playlist", strings.NewReader("\n\n\n"))

		g.Expect(err).To(MatchError(`no tracks found, playlist "playlist" not created`))
		g.Expect(client.Playlists).To(BeEmpty())
	})

	t.Run("when creating the playlist fails", func(t *testing.T) {
		client := NewFakeClient()
		client.Tracks = map[string]spot.Track{
			"foo": {ID: "bar"},
		}
		client.CreatePlaylistErr = errors.New("create-playlist-error")

		cmd := commands.NewParseTracks(client)
		err := cmd.Run(context.Background(), "playlist", strings.NewReader("foo"))

		g.Expect(err).To(MatchError("create-playlist-error"))
	})

	t.Run("when creating the playlist fails", func(t *testing.T) {
		client := NewFakeClient()
		client.Tracks = map[string]spot.Track{
			"foo": {ID: "bar"},
		}
		client.AddTrackToPlaylistErr = errors.New("add-track-error")

		cmd := commands.NewParseTracks(client)
		err := cmd.Run(context.Background(), "playlist", strings.NewReader("foo"))

		g.Expect(err).To(MatchError("add-track-error"))
	})
}
