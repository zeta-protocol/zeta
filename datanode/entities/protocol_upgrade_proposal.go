// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.DATANODE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package entities

import (
	"encoding/json"
	"fmt"
	"time"

	v2 "github.com/zeta-protocol/zeta/protos/data-node/api/v2"
	eventspb "github.com/zeta-protocol/zeta/protos/zeta/events/v1"
)

type ProtocolUpgradeProposal struct {
	UpgradeBlockHeight uint64
	ZetaReleaseTag     string
	Approvers          []string
	Status             ProtocolUpgradeProposalStatus
	TxHash             TxHash
	ZetaTime           time.Time
}

func ProtocolUpgradeProposalFromProto(p *eventspb.ProtocolUpgradeEvent, txHash TxHash, zetaTime time.Time) ProtocolUpgradeProposal {
	proposal := ProtocolUpgradeProposal{
		UpgradeBlockHeight: p.UpgradeBlockHeight,
		ZetaReleaseTag:     p.ZetaReleaseTag,
		Approvers:          p.Approvers,
		Status:             ProtocolUpgradeProposalStatus(p.Status),
		TxHash:             txHash,
		ZetaTime:           zetaTime,
	}
	return proposal
}

func (p ProtocolUpgradeProposal) ToProto() *eventspb.ProtocolUpgradeEvent {
	return &eventspb.ProtocolUpgradeEvent{
		UpgradeBlockHeight: p.UpgradeBlockHeight,
		ZetaReleaseTag:     p.ZetaReleaseTag,
		Approvers:          p.Approvers,
		Status:             eventspb.ProtocolUpgradeProposalStatus(p.Status),
	}
}

func (p ProtocolUpgradeProposal) Cursor() *Cursor {
	pc := ProtocolUpgradeProposalCursor{
		ZetaTime:           p.ZetaTime,
		UpgradeBlockHeight: p.UpgradeBlockHeight,
		ZetaReleaseTag:     p.ZetaReleaseTag,
	}
	return NewCursor(pc.String())
}

func (p ProtocolUpgradeProposal) ToProtoEdge(_ ...any) (*v2.ProtocolUpgradeProposalEdge, error) {
	return &v2.ProtocolUpgradeProposalEdge{
		Node:   p.ToProto(),
		Cursor: p.Cursor().Encode(),
	}, nil
}

type ProtocolUpgradeProposalCursor struct {
	ZetaTime           time.Time
	UpgradeBlockHeight uint64
	ZetaReleaseTag     string
}

func (pc ProtocolUpgradeProposalCursor) String() string {
	bs, err := json.Marshal(pc)
	if err != nil {
		panic(fmt.Errorf("failed to marshal protocol upgrade proposal cursor: %w", err))
	}
	return string(bs)
}

func (pc *ProtocolUpgradeProposalCursor) Parse(cursorString string) error {
	if cursorString == "" {
		return nil
	}
	return json.Unmarshal([]byte(cursorString), pc)
}
