//revive:disable:package-comments
package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v11"
)

// Resolve builds the Database from its configuration in the environment and
// reaches it before answering: a service that cannot reach its database
// answers nothing usefully, so that is a startup failure rather than
// something every request discovers.
func Resolve(ctx context.Context, log *slog.Logger) (*Database, error) {
	cfg, err := env.ParseAs[Configuration]()
	if err != nil {
		return nil, err
	}

	db, err := New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()

		return nil, fmt.Errorf("database not reachable: %w", err)
	}

	log.InfoContext(ctx, "Using postgres", "address", db.Address())

	return db, nil
}
