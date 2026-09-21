//revive:disable:package-comments
package service

import (
	"testing"

	"connectrpc.com/connect/v2"
)

func TestDelete(t *testing.T) {
	t.Run("removes the identity", func(t *testing.T) {
		server := newServer(&queriesStub{rows: 1}, pingerStub{})

		if _, err := server.Delete(t.Context(), deleteRequest(principalID)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("refuses a key the type cannot hold", func(t *testing.T) {
		server := newServer(&queriesStub{}, pingerStub{})

		_, err := server.Delete(t.Context(), deleteRequest(malformed))

		assertCode(t, err, connect.CodeInvalidArgument)
	})

	t.Run("refuses an unknown identity", func(t *testing.T) {
		server := newServer(&queriesStub{rows: 0}, pingerStub{})

		_, err := server.Delete(t.Context(), deleteRequest(principalID))

		assertCode(t, err, connect.CodeNotFound)
	})

	t.Run("reports a statement that failed", func(t *testing.T) {
		server := newServer(&queriesStub{err: errStore}, pingerStub{})

		_, err := server.Delete(t.Context(), deleteRequest(principalID))

		assertCode(t, err, connect.CodeInternal)
	})
}
