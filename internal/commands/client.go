package commands

import (
	"context"

	spot "github.com/gcapizzi/spot/internal"
)

type Client interface {
	FindTrack(context.Context, string) (spot.Track, error)
	FindAlbum(context.Context, string) (spot.Album, error)
	CreatePlaylist(context.Context, string) (spot.Playlist, error)
	AddTracksToPlaylist(context.Context, spot.Playlist, []spot.Track) error
	SavedAlbums(context.Context) ([]spot.Album, error)
}
