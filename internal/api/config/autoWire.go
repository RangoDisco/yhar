package config

import (
	"github.com/rangodisco/yhar/internal/api/handlers"
	"github.com/rangodisco/yhar/internal/api/importers"
	"github.com/rangodisco/yhar/internal/api/pollers"
	"github.com/rangodisco/yhar/internal/api/providers"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/services"
	"gorm.io/gorm"
)

type Repositories struct {
	Album    *repositories.AlbumRepository
	Artist   *repositories.ArtistRepository
	Genre    *repositories.GenreRepository
	Image    *repositories.ImageRepository
	Scrobble *repositories.ScrobbleRepository
	Session  *repositories.SessionRepository
	Stats    *repositories.StatsRepository
	Track    *repositories.TrackRepository
	User     *repositories.UserRepository
}

type Services struct {
	Album         *services.AlbumService
	Artist        *services.ArtistService
	Auth          *services.AuthService
	Genre         *services.GenreService
	Image         *services.ImageService
	Metadata      *services.MetadataService
	Scrobble      *services.ScrobbleService
	ScrobbleStats *services.ScrobbleStatsService
	Session       *services.SessionService
	Track         *services.TrackService
	User          *services.UserService
}

type Handlers struct {
	Scrobble      *handlers.ScrobbleHandler
	ScrobbleStats *handlers.ScrobbleStatsHandler
	Auth          *handlers.AuthHandler
	User          *handlers.UserHandler
	Image         *handlers.ImageHandler
	Artist        *handlers.ArtistHandler
}

type Pollers struct {
	Subsonic pollers.PlayerPoller
}

type Importers struct {
	Maloja importers.Importer
}

func AutoWire(db *gorm.DB) (*Repositories, *Services, *Handlers, *Pollers, *Importers) {
	repos := &Repositories{
		Scrobble: repositories.NewScrobbleRepository(db),
		Album:    repositories.NewAlbumRepository(db),
		Artist:   repositories.NewArtistRepository(db),
		Genre:    repositories.NewGenreRepository(db),
		Image:    repositories.NewImageRepository(db),
		User:     repositories.NewUserRepository(db),
		Track:    repositories.NewTrackRepository(db),
		Stats:    repositories.NewStatsRepository(db),
		Session:  repositories.NewSessionRepository(db),
	}

	pvds := []providers.MetadataProvider{
		providers.NewMusicBrainzProvider(),
		providers.NewDeezerProvider(),
	}

	imageService := services.NewImageService(repos.Image)
	genreService := services.NewGenreService(repos.Genre)
	authService := services.NewAuthService(repos.User)
	albumService := services.NewAlbumService(repos.Album, imageService)
	artistService := services.NewArtistService(repos.Artist, imageService, genreService)
	metaService := services.NewMetadataService(pvds)
	trackService := services.NewTrackService(repos.Track)
	userService := services.NewUserService(repos.User)
	scrobbleStatsService := services.NewScrobbleStatsService(repos.Stats)
	scrobbleService := services.NewScrobbleService(repos.Scrobble, userService, trackService, artistService, albumService, metaService)
	sessionService := services.NewSessionService(repos.Session, scrobbleService)

	svs := &Services{
		Album:         albumService,
		Artist:        artistService,
		Auth:          authService,
		Genre:         genreService,
		Image:         imageService,
		Metadata:      metaService,
		Scrobble:      scrobbleService,
		ScrobbleStats: scrobbleStatsService,
		Session:       sessionService,
		Track:         trackService,
		User:          userService,
	}

	plrs := &Pollers{
		Subsonic: pollers.NewSubsonicPoller(svs.Session),
	}

	hdls := &Handlers{
		Auth:          handlers.NewAuthHandler(svs.Auth),
		Artist:        handlers.NewArtistHandler(svs.Artist),
		User:          handlers.NewUserHandler(svs.Auth),
		Scrobble:      handlers.NewScrobbleHandler(plrs.Subsonic),
		ScrobbleStats: handlers.NewScrobbleStatsHandler(svs.ScrobbleStats),
		Image:         handlers.NewImageHandler(svs.Image),
	}

	impts := &Importers{
		Maloja: importers.NewMalojaImporter(svs.Scrobble),
	}

	return repos, svs, hdls, plrs, impts
}
