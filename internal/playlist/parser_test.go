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

	t.Run("CreatePlaylistFromTrackList", func(t *testing.T) {
		t.Run("when everything works", func(t *testing.T) {
			client := NewFakeClient()
			client.Tracks = map[string]playlist.Track{
				"one":   {ID: "1"},
				"two":   {ID: "2"},
				"three": {ID: "3"},
			}
			playlistParser := playlist.NewParser(client)

			playlistParser.CreatePlaylistFromTrackList(context.Background(), "playlist", strings.NewReader("one\ntwo\nthree"))

			g.Expect(client.Playlists).To(HaveLen(1))
			g.Expect(client.Playlists["playlist"]).To(Equal([]playlist.Track{
				{ID: "1"},
				{ID: "2"},
				{ID: "3"},
			}))
		})

		t.Run("when the text is empty", func(t *testing.T) {
			client := NewFakeClient()
			playlistParser := playlist.NewParser(client)

			err := playlistParser.CreatePlaylistFromTrackList(context.Background(), "playlist", strings.NewReader("\n\n\n"))

			g.Expect(err).To(MatchError(`no tracks found, playlist "playlist" not created`))
			g.Expect(client.Playlists).To(BeEmpty())
		})

		t.Run("when creating the playlist fails", func(t *testing.T) {
			client := NewFakeClient()
			client.Tracks = map[string]playlist.Track{
				"foo": {ID: "bar"},
			}
			client.CreatePlaylistErr = errors.New("create-playlist-error")

			playlistParser := playlist.NewParser(client)
			err := playlistParser.CreatePlaylistFromTrackList(context.Background(), "playlist", strings.NewReader("foo"))

			g.Expect(err).To(MatchError("create-playlist-error"))
		})

		t.Run("when creating the playlist fails", func(t *testing.T) {
			client := NewFakeClient()
			client.Tracks = map[string]playlist.Track{
				"foo": {ID: "bar"},
			}
			client.AddTrackToPlaylistErr = errors.New("add-track-error")

			playlistParser := playlist.NewParser(client)
			err := playlistParser.CreatePlaylistFromTrackList(context.Background(), "playlist", strings.NewReader("foo"))

			g.Expect(err).To(MatchError("add-track-error"))
		})
	})

	t.Run("CreatePlaylistFromAlbumList", func(t *testing.T) {
		t.Run("when everything works", func(t *testing.T) {
			client := NewFakeClient()
			client.Albums = map[string]playlist.Album{
				"A": {
					ID: "a",
					Tracks: []playlist.Track{
						{ID: "1"},
						{ID: "2"},
						{ID: "3"},
					},
				},
				"B": {
					ID: "b",
					Tracks: []playlist.Track{
						{ID: "4"},
						{ID: "5"},
						{ID: "6"},
					},
				},
			}
			playlistParser := playlist.NewParser(client)

			playlistParser.CreatePlaylistFromAlbumList(context.Background(), "playlist", strings.NewReader("A\nB"))

			g.Expect(client.Playlists).To(HaveLen(1))
			g.Expect(client.Playlists["playlist"]).To(Equal([]playlist.Track{
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
			playlistParser := playlist.NewParser(client)

			err := playlistParser.CreatePlaylistFromAlbumList(context.Background(), "playlist", strings.NewReader("\n\n\n"))

			g.Expect(err).To(MatchError(`no albums found, playlist "playlist" not created`))
			g.Expect(client.Playlists).To(BeEmpty())
		})

		t.Run("when creating the playlist fails", func(t *testing.T) {
			client := NewFakeClient()
			client.Albums = map[string]playlist.Album{
				"Foo": {
					ID: "foo",
					Tracks: []playlist.Track{
						{ID: "foo"},
					},
				},
			}
			client.CreatePlaylistErr = errors.New("create-playlist-error")

			playlistParser := playlist.NewParser(client)
			err := playlistParser.CreatePlaylistFromAlbumList(context.Background(), "playlist", strings.NewReader("Foo"))

			g.Expect(err).To(MatchError("create-playlist-error"))
		})

		t.Run("when creating the playlist fails", func(t *testing.T) {
			client := NewFakeClient()
			client.Albums = map[string]playlist.Album{
				"Foo": {
					ID: "foo",
					Tracks: []playlist.Track{
						{ID: "foo"},
					},
				},
			}
			client.AddTrackToPlaylistErr = errors.New("add-track-error")

			playlistParser := playlist.NewParser(client)
			err := playlistParser.CreatePlaylistFromAlbumList(context.Background(), "playlist", strings.NewReader("Foo"))

			g.Expect(err).To(MatchError("add-track-error"))
		})
	})
}

func NewFakeClient() *FakeClient {
	return &FakeClient{
		Playlists: map[string][]playlist.Track{},
	}
}

type FakeClient struct {
	Tracks                map[string]playlist.Track
	Albums                map[string]playlist.Album
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

func (c *FakeClient) FindAlbum(ctx context.Context, query string) (playlist.Album, error) {
	a, ok := c.Albums[query]
	if !ok {
		return playlist.Album{}, fmt.Errorf("cannot find album '%s'", query)
	}

	return a, nil
}

func (c *FakeClient) CreatePlaylist(ctx context.Context, name string) (playlist.Playlist, error) {
	if c.CreatePlaylistErr != nil {
		return playlist.Playlist{}, c.CreatePlaylistErr
	}

	c.Playlists[name] = []playlist.Track{}
	return playlist.Playlist{ID: name, Name: name}, nil
}

func (c *FakeClient) AddTracksToPlaylist(ctx context.Context, playlist playlist.Playlist, tracks []playlist.Track) error {
	if c.AddTrackToPlaylistErr != nil {
		return c.AddTrackToPlaylistErr
	}

	c.Playlists[playlist.ID] = append(c.Playlists[playlist.ID], tracks...)
	return nil
}
