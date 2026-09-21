//revive:disable:package-comments
package service

import (
	"context"

	errors "github.com/pbrpc/connect-errors"

	"buf.build/gen/go/authaas/identity-data/protocolbuffers/go/identity/data"
)

// Delete removes the identity.
func (s *Server) Delete(ctx context.Context, req *data.DeleteRequest) (*data.DeleteResponse, error) {
	id, err := key(req.GetPrincipal())
	if err != nil {
		return nil, invalidKey(ctx)
	}

	rows, err := s.queries.DeleteIdentity(ctx, id)
	if err != nil {
		return nil, storeFailed(ctx, "delete the identity", err)
	}

	if rows == 0 {
		return nil, errors.NotFound(ctx, "identity", req.GetPrincipal().GetId())
	}

	return &data.DeleteResponse{}, nil
}
