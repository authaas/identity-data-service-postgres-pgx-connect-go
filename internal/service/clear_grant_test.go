//revive:disable:package-comments
package service

import (
	"testing"

	"connectrpc.com/connect/v2"
)

func TestClearGrant(t *testing.T) {
	t.Run("clears the grant", func(t *testing.T) {
		server := newServer(&queriesStub{rows: 1}, pingerStub{})

		_, err := server.ClearGrant(t.Context(), clearGrantRequest(principalID, digest))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("refuses a key the type cannot hold", func(t *testing.T) {
		server := newServer(&queriesStub{}, pingerStub{})

		_, err := server.ClearGrant(t.Context(), clearGrantRequest(malformed, digest))

		assertCode(t, err, connect.CodeInvalidArgument)
	})

	t.Run("refuses when the compare did not hold", func(t *testing.T) {
		server := newServer(&queriesStub{rows: 0}, pingerStub{})

		_, err := server.ClearGrant(t.Context(), clearGrantRequest(principalID, digest))

		assertCode(t, err, connect.CodeFailedPrecondition)
	})

	t.Run("reports a statement that failed", func(t *testing.T) {
		server := newServer(&queriesStub{err: errStore}, pingerStub{})

		_, err := server.ClearGrant(t.Context(), clearGrantRequest(principalID, digest))

		assertCode(t, err, connect.CodeInternal)
	})
}
