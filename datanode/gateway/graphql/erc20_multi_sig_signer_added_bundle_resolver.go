package gql

import (
	"context"
	"fmt"

	v2 "github.com/zeta-protocol/zeta/protos/data-node/api/v2"
)

type erc20MultiSigSignerAddedBundleResolver ZetaResolverRoot

func (e erc20MultiSigSignerAddedBundleResolver) Timestamp(ctx context.Context, obj *v2.ERC20MultiSigSignerAddedBundle) (string, error) {
	return fmt.Sprint(obj.Timestamp), nil
}
