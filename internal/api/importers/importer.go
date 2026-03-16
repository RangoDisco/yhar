package importers

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/services"
)

type Importer interface {
	Name() string
	Import(ctx context.Context) error
	parseToUnifiedScrobble(ctx context.Context, entry interface{}) (*services.UnifiedScrobbleEntry, error)
}
