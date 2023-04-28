package commands

import (
	commandspb "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"
)

func CheckProtocolUpgradeProposal(cmd *commandspb.ProtocolUpgradeProposal) error {
	return checkProtocolUpgradeProposal(cmd).ErrorOrNil()
}

func checkProtocolUpgradeProposal(cmd *commandspb.ProtocolUpgradeProposal) Errors {
	errs := NewErrors()
	if cmd == nil {
		return errs.FinalAddForProperty("protocol_upgrade_proposal", ErrIsRequired)
	}

	if len(cmd.ZetaReleaseTag) == 0 {
		errs.AddForProperty("protocol_upgrade_proposal.zeta_release_tag", ErrIsRequired)
	}

	return errs
}
