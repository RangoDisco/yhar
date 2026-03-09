package pollers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/services"
	"github.com/rangodisco/yhar/internal/api/types/subsonic"
)

type SubsonicPoller struct {
	name           string
	baseUrl        string
	username       string
	token          string
	salt           string
	sessionService *services.SessionService
}

func NewSubsonicPoller(sessionService *services.SessionService) PlayerPoller {
	baseUrl := os.Getenv("SUBSONIC_BASE_URL")
	username := os.Getenv("SUBSONIC_USER")
	password := os.Getenv("SUBSONIC_PASSWORD")
	salt := os.Getenv("SUBSONIC_SALT")
	hash := md5.Sum([]byte(fmt.Sprintf("%s+%s", password, salt)))
	token := hex.EncodeToString(hash[:])
	return &SubsonicPoller{
		name:           "subsonic",
		baseUrl:        baseUrl,
		username:       username,
		token:          token,
		salt:           salt,
		sessionService: sessionService,
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
			// TODO: handle chan
			p.PollPlaying(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// PollPlaying calls the subsonic server, and for each entry being played, gets or creates the session from/in db, checks if the session is completed, and creates a new scrobble if so
func (p *SubsonicPoller) PollPlaying(ctx context.Context) error {
	var nowPlayingRes subsonic.GetNowPlayingResponse
	pollUrl := fmt.Sprintf("%s/rest/getNowPlaying?u=%s&c=yhar&s=%s&t=%s", p.baseUrl, p.username, p.salt, p.token)
	req, err := http.NewRequest(http.MethodGet, pollUrl, nil)
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
			fmt.Println(err)
		}
	}(res.Body)

	err = xml.NewDecoder(res.Body).Decode(&nowPlayingRes)
	if err != nil {
		return fmt.Errorf("unable to decode response: %w", err)
	}

	// No need to go further if nothing is being played
	if len(nowPlayingRes.NowPlaying.Entry) == 0 {
		return nil
	}

	errChan := make(chan error, len(nowPlayingRes.NowPlaying.Entry))
	for _, entry := range nowPlayingRes.NowPlaying.Entry {
		go p.handleEntry(ctx, entry, errChan)
	}

	var errs []error
	for range nowPlayingRes.NowPlaying.Entry {
		err = <-errChan
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// handleEntry gets or create the session from/in db, checks if the session is completed, and creates a new scrobble if so
func (p *SubsonicPoller) handleEntry(ctx context.Context, entry subsonic.Entry, errChan chan<- error) {
	session, err := p.sessionService.GetOrCreateSession(ctx, entry)
	if err != nil {
		errChan <- fmt.Errorf("unable to get session for entry: %s : %w", entry.Title, err)
		return
	}

	// The session could be nil in case we're in its last 30% of the track, which we consider as "completed" and don't want to start a new session for it
	if session == nil {
		errChan <- nil
		return
	}

	if p.isCompleted(session) {
		err = p.sessionService.HandleCompletedSession(ctx, session, entry)
		if err != nil {
			errChan <- fmt.Errorf("unable to complete session for entry: %s : %w", entry.Title, err)
			return
		}
	} else {
		session.LastSeenAt = time.Now()
		err = p.sessionService.UpdateSession(ctx, session)
		if err != nil {
			errChan <- fmt.Errorf("unable to update session for entry: %s : %w", entry.Title, err)
			return
		}
	}

	errChan <- nil
	return
}

// isCompleted checks if at least 70% of the track has been played
func (p *SubsonicPoller) isCompleted(session *models.Session) bool {
	elapsed := time.Now().Sub(session.StartedAt).Milliseconds()
	threshold := float64(session.Duration) * 0.7
	return float64(elapsed) >= threshold
}
