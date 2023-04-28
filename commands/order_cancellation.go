package commands

import (
	commandspb "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"
)

func CheckOrderCancellation(cmd *commandspb.OrderCancellation) error {
	return checkOrderCancellation(cmd).ErrorOrNil()
}

func checkOrderCancellation(cmd *commandspb.OrderCancellation) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("order_cancellation", ErrIsRequired)
	}

	if len(cmd.MarketId) > 0 && !IsZetaPubkey(cmd.MarketId) {
		errs.AddForProperty("order_cancellation.market_id", ErrShouldBeAValidZetaID)
	}

	if len(cmd.OrderId) > 0 && !IsZetaPubkey(cmd.OrderId) {
		errs.AddForProperty("order_cancellation.order_id", ErrShouldBeAValidZetaID)
	}

	return errs
}
