package pollers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/services"
	"github.com/rangodisco/yhar/internal/api/types/subsonic"
)

type SubsonicPoller struct {
	name           string
	baseUrl        string
	username       string
	password       string
	sessionService *services.SessionService
	sessionRepo    *repositories.SessionRepository
}

// Used for salt generation
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NewSubsonicPoller(sessionService *services.SessionService, sessionRepo *repositories.SessionRepository) PlayerPoller {
	baseUrl := os.Getenv("SUBSONIC_BASE_URL")
	username := os.Getenv("SUBSONIC_USER")
	password := os.Getenv("SUBSONIC_PASSWORD")
	return &SubsonicPoller{
		name:           "subsonic",
		baseUrl:        baseUrl,
		username:       username,
		sessionService: sessionService,
		password:       password,
		sessionRepo:    sessionRepo,
	}
}

func (p *SubsonicPoller) Name() string {
	return p.name
}

func (p *SubsonicPoller) Start(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := p.PollPlaying(ctx)
			if err != nil {
				slog.Error("unable to poll subsonic", slog.String("error", err.Error()))
			}
		case <-ctx.Done():
			slog.Debug("Subsonic poller stopped")
			return
		}
	}
}

// PollPlaying calls the subsonic server, and for each entry being played, gets or creates the session from/in db, checks if the session is completed, and creates a new scrobble if so
func (p *SubsonicPoller) PollPlaying(ctx context.Context) error {
	var nowPlayingRes subsonic.GetNowPlayingResponse

	// As per Subsonic's docs, for each REST call, generate a random string called the salt. Send this as parameter s.
	//Use a salt length of at least six characters.
	salt := p.generateSalt()
	hash := md5.Sum([]byte(fmt.Sprintf("%s%s", p.password, salt)))
	token := hex.EncodeToString(hash[:])

	pollUrl := fmt.Sprintf("%s/rest/getNowPlaying?u=%s&c=yhar&s=%s&t=%s&v=%s", p.baseUrl, p.username, salt, token, "1.16.1")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pollUrl, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to send request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			slog.Error("unable to close body", slog.String("error", err.Error()))
		}
	}(res.Body)

	err = xml.NewDecoder(res.Body).Decode(&nowPlayingRes)
	if err != nil {
		return fmt.Errorf("unable to decode response: %w", err)
	}

	if nowPlayingRes.Status == "failed" {
		return fmt.Errorf("subsonic auth error: %s", nowPlayingRes.Error.Message)
	}

	// No need to go further if nothing is being played
	if len(nowPlayingRes.NowPlaying.Entry) == 0 {
		slog.Debug("polled subsonic, no track were currently played")
		return nil
	}

	errChan := make(chan error, len(nowPlayingRes.NowPlaying.Entry))
	for _, e := range nowPlayingRes.NowPlaying.Entry {
		slog.Debug("polled subsonic, found track being played",
			slog.String("track", e.Title),
			slog.String("album", e.Album),
			slog.String("artist", e.Artist))
		go func() {
			errChan <- p.handleEntry(ctx, e)
		}()
	}

	var errs []error
	for range nowPlayingRes.NowPlaying.Entry {
		err := <-errChan
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// handleEntry gets or create the session from/in db, checks if the session is completed, and creates a new scrobble if so
func (p *SubsonicPoller) handleEntry(ctx context.Context, entry subsonic.Entry) error {
	session, err := p.sessionService.GetOrCreateSession(ctx, entry)
	if err != nil {
		return fmt.Errorf("unable to get session for entry: %s : %w", entry.Title, err)
	}

	// The session could be nil in case we're in its last 30% of the track, which we consider as "completed" and don't want to start a new session for it
	if session == nil {
		return nil
	}

	if p.isCompleted(session) {
		scrobble, err := p.parseToUnifiedScrobble(entry)
		if err != nil {
			return fmt.Errorf("unable to parse into unified scrobble: %w", err)
		}
		err = p.sessionService.HandleCompletedSession(ctx, session, scrobble)
		if err != nil {
			return fmt.Errorf("unable to complete session for entry: %s : %w", entry.Title, err)
		}
	} else {
		session.LastSeenAt = time.Now()
		updates := map[string]interface{}{"last_seen_at": time.Now()}
		err := p.sessionRepo.Update(ctx, session.ID, updates)
		if err != nil {
			return fmt.Errorf("unable to update session for entry: %s : %w", entry.Title, err)
		}
	}

	return nil
}

func (p *SubsonicPoller) parseToUnifiedScrobble(entry subsonic.Entry) (*services.UnifiedScrobbleEntry, error) {
	return &services.UnifiedScrobbleEntry{
		Username:      entry.Username,
		Title:         entry.Title,
		MusicBrainzID: entry.MusicBrainzID,
		Album:         entry.Album,
		Artist:        entry.Artist,
		Duration:      entry.Duration,
	}, nil

}

func (p *SubsonicPoller) generateSalt() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// isCompleted checks if at least 80% of the track has been played
func (p *SubsonicPoller) isCompleted(session *models.Session) bool {
	elapsed := time.Since(session.StartedAt)
	threshold := time.Duration(float64(session.Duration) * 0.8)
	return elapsed >= threshold
}
