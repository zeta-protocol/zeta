package adapters

import (
	apipb "github.com/zeta-protocol/zeta/protos/zeta/api/v1"
	nodetypes "github.com/zeta-protocol/zeta/wallet/api/node/types"
)

func toSpamStatistic(st *apipb.SpamStatistic) *nodetypes.SpamStatistic {
	if st == nil {
		// can happen if pointing to an older version of core where this
		// particular spam statistic doesn't exist yet
		return &nodetypes.SpamStatistic{}
	}
	return &nodetypes.SpamStatistic{
		CountForEpoch: st.CountForEpoch,
		MaxForEpoch:   st.MaxForEpoch,
		BannedUntil:   st.BannedUntil,
	}
}
