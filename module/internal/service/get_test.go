//revive:disable:package-comments
package service

import (
	"bytes"
	"testing"

	"connectrpc.com/connect/v2"
	"github.com/jackc/pgx/v5"

	ops "github.com/authaas/identity-schema-postgres-bindings-pgx-go"
)

func TestGet(t *testing.T) {
	t.Run("answers with the record", func(t *testing.T) {
		row := ops.Principal{
			ID:                    storedKey(t),
			Name:                  "name",
			DisplayName:           "display",
			CreationDate:          1,
			LastAuthenticatedDate: 2,
			GrantHash:             digest,
		}
		server := newServer(&queriesStub{identity: row}, pingerStub{})

		res, err := server.Get(t.Context(), getRequest(principalID))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		record := res.GetRecord()

		if got := record.GetPrincipal().GetId(); got != principalID {
			t.Errorf("principal = %q, want %q", got, principalID)
		}

		profile := record.GetProfile()
		if profile.GetName() != "name" || profile.GetDisplayName() != "display" {
			t.Errorf("profile = %v, want name and display", profile)
		}

		created, authenticated := record.GetCreationDate(), record.GetLastAuthenticatedDate()
		if created != 1 || authenticated != 2 {
			t.Errorf("dates = %d, %d, want 1, 2", created, authenticated)
		}

		if got := record.GetGrantHash().GetBytes(); !bytes.Equal(got, digest) {
			t.Errorf("grant hash = %x, want %x", got, digest)
		}
	})

	t.Run("leaves the grant hash unset when none is outstanding", func(t *testing.T) {
		row := ops.Principal{ID: storedKey(t)}
		server := newServer(&queriesStub{identity: row}, pingerStub{})

		res, err := server.Get(t.Context(), getRequest(principalID))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.GetRecord().HasGrantHash() {
			t.Error("expected no grant hash")
		}
	})

	t.Run("refuses a key the type cannot hold", func(t *testing.T) {
		server := newServer(&queriesStub{}, pingerStub{})

		_, err := server.Get(t.Context(), getRequest(malformed))

		assertCode(t, err, connect.CodeInvalidArgument)
	})

	t.Run("refuses an unknown identity", func(t *testing.T) {
		server := newServer(&queriesStub{err: pgx.ErrNoRows}, pingerStub{})

		_, err := server.Get(t.Context(), getRequest(principalID))

		assertCode(t, err, connect.CodeNotFound)
	})

	t.Run("reports a statement that failed", func(t *testing.T) {
		server := newServer(&queriesStub{err: errStore}, pingerStub{})

		_, err := server.Get(t.Context(), getRequest(principalID))

		assertCode(t, err, connect.CodeInternal)
	})

	t.Run("reports a row whose key does not read back", func(t *testing.T) {
		server := newServer(&queriesStub{identity: ops.Principal{}}, pingerStub{})

		_, err := server.Get(t.Context(), getRequest(principalID))

		assertCode(t, err, connect.CodeInternal)
	})
}
