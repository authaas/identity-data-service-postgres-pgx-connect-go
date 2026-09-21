//revive:disable:package-comments
package service

import (
	"errors"

	"github.com/jackc/pgx/v5/pgtype"

	"buf.build/gen/go/authaas/identity/protocolbuffers/go/identity"
)

// ErrInvalidKey reports a stored key that does not read back as text, which
// a NOT NULL uuid column never yields.
var ErrInvalidKey = errors.New("invalid principal key")

// key answers with the stored form of a principal's id, or the type's own
// refusal of a value it cannot hold.
func key(principal *identity.Principal) (pgtype.UUID, error) {
	var id pgtype.UUID

	err := id.Scan(principal.GetId())

	return id, err
}

// principal answers with the message for a stored key.
//
// pgtype.UUID.Value has no path that returns an error: a key with no value
// answers nil, and encoding sixteen bytes to text cannot fail. The error is
// not consulted, and nil in place of the text is the one failure.
func principal(id pgtype.UUID) (*identity.Principal, error) {
	value, _ := id.Value()

	text, ok := value.(string)
	if !ok {
		return nil, ErrInvalidKey
	}

	return identity.Principal_builder{Id: text}.Build(), nil
}
