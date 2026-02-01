package config

import (
	"github.com/rangodisco/yhar/internal/api/handlers"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/services"
	"gorm.io/gorm"
)

type Repositories struct {
	Scrobble repositories.IScrobbleRepository
	Album    repositories.IAlbumRepository
}

type Services struct {
	Scrobble services.IScrobbleService
}

type Handlers struct {
	Scrobble handlers.IScrobbleHandler
}

func AutoWire(db *gorm.DB) *Handlers {
	repos := &Repositories{
		Scrobble: repositories.NewScrobbleRepository(db),
	}

	services := &Services{
		Scrobble: services.NewScrobbleService(repos.Scrobble),
	}

	return &Handlers{
		Scrobble: handlers.NewScrobbleHandler(services.Scrobble),
	}
}
