package commands

import (
	"encoding/hex"

	"github.com/zeta-protocol/zeta/libs/crypto"
	commandspb "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"
)

func CheckAnnounceNode(cmd *commandspb.AnnounceNode) error {
	return checkAnnounceNode(cmd).ErrorOrNil()
}

func checkAnnounceNode(cmd *commandspb.AnnounceNode) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("announce_node", ErrIsRequired)
	}

	if len(cmd.ZetaPubKey) == 0 {
		errs.AddForProperty("announce_node.zeta_pub_key", ErrIsRequired)
	} else if !IsZetaPubkey(cmd.ZetaPubKey) {
		errs.AddForProperty("announce_node.zeta_pub_key", ErrShouldBeAValidZetaPubkey)
	}

	if len(cmd.Id) == 0 {
		errs.AddForProperty("announce_node.id", ErrIsRequired)
	} else if !IsZetaPubkey(cmd.Id) {
		errs.AddForProperty("announce_node.id", ErrShouldBeAValidZetaPubkey)
	}

	if len(cmd.EthereumAddress) == 0 {
		errs.AddForProperty("announce_node.ethereum_address", ErrIsRequired)
	} else if !crypto.EthereumIsValidAddress(cmd.EthereumAddress) {
		errs.AddForProperty("announce_node.ethereum_address", ErrIsNotValidEthereumAddress)
	}

	if len(cmd.ChainPubKey) == 0 {
		errs.AddForProperty("announce_node.chain_pub_key", ErrIsRequired)
	}

	if cmd.EthereumSignature == nil || len(cmd.EthereumSignature.Value) == 0 {
		errs.AddForProperty("announce_node.ethereum_signature", ErrIsRequired)
	} else {
		_, err := hex.DecodeString(cmd.EthereumSignature.Value)
		if err != nil {
			errs.AddForProperty("announce_node.ethereum_signature.value", ErrShouldBeHexEncoded)
		}
	}

	if cmd.ZetaSignature == nil || len(cmd.ZetaSignature.Value) == 0 {
		errs.AddForProperty("announce_node.zeta_signature", ErrIsRequired)
	} else {
		_, err := hex.DecodeString(cmd.ZetaSignature.Value)
		if err != nil {
			errs.AddForProperty("announce_node.zeta_signature.value", ErrShouldBeHexEncoded)
		}
	}

	if len(cmd.SubmitterAddress) != 0 && !crypto.EthereumIsValidAddress(cmd.SubmitterAddress) {
		errs.AddForProperty("announce_node.submitter_address", ErrIsNotValidEthereumAddress)
	}

	return errs
}
