package wallets

import (
	"fmt"

	"github.com/zeta-protocol/zeta/paths"
	wstorev1 "github.com/zeta-protocol/zeta/wallet/wallet/store/v1"
)

// InitialiseStore builds a wallet Store specifically for users wallets.
func InitialiseStore(zetaHome string, withFileWatcher bool) (*wstorev1.FileStore, error) {
	p := paths.New(zetaHome)
	return InitialiseStoreFromPaths(p, withFileWatcher)
}

// InitialiseStoreFromPaths builds a wallet Store specifically for users wallets.
func InitialiseStoreFromPaths(zetaPaths paths.Paths, withFileWatcher bool) (*wstorev1.FileStore, error) {
	walletsHome, err := zetaPaths.CreateDataPathFor(paths.WalletsDataHome)
	if err != nil {
		return nil, fmt.Errorf("couldn't get wallets data home path: %w", err)
	}
	return wstorev1.InitialiseStore(walletsHome, withFileWatcher)
}
