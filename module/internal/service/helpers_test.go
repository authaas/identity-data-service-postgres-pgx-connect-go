//revive:disable:package-comments
package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"connectrpc.com/connect/v2"
	"github.com/jackc/pgx/v5/pgtype"

	"buf.build/gen/go/authaas/identity-data/protocolbuffers/go/identity/data"
	"buf.build/gen/go/authaas/identity/protocolbuffers/go/identity"
	"buf.build/gen/go/authaas/token/protocolbuffers/go/token"
	ops "github.com/authaas/identity-operations-postgres-pgx-go"
)

const (
	// principalID is a canonical UUID.
	principalID = "01234567-89ab-4def-8123-456789abcdef"

	// malformed is a string the key type cannot hold.
	malformed = "not-a-uuid"
)

// digest stands in for a stored grant hash.
var digest = []byte("0123456789abcdef0123456789abcdef")

var errStore = errors.New("store failed")

// queriesStub answers each statement with what the test set.
type queriesStub struct {
	identity ops.Principal
	rows     int64
	hash     []byte
	err      error
}

func (q *queriesStub) GetIdentity(context.Context, pgtype.UUID) (ops.Principal, error) {
	return q.identity, q.err
}

func (q *queriesStub) DeleteIdentity(context.Context, pgtype.UUID) (int64, error) {
	return q.rows, q.err
}

func (q *queriesStub) GetGrantHash(context.Context, pgtype.UUID) ([]byte, error) {
	return q.hash, q.err
}

func (q *queriesStub) ClearGrant(context.Context, ops.ClearGrantParams) (int64, error) {
	return q.rows, q.err
}

// pingerStub answers Ping with what the test set.
type pingerStub struct {
	err error
}

func (p pingerStub) Ping(context.Context) error { return p.err }

// newServer builds a Server on the stubs.
func newServer(queries *queriesStub, db pingerStub) *Server {
	return New(slog.New(slog.DiscardHandler), queries, db, "db:5432")
}

// principalOf builds the request key for id.
func principalOf(id string) *identity.Principal {
	return identity.Principal_builder{Id: id}.Build()
}

// grantHashOf builds the request digest.
func grantHashOf(b []byte) *token.GrantHash {
	return token.GrantHash_builder{Bytes: b}.Build()
}

func getRequest(id string) *data.GetRequest {
	return data.GetRequest_builder{Principal: principalOf(id)}.Build()
}

func deleteRequest(id string) *data.DeleteRequest {
	return data.DeleteRequest_builder{Principal: principalOf(id)}.Build()
}

func getGrantHashRequest(id string) *data.GetGrantHashRequest {
	return data.GetGrantHashRequest_builder{Principal: principalOf(id)}.Build()
}

func clearGrantRequest(id string, b []byte) *data.ClearGrantRequest {
	return data.ClearGrantRequest_builder{
		Principal: principalOf(id),
		GrantHash: grantHashOf(b),
	}.Build()
}

// storedKey is principalID in the stored form.
func storedKey(t *testing.T) pgtype.UUID {
	t.Helper()

	id, err := key(principalOf(principalID))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return id
}

// assertCode fails the test unless err is a *connect.Error carrying want.
func assertCode(t *testing.T, err error, want connect.Code) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected %v, got no error", want)
	}

	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		t.Fatalf("expected a *connect.Error, got %v", err)
	}

	if got := connect.CodeOf(err); got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}
