package config

import (
	"github.com/rangodisco/yhar/internal/api/handlers"
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
}

type Pollers struct {
	Subsonic pollers.PlayerPoller
}

func AutoWire(db *gorm.DB) (*Repositories, *Services, *Handlers) {
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

	svs := &Services{
		Album:         albumService,
		Artist:        artistService,
		Auth:          authService,
		Genre:         genreService,
		Image:         imageService,
		Metadata:      metaService,
		Scrobble:      scrobbleService,
		ScrobbleStats: scrobbleStatsService,
		Track:         trackService,
		User:          userService,
	}

	plrs := &Pollers{
		Subsonic: pollers.NewSubsonicPoller(svs.Session),
	}

	hdls := &Handlers{
		Auth:          handlers.NewAuthHandler(svs.Auth),
		User:          handlers.NewUserHandler(svs.Auth),
		Scrobble:      handlers.NewScrobbleHandler(plrs.Subsonic),
		ScrobbleStats: handlers.NewScrobbleStatsHandler(svs.ScrobbleStats),
	}

	return repos, svs, hdls
}
