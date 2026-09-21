//revive:disable:package-comments
package service

import (
	"buf.build/gen/go/authaas/identity-data/protocolbuffers/go/identity/data"
	"buf.build/gen/go/authaas/identity/protocolbuffers/go/identity"
	"buf.build/gen/go/authaas/token/protocolbuffers/go/token"
	"github.com/authaas/identity-pgx-go/principal"
	ops "github.com/authaas/identity-schema-postgres-bindings-pgx-go"
)

// record answers with the contract's record for a stored row.
func record(row ops.Principal) (*data.Record, error) {
	key, err := principal.FromKey(row.ID)
	if err != nil {
		return nil, err
	}

	return data.Record_builder{
		Principal:             key,
		Profile:               identity.Profile_builder{Name: row.Name, DisplayName: row.DisplayName}.Build(),
		CreationDate:          row.CreationDate,
		LastAuthenticatedDate: row.LastAuthenticatedDate,
		GrantHash:             grantHash(row.GrantHash),
	}.Build(), nil
}

// grantHash answers with the digest as the contract carries it: unset when
// none is outstanding.
func grantHash(digest []byte) *token.GrantHash {
	if digest == nil {
		return nil
	}

	return token.GrantHash_builder{Bytes: digest}.Build()
}
