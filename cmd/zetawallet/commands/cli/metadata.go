package cli

import (
	"strings"

	"github.com/zeta-protocol/zeta/cmd/zetawallet/commands/flags"
	"github.com/zeta-protocol/zeta/wallet/wallet"
)

func ParseMetadata(rawMetadata []string) ([]wallet.Metadata, error) {
	if len(rawMetadata) == 0 {
		return nil, nil
	}

	metadata := make([]wallet.Metadata, 0, len(rawMetadata))
	for _, v := range rawMetadata {
		rawMeta := strings.Split(v, ":")
		if len(rawMeta) != 2 { //nolint:gomnd
			return nil, flags.InvalidFlagFormatError("meta")
		}
		metadata = append(metadata, wallet.Metadata{Key: rawMeta[0], Value: rawMeta[1]})
	}

	return metadata, nil
}
