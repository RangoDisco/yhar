package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/filters"
	"github.com/rangodisco/yhar/internal/api/types/subsonic"
	"gorm.io/gorm"
)

type SessionService struct {
	repo     *repositories.SessionRepository
	scrobble *ScrobbleService
}

func NewSessionService(repo *repositories.SessionRepository, scrobble *ScrobbleService) *SessionService {
	return &SessionService{
		repo:     repo,
		scrobble: scrobble,
	}
}

func (s *SessionService) GetOrCreateSession(ctx context.Context, entry subsonic.Entry) (*models.Session, error) {
	queryFilters := []filters.QueryFilter{
		{Key: "username", Value: entry.Username},
		{Key: "player_id", Value: entry.PlayerID},
		{Key: "title", Value: entry.Title},
	}

	session, err := s.repo.FindByFilters(ctx, queryFilters)
	// In case an error occurred, and it's not a gorm not found, skip and return error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if session != nil {
		return session, nil
	}

	// Shitty fix to prevent starting a new session for the last remaining 30% of a song
	// bc we consider a track as "completed" when listened for at least 70% of its duration
	if entry.MinutesAgo != "0" {
		return nil, nil
	}

	model := &models.Session{
		PlayerID:  entry.PlayerID,
		Username:  entry.Username,
		Title:     entry.Title,
		Artist:    entry.Artist,
		Album:     entry.Album,
		StartedAt: time.Now(),
	}

	duration, err := time.ParseDuration(fmt.Sprintf("%ss", entry.Duration))
	if err != nil {
		return nil, fmt.Errorf("unable to parse duration: %w", err)
	}
	model.Duration = duration

	err = s.repo.Persist(ctx, model)
	if err != nil {
		return nil, fmt.Errorf("unable to persist new model: %w", err)
	}

	return model, nil
}

// HandleCompletedSession updates the session timestamp fields and trigger a scrobble creation
func (s *SessionService) HandleCompletedSession(ctx context.Context, session *models.Session, entry *UnifiedScrobbleEntry) error {
	now := time.Now()
	session.CompletedAt = &now
	session.LastSeenAt = now

	_, err := s.scrobble.HandleNewScrobble(ctx, entry)
	if err != nil {
		return fmt.Errorf("unable to create new scrobble: %w", err)
	}

	err = s.UpdateSession(ctx, session)
	if err != nil {
		return fmt.Errorf("unable to update session: %w", err)
	}

	return nil
}

func (s *SessionService) UpdateSession(ctx context.Context, session *models.Session) error {
	return s.repo.Update(ctx, session)
}
