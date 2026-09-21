//revive:disable:package-comments
package service

import (
	"context"

	errors "github.com/pbrpc/connect-errors"

	"buf.build/gen/go/authaas/identity-data/protocolbuffers/go/identity/data"
	store "github.com/authaas/data-connect-go"
	"github.com/authaas/identity-data-bindings-connect-go/identity/data/dataconnect"
	identity "github.com/authaas/identity-connect-go"
	"github.com/authaas/identity-pgx-go/principal"
	ops "github.com/authaas/identity-schema-postgres-bindings-pgx-go"
)

// ClearGrant clears the outstanding grant while its digest still equals the
// one presented. The compare and the write are one statement; zero rows
// means no row satisfied it, and nothing was written.
func (s *Server) ClearGrant(ctx context.Context, req *data.ClearGrantRequest) (*data.ClearGrantResponse, error) {
	id, err := principal.Key(req.GetPrincipal())
	if err != nil {
		return nil, identity.InvalidPrincipal(ctx)
	}

	rows, err := s.queries.ClearGrant(ctx, ops.ClearGrantParams{
		ID:        id,
		GrantHash: req.GetGrantHash().GetBytes(),
	})
	if err != nil {
		return nil, store.StoreFailed(ctx, dataconnect.ServiceName, "clear the grant", err)
	}

	if rows == 0 {
		return nil, errors.PreconditionFailed(
			ctx, "grant not cleared", "GRANT_NOT_CLEARED", "grant_hash",
			"no outstanding grant with this digest",
		)
	}

	return &data.ClearGrantResponse{}, nil
}
