//revive:disable:package-comments
package service

import (
	"context"
	stderrors "errors"

	"github.com/jackc/pgx/v5"

	errors "github.com/pbrpc/connect-errors"

	"buf.build/gen/go/authaas/identity-data/protocolbuffers/go/identity/data"
	principal "github.com/authaas/identity-pgx-go"
)

// GetGrantHash answers with the digest of the identity's outstanding grant,
// unset when none is outstanding.
func (s *Server) GetGrantHash(
	ctx context.Context, req *data.GetGrantHashRequest,
) (*data.GetGrantHashResponse, error) {
	id, err := principal.Key(req.GetPrincipal())
	if err != nil {
		return nil, invalidKey(ctx)
	}

	digest, err := s.queries.GetGrantHash(ctx, id)
	if stderrors.Is(err, pgx.ErrNoRows) {
		return nil, errors.NotFound(ctx, "identity", req.GetPrincipal().GetId())
	}

	if err != nil {
		return nil, storeFailed(ctx, "read the grant hash", err)
	}

	return data.GetGrantHashResponse_builder{GrantHash: grantHash(digest)}.Build(), nil
}
