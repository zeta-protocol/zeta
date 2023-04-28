package entities

import (
	"time"

	v2 "github.com/zeta-protocol/zeta/protos/data-node/api/v2"
	"github.com/zeta-protocol/zeta/protos/zeta"
)

type NodeBasic struct {
	ID              NodeID
	PubKey          ZetaPublicKey       `db:"zeta_pub_key"`
	TmPubKey        TendermintPublicKey `db:"tendermint_pub_key"`
	EthereumAddress EthereumAddress
	InfoURL         string
	Location        string
	Status          NodeStatus
	Name            string
	AvatarURL       string
	TxHash          TxHash
	ZetaTime        time.Time
}

func (n NodeBasic) ToProto() *v2.NodeBasic {
	return &v2.NodeBasic{
		Id:              n.ID.String(),
		PubKey:          n.PubKey.String(),
		TmPubKey:        n.TmPubKey.String(),
		EthereumAddress: n.EthereumAddress.String(),
		InfoUrl:         n.InfoURL,
		Location:        n.Location,
		Status:          zeta.NodeStatus(n.Status),
		Name:            n.Name,
		AvatarUrl:       n.AvatarURL,
	}
}
