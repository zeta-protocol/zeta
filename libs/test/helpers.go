package test

import (
	"path/filepath"

	vgrand "github.com/zeta-protocol/zeta/libs/rand"
)

func RandomPath() string {
	return filepath.Join("/tmp", "zeta_tests", vgrand.RandomStr(10))
}
