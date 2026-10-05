package importers

import (
	"context"
)

type Importer interface {
	Name() string
	Import(ctx context.Context) error
}
