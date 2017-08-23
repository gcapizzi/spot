package command

import (
	"testing"

	"github.com/gcapizzi/spot/spotify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePlaylistCommand(t *testing.T) {
	assert := assert.New(t)

	client := new(FakeClient)
	command := NewCreatePlaylistCommand(client)

	playlist := spotify.Playlist{}
	client.On("CreatePlaylist", "spot").Return(playlist, nil)
	client.On("AddTrackToPlaylist").Return(nil)

	error := command.Run()

	assert.Nil(error)
}

type FakeClient struct {
	mock.Mock
}

func (client *FakeClient) CurrentUserId() (string, error) {
	returnValues := client.Called()
	return returnValues.String(0), returnValues.Error(1)
}

func (client *FakeClient) FindTrack(query string) (spotify.Track, error) {
	returnValues := client.Called(query)
	return returnValues.Get(0).(spotify.Track), returnValues.Error(1)
}

func (client *FakeClient) CreatePlaylist(name string) (spotify.Playlist, error) {
	returnValues := client.Called(name)
	return returnValues.Get(0).(spotify.Playlist), returnValues.Error(1)
}

func (client *FakeClient) AddTrackToPlaylist(playlist spotify.Playlist, track spotify.Track) error {
	returnValues := client.Called(playlist, track)
	return returnValues.Error(0)
}
