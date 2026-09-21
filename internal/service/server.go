//revive:disable:package-comments
package service

import (
	"log/slog"

	"git.sonicoriginal.software/logger"

	"github.com/authaas/identity-data-bindings-connect-go/identity/data/dataconnect"
)

// Server serves identity.data.Service over the generated statements.
type Server struct {
	dataconnect.UnimplementedServiceHandler

	log     *slog.Logger
	queries Queries
	db      Pinger
	address string
}

// New returns a Server over queries, reporting db under StorageCheckName at
// address.
func New(log *slog.Logger, queries Queries, db Pinger, address string) *Server {
	if log == nil {
		log = logger.NewNullLogger()
	}

	return &Server{log: log, queries: queries, db: db, address: address}
}
