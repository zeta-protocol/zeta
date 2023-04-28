package rand_test

import (
	"testing"

	vgrand "github.com/zeta-protocol/zeta/libs/rand"

	"github.com/stretchr/testify/assert"
)

func TestNonce(t *testing.T) {
	t.Run("Create a new nonce succeeds", testCreatingNewNonceSucceeds)
}

func testCreatingNewNonceSucceeds(t *testing.T) {
	assert.NotPanics(t, func() {
		vgrand.NewNonce()
	})
}
