//revive:disable:package-comments
package postgres

import (
	"context"
	"net"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database is the pool the statements run on.
type Database struct {
	*pgxpool.Pool
}

// New returns a Database over a pool built for the configured connection
// string. Building the pool dials nothing; a string the pool cannot parse is
// refused here.
func New(ctx context.Context, cfg Configuration) (*Database, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	return &Database{Pool: pool}, nil
}

// Address of the database, without the credentials the connection string
// carries.
func (d *Database) Address() string {
	config := d.Config().ConnConfig

	return net.JoinHostPort(config.Host, strconv.Itoa(int(config.Port)))
}
