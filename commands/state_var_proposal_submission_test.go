package commands_test

import (
	"testing"

	"github.com/zeta-protocol/zeta/commands"
	vgcrypto "github.com/zeta-protocol/zeta/libs/crypto"
	"github.com/zeta-protocol/zeta/protos/zeta"
	commandspb "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"

	"github.com/stretchr/testify/assert"
)

func TestNilStateVarProposalFundsFails(t *testing.T) {
	err := checkStateVarProposal(nil)
	assert.Contains(t, err.Get("state_variable_proposal"), commands.ErrIsRequired)
}

func TestStateVarProposals(t *testing.T) {
	cases := []struct {
		stateVar  commandspb.StateVariableProposal
		errString string
	}{
		{
			stateVar: commandspb.StateVariableProposal{
				Proposal: &zeta.StateValueProposal{
					StateVarId: vgcrypto.RandomHash(),
					EventId:    "",
					Kvb: []*zeta.KeyValueBundle{
						{
							Key:       vgcrypto.RandomHash(),
							Tolerance: "11000",
							Value:     &zeta.StateVarValue{},
						},
					},
				},
			},
			errString: "state_variable_proposal.event_id (is required)",
		},
		{
			stateVar: commandspb.StateVariableProposal{
				Proposal: &zeta.StateValueProposal{
					StateVarId: "",
					EventId:    vgcrypto.RandomHash(),
					Kvb: []*zeta.KeyValueBundle{
						{
							Key:       vgcrypto.RandomHash(),
							Tolerance: "11000",
							Value:     &zeta.StateVarValue{},
						},
					},
				},
			},
			errString: "state_variable_proposal.state_var_id (is required)",
		},
		{
			stateVar: commandspb.StateVariableProposal{
				Proposal: &zeta.StateValueProposal{
					StateVarId: "",
					EventId:    vgcrypto.RandomHash(),
					Kvb: []*zeta.KeyValueBundle{
						{
							Key:       vgcrypto.RandomHash(),
							Tolerance: "11000",
							Value:     nil,
						},
					},
				},
			},
			errString: "state_variable_proposal.key_value_bundle.0.value (is required)",
		},
		{
			stateVar: commandspb.StateVariableProposal{
				Proposal: &zeta.StateValueProposal{
					StateVarId: vgcrypto.RandomHash(),
					EventId:    vgcrypto.RandomHash(),
					Kvb: []*zeta.KeyValueBundle{
						{
							Key:       vgcrypto.RandomHash(),
							Tolerance: "11000",
							Value:     &zeta.StateVarValue{},
						},
					},
				},
			},
			errString: "",
		},
	}

	for _, c := range cases {
		err := commands.CheckStateVariableProposal(&c.stateVar)
		if len(c.errString) <= 0 {
			assert.NoError(t, err)
			continue
		}
		assert.Contains(t, err.Error(), c.errString)
	}
}

func checkStateVarProposal(cmd *commandspb.StateVariableProposal) commands.Errors {
	err := commands.CheckStateVariableProposal(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}

	return e
}
