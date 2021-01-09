package playlist_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gcapizzi/spot/internal/playlist"

	. "github.com/onsi/gomega"
)

func TestParser(t *testing.T) {
	g := NewGomegaWithT(t)

	client := NewFakeClient(map[string]playlist.Track{
		"one":   {ID: "1"},
		"two":   {ID: "2"},
		"three": {ID: "3"},
	})

	playlistParser := playlist.NewParser(client)
	playlistParser.CreatePlaylistFromText("playlist", strings.NewReader("one\ntwo\nthree"))

	g.Expect(client.Playlists).To(HaveLen(1))
	g.Expect(client.Playlists["playlist"]).To(Equal([]playlist.Track{
		{ID: "1"},
		{ID: "2"},
		{ID: "3"},
	}))
}

func TestParserWithEmptyText(t *testing.T) {
	g := NewGomegaWithT(t)

	client := NewFakeClient(map[string]playlist.Track{})

	playlistParser := playlist.NewParser(client)
	err := playlistParser.CreatePlaylistFromText("playlist", strings.NewReader("\n\n\n"))

	g.Expect(err).To(MatchError(`no tracks found, playlist "playlist" not created`))
	g.Expect(client.Playlists).To(BeEmpty())
}

func NewFakeClient(tracks map[string]playlist.Track) *FakeClient {
	return &FakeClient{
		tracks:    tracks,
		Playlists: map[string][]playlist.Track{},
	}
}

type FakeClient struct {
	tracks    map[string]playlist.Track
	Playlists map[string][]playlist.Track
}

func (c *FakeClient) FindTrack(query string) (playlist.Track, error) {
	t, ok := c.tracks[query]
	if !ok {
		return playlist.Track{}, fmt.Errorf("cannot find track '%s'", query)
	}

	return t, nil
}

func (c *FakeClient) CreatePlaylist(name string) (playlist.Playlist, error) {
	c.Playlists[name] = []playlist.Track{}
	return playlist.Playlist{ID: name, Name: name}, nil
}

func (c *FakeClient) AddTrackToPlaylist(playlist playlist.Playlist, track playlist.Track) error {
	c.Playlists[playlist.ID] = append(c.Playlists[playlist.ID], track)
	return nil
}
