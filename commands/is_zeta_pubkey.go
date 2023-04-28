package commands

import (
	"encoding/hex"
	"errors"
)

const zetaPubkeyLen = 64

var (
	ErrShouldBeAValidZetaPubkey = errors.New("should be a valid zeta public key")
	ErrShouldBeAValidZetaID     = errors.New("should be a valid zeta ID")
)

// IsZetaPubkey check if a string is a valid zeta public key.
// A zeta public key is a string of 64 characters containing only hexadecimal characters.
func IsZetaPubkey(pk string) bool {
	pkLen := len(pk)
	_, err := hex.DecodeString(pk)
	return pkLen == zetaPubkeyLen && err == nil
}
