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

func TestParseAlbums(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("when everything works", func(t *testing.T) {
		client := NewFakeClient()
		client.Albums = map[string]spot.Album{
			"A": {
				ID: "a",
				Tracks: []spot.Track{
					{ID: "1"},
					{ID: "2"},
					{ID: "3"},
				},
			},
			"B": {
				ID: "b",
				Tracks: []spot.Track{
					{ID: "4"},
					{ID: "5"},
					{ID: "6"},
				},
			},
		}
		cmd := commands.NewParseAlbums(client)

		cmd.Run(context.Background(), "playlist", strings.NewReader("A\nB"))

		g.Expect(client.Playlists).To(HaveLen(1))
		g.Expect(client.Playlists["playlist"]).To(Equal([]spot.Track{
			{ID: "1"},
			{ID: "2"},
			{ID: "3"},
			{ID: "4"},
			{ID: "5"},
			{ID: "6"},
		}))
	})

	t.Run("when the text is empty", func(t *testing.T) {
		client := NewFakeClient()
		cmd := commands.NewParseAlbums(client)

		err := cmd.Run(context.Background(), "playlist", strings.NewReader("\n\n\n"))

		g.Expect(err).To(MatchError(`no albums found, playlist "playlist" not created`))
		g.Expect(client.Playlists).To(BeEmpty())
	})

	t.Run("when creating the playlist fails", func(t *testing.T) {
		client := NewFakeClient()
		client.Albums = map[string]spot.Album{
			"Foo": {
				ID: "foo",
				Tracks: []spot.Track{
					{ID: "foo"},
				},
			},
		}
		client.CreatePlaylistErr = errors.New("create-playlist-error")

		cmd := commands.NewParseAlbums(client)
		err := cmd.Run(context.Background(), "playlist", strings.NewReader("Foo"))

		g.Expect(err).To(MatchError("create-playlist-error"))
	})

	t.Run("when creating the playlist fails", func(t *testing.T) {
		client := NewFakeClient()
		client.Albums = map[string]spot.Album{
			"Foo": {
				ID: "foo",
				Tracks: []spot.Track{
					{ID: "foo"},
				},
			},
		}
		client.AddTrackToPlaylistErr = errors.New("add-track-error")

		cmd := commands.NewParseAlbums(client)
		err := cmd.Run(context.Background(), "playlist", strings.NewReader("Foo"))

		g.Expect(err).To(MatchError("add-track-error"))
	})
}
