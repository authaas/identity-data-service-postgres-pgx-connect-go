//revive:disable:package-comments
package service

import (
	"context"
	stderrors "errors"

	"github.com/jackc/pgx/v5"

	errors "github.com/pbrpc/connect-errors"

	"buf.build/gen/go/authaas/identity-data/protocolbuffers/go/identity/data"
	store "github.com/authaas/data-connect-go"
	"github.com/authaas/identity-data-bindings-connect-go/identity/data/dataconnect"
	identity "github.com/authaas/identity-connect-go"
	"github.com/authaas/identity-pgx-go/principal"
)

// Get answers with the identity's record.
func (s *Server) Get(ctx context.Context, req *data.GetRequest) (*data.GetResponse, error) {
	id, err := principal.Key(req.GetPrincipal())
	if err != nil {
		return nil, identity.InvalidPrincipal(ctx)
	}

	row, err := s.queries.GetIdentity(ctx, id)
	if stderrors.Is(err, pgx.ErrNoRows) {
		return nil, errors.NotFound(ctx, "identity", req.GetPrincipal().GetId())
	}

	if err != nil {
		return nil, store.StoreFailed(ctx, dataconnect.ServiceName, "read the identity", err)
	}

	record, err := record(row)
	if err != nil {
		return nil, store.StoreFailed(ctx, dataconnect.ServiceName, "read the identity", err)
	}

	return data.GetResponse_builder{Record: record}.Build(), nil
}
