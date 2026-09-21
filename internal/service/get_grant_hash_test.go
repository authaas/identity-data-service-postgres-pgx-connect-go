//revive:disable:package-comments
package service

import (
	"bytes"
	"testing"

	"connectrpc.com/connect/v2"
	"github.com/jackc/pgx/v5"
)

func TestGetGrantHash(t *testing.T) {
	t.Run("answers with the digest", func(t *testing.T) {
		server := newServer(&queriesStub{hash: digest}, pingerStub{})

		res, err := server.GetGrantHash(t.Context(), getGrantHashRequest(principalID))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := res.GetGrantHash().GetBytes(); !bytes.Equal(got, digest) {
			t.Errorf("grant hash = %x, want %x", got, digest)
		}
	})

	t.Run("answers with nothing when none is outstanding", func(t *testing.T) {
		server := newServer(&queriesStub{hash: nil}, pingerStub{})

		res, err := server.GetGrantHash(t.Context(), getGrantHashRequest(principalID))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.HasGrantHash() {
			t.Error("expected no grant hash")
		}
	})

	t.Run("refuses a key the type cannot hold", func(t *testing.T) {
		server := newServer(&queriesStub{}, pingerStub{})

		_, err := server.GetGrantHash(t.Context(), getGrantHashRequest(malformed))

		assertCode(t, err, connect.CodeInvalidArgument)
	})

	t.Run("refuses an unknown identity", func(t *testing.T) {
		server := newServer(&queriesStub{err: pgx.ErrNoRows}, pingerStub{})

		_, err := server.GetGrantHash(t.Context(), getGrantHashRequest(principalID))

		assertCode(t, err, connect.CodeNotFound)
	})

	t.Run("reports a statement that failed", func(t *testing.T) {
		server := newServer(&queriesStub{err: errStore}, pingerStub{})

		_, err := server.GetGrantHash(t.Context(), getGrantHashRequest(principalID))

		assertCode(t, err, connect.CodeInternal)
	})
}
