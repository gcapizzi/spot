package commands

import (
	"context"
	"fmt"
	"io"
)

type ExportAlbums struct {
	client Client
}

func NewExportAlbums(client Client) ExportAlbums {
	return ExportAlbums{client: client}
}

func (c ExportAlbums) Run(ctx context.Context, writer io.Writer) error {
	albums, _ := c.client.SavedAlbums(ctx)
	for _, album := range albums {
		fmt.Fprintln(writer, album.URL)
	}
	return nil
}
