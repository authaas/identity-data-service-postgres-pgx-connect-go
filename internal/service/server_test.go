//revive:disable:package-comments
package service

import (
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("create a server", func(_ *testing.T) {
		New(&queriesStub{}, pingerStub{}, "db:5432")
	})
}
