//revive:disable:package-comments
package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	ops "github.com/authaas/identity-schema-postgres-bindings-pgx-go"
)

// Queries is what the handlers run: the generated statements, one per RPC.
type Queries interface {
	GetIdentity(ctx context.Context, id pgtype.UUID) (ops.Principal, error)
	DeleteIdentity(ctx context.Context, id pgtype.UUID) (int64, error)
	GetGrantHash(ctx context.Context, id pgtype.UUID) ([]byte, error)
	ClearGrant(ctx context.Context, arg ops.ClearGrantParams) (int64, error)
}
