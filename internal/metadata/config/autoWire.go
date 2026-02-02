package config

import (
	"github.com/rangodisco/yhar/internal/metadata/repositories"
	"github.com/rangodisco/yhar/internal/metadata/services"
	"gorm.io/gorm"
)

type Repositories struct {
	Album *repositories.AlbumRepository
	Track *repositories.TrackRepository
}

type Services struct {
	Album    *services.AlbumService
	Artist   *services.ArtistService
	Scrobble *services.ScrobbleService
}

func AutoWire(db *gorm.DB) *Services {
	repos := &Repositories{
		Album: repositories.NewAlbumRepository(db),
		Track: repositories.NewTrackRepository(db),
	}

	artistService := services.NewArtistService()
	albumService := services.NewAlbumService(artistService)
	scrobbleService := services.NewScrobbleService(repos.Track, repos.Album, artistService, albumService)

	svs := &Services{
		Album:    albumService,
		Artist:   artistService,
		Scrobble: scrobbleService,
	}

	return svs
}
