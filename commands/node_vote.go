package commands

import commandspb "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"

func CheckNodeVote(cmd *commandspb.NodeVote) error {
	return checkNodeVote(cmd).ErrorOrNil()
}

func checkNodeVote(cmd *commandspb.NodeVote) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("node_vote", ErrIsRequired)
	}

	if len(cmd.Reference) == 0 {
		errs.AddForProperty("node_vote.reference", ErrIsRequired)
	} else if len(cmd.Reference) > 1000 {
		errs.AddForProperty("node_vote.reference", ErrReferenceTooLong)
	}

	return errs
}
