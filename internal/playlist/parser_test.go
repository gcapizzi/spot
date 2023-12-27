package playlist_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/gcapizzi/spot/internal/playlist"

	. "github.com/onsi/gomega"
)

func TestParser(t *testing.T) {
	g := NewGomegaWithT(t)

	client := NewFakeClient()
	client.Tracks = map[string]playlist.Track{
		"one":   {ID: "1"},
		"two":   {ID: "2"},
		"three": {ID: "3"},
	}

	playlistParser := playlist.NewParser(client)
	playlistParser.CreatePlaylistFromText(context.Background(), "playlist", strings.NewReader("one\ntwo\nthree"))

	g.Expect(client.Playlists).To(HaveLen(1))
	g.Expect(client.Playlists["playlist"]).To(Equal([]playlist.Track{
		{ID: "1"},
		{ID: "2"},
		{ID: "3"},
	}))
}

func TestParserWithEmptyText(t *testing.T) {
	g := NewGomegaWithT(t)

	client := NewFakeClient()

	playlistParser := playlist.NewParser(client)
	err := playlistParser.CreatePlaylistFromText(context.Background(), "playlist", strings.NewReader("\n\n\n"))

	g.Expect(err).To(MatchError(`no tracks found, playlist "playlist" not created`))
	g.Expect(client.Playlists).To(BeEmpty())
}

func TestParserFailingToCreatePlaylist(t *testing.T) {
	g := NewGomegaWithT(t)

	client := NewFakeClient()
	client.Tracks = map[string]playlist.Track{
		"foo": {ID: "bar"},
	}
	client.CreatePlaylistErr = errors.New("create-playlist-error")

	playlistParser := playlist.NewParser(client)
	err := playlistParser.CreatePlaylistFromText(context.Background(), "playlist", strings.NewReader("foo"))

	g.Expect(err).To(MatchError("create-playlist-error"))
}

func TestParserFailingToAddTrackToPlaylist(t *testing.T) {
	g := NewGomegaWithT(t)

	client := NewFakeClient()
	client.Tracks = map[string]playlist.Track{
		"foo": {ID: "bar"},
	}
	client.AddTrackToPlaylistErr = errors.New("add-track-error")

	playlistParser := playlist.NewParser(client)
	err := playlistParser.CreatePlaylistFromText(context.Background(), "playlist", strings.NewReader("foo"))

	g.Expect(err).To(MatchError("add-track-error"))
}

func NewFakeClient() *FakeClient {
	return &FakeClient{
		Playlists: map[string][]playlist.Track{},
	}
}

type FakeClient struct {
	Tracks                map[string]playlist.Track
	Playlists             map[string][]playlist.Track
	CreatePlaylistErr     error
	AddTrackToPlaylistErr error
}

func (c *FakeClient) FindTrack(ctx context.Context, query string) (playlist.Track, error) {
	t, ok := c.Tracks[query]
	if !ok {
		return playlist.Track{}, fmt.Errorf("cannot find track '%s'", query)
	}

	return t, nil
}

func (c *FakeClient) CreatePlaylist(ctx context.Context, name string) (playlist.Playlist, error) {
	if c.CreatePlaylistErr != nil {
		return playlist.Playlist{}, c.CreatePlaylistErr
	}

	c.Playlists[name] = []playlist.Track{}
	return playlist.Playlist{ID: name, Name: name}, nil
}

func (c *FakeClient) AddTrackToPlaylist(ctx context.Context, playlist playlist.Playlist, track playlist.Track) error {
	if c.AddTrackToPlaylistErr != nil {
		return c.AddTrackToPlaylistErr
	}

	c.Playlists[playlist.ID] = append(c.Playlists[playlist.ID], track)
	return nil
}
