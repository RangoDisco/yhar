package importers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rangodisco/yhar/internal/api/services"
)

type MalojaImporter struct {
	name            string
	path            string
	scrobbleService *services.ScrobbleService
}

type malojaExport struct {
	Scrobbles []malojaScrobble `json:"scrobbles"`
}

type malojaScrobble struct {
	Time  int64       `json:"time"`
	Track malojaTrack `json:"track"`
}
type malojaTrack struct {
	Artists []string    `json:"artists"`
	Title   string      `json:"title"`
	Album   malojaAlbum `json:"album"`
}

type malojaAlbum struct {
	Artists []string `json:"artists"`
	Title   string   `json:"albumtitle"`
}

func NewMalojaImporter(scrobbleService *services.ScrobbleService) Importer {
	return &MalojaImporter{
		name:            "Maloja",
		path:            "import/maloja",
		scrobbleService: scrobbleService,
	}
}

func (i *MalojaImporter) Name() string {
	return i.name
}

func (i *MalojaImporter) Import(ctx context.Context) error {
	file, err := os.ReadFile(fmt.Sprintf("%s/%s", i.path, "maloja_export.json"))
	if err != nil {
		return fmt.Errorf("unable to open export file")
	}

	var data malojaExport
	err = json.Unmarshal(file, &data)
	if err != nil {
		return fmt.Errorf("unable to unmarshal file to struct: %w", err)
	}

	if len(data.Scrobbles) == 0 {
		return fmt.Errorf("no scrobbles found in export file")
	}

	//errChan := make(chan error, len(data.Scrobbles))
	var errs []error
	for _, scrobble := range data.Scrobbles {
		//go func() {
		//errChan <- i.importScrobble(ctx, scrobble)
		err = i.importScrobble(ctx, scrobble)
		if err != nil {
			fmt.Printf("unable to import scrobble: %s - %s : %s \n", scrobble.Track.Artists, scrobble.Track.Title, err)
			errs = append(errs, err)
		}
		//}()
	}

	//var errs []error
	//for range data.Scrobbles {
	//	err := <-errChan
	//	if err != nil {
	//		errs = append(errs, err)
	//	}
	//}

	return errors.Join(errs...)
}

func (i *MalojaImporter) importScrobble(ctx context.Context, scrobble malojaScrobble) error {
	uScrobble, err := i.parseToUnifiedScrobble(ctx, scrobble)
	if err != nil {
		return fmt.Errorf("unable to parse scrobble into unified scrobble: %w", err)
	}

	_, err = i.scrobbleService.HandleNewScrobble(ctx, uScrobble)
	if err != nil {
		return fmt.Errorf("unable to handle to new scrobble: %w", err)
	}

	return nil
}

func (i *MalojaImporter) parseToUnifiedScrobble(ctx context.Context, entry interface{}) (*services.UnifiedScrobbleEntry, error) {
	scrobble, ok := entry.(malojaScrobble)
	if !ok {
		return nil, fmt.Errorf("unable to parse into unified scrobble")
	}

	return &services.UnifiedScrobbleEntry{
		Title:      scrobble.Track.Title,
		ListenedAt: time.Unix(scrobble.Time, 0),
		Artist:     scrobble.Track.Artists[0],
		Album:      scrobble.Track.Album.Title,
	}, nil
}
