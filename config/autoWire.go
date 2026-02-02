package config

import (
	"github.com/rangodisco/yhar/internal/api/handlers"
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
	Track    *repositories.TrackRepository
	User     *repositories.UserRepository
}

type Services struct {
	Album         *services.AlbumService
	Artist        *services.ArtistService
	Auth          *services.AuthService
	Genre         *services.GenreService
	Image         *services.ImageService
	Scrobble      *services.ScrobbleService
	ScrobbleStats *services.ScrobbleStatsService
	Track         *services.TrackService
	User          *services.UserService
}

type Handlers struct {
	Scrobble *handlers.ScrobbleHandler
	Auth     *handlers.AuthHandler
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
	}

	imageService := services.NewImageService(repos.Image)
	genreService := services.NewGenreService(repos.Genre)
	authService := services.NewAuthService(repos.User)
	albumService := services.NewAlbumService(repos.Album, imageService)
	artistService := services.NewArtistService(repos.Artist, imageService, genreService)
	trackService := services.NewTrackService(repos.Track)
	userService := services.NewUserService(repos.User)
	scrobbleStatsService := services.NewScrobbleStatsService(repos.Scrobble)
	scrobbleService := services.NewScrobbleService(repos.Scrobble, userService, trackService, artistService, albumService)

	svs := &Services{
		Album:         albumService,
		Artist:        artistService,
		Auth:          authService,
		Genre:         genreService,
		Image:         imageService,
		Scrobble:      scrobbleService,
		ScrobbleStats: scrobbleStatsService,
		Track:         trackService,
		User:          userService,
	}

	hdls := &Handlers{
		Scrobble: handlers.NewScrobbleHandler(svs.Scrobble, svs.ScrobbleStats),
		Auth:     handlers.NewAuthHandler(svs.Auth),
	}

	return repos, svs, hdls
}
