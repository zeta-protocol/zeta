package commands

import (
	"encoding/hex"

	commandspb "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"
)

func CheckValidatorHeartbeat(cmd *commandspb.ValidatorHeartbeat) error {
	return checkValidatorHeartbeat(cmd).ErrorOrNil()
}

func checkValidatorHeartbeat(cmd *commandspb.ValidatorHeartbeat) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("validator_heartbeat", ErrIsRequired)
	}

	if len(cmd.NodeId) == 0 {
		errs.AddForProperty("validator_heartbeat.node_id", ErrIsRequired)
	} else {
		if !IsZetaPubkey(cmd.NodeId) {
			errs.AddForProperty("validator_heartbeat.node_id", ErrShouldBeAValidZetaPubkey)
		}
	}

	if cmd.EthereumSignature == nil || len(cmd.EthereumSignature.Value) == 0 {
		errs.AddForProperty("validator_heartbeat.ethereum_signature.value", ErrIsRequired)
	} else {
		_, err := hex.DecodeString(cmd.EthereumSignature.Value)
		if err != nil {
			errs.AddForProperty("validator_heartbeat.ethereum_signature.value", ErrShouldBeHexEncoded)
		}
	}

	if cmd.ZetaSignature == nil || len(cmd.ZetaSignature.Value) == 0 {
		errs.AddForProperty("validator_heartbeat.zeta_signature.value", ErrIsRequired)
	} else {
		_, err := hex.DecodeString(cmd.ZetaSignature.Value)
		if err != nil {
			errs.AddForProperty("validator_heartbeat.zeta_signature.value", ErrShouldBeHexEncoded)
		}
	}

	if len(cmd.ZetaSignature.Algo) == 0 {
		errs.AddForProperty("validator_heartbeat.zeta_signature.algo", ErrIsRequired)
	}

	return errs
}
