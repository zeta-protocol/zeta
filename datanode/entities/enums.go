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
	"fmt"

	"github.com/zeta-protocol/zeta/protos/zeta"
	zetapb "github.com/zeta-protocol/zeta/protos/zeta"
	commandspb "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"
	eventspb "github.com/zeta-protocol/zeta/protos/zeta/events/v1"
	"github.com/jackc/pgtype"
)

type Side = zeta.Side

const (
	// Default value, always invalid.
	SideUnspecified Side = zeta.Side_SIDE_UNSPECIFIED
	// Buy order.
	SideBuy Side = zeta.Side_SIDE_BUY
	// Sell order.
	SideSell Side = zeta.Side_SIDE_SELL
)

type TradeType = zeta.Trade_Type

const (
	// Default value, always invalid.
	TradeTypeUnspecified TradeType = zeta.Trade_TYPE_UNSPECIFIED
	// Normal trading between two parties.
	TradeTypeDefault TradeType = zeta.Trade_TYPE_DEFAULT
	// Trading initiated by the network with another party on the book,
	// which helps to zero-out the positions of one or more distressed parties.
	TradeTypeNetworkCloseOutGood TradeType = zeta.Trade_TYPE_NETWORK_CLOSE_OUT_GOOD
	// Trading initiated by the network with another party off the book,
	// with a distressed party in order to zero-out the position of the party.
	TradeTypeNetworkCloseOutBad TradeType = zeta.Trade_TYPE_NETWORK_CLOSE_OUT_BAD
)

type PeggedReference = zeta.PeggedReference

const (
	// Default value for PeggedReference, no reference given.
	PeggedReferenceUnspecified PeggedReference = zeta.PeggedReference_PEGGED_REFERENCE_UNSPECIFIED
	// Mid price reference.
	PeggedReferenceMid PeggedReference = zeta.PeggedReference_PEGGED_REFERENCE_MID
	// Best bid price reference.
	PeggedReferenceBestBid PeggedReference = zeta.PeggedReference_PEGGED_REFERENCE_BEST_BID
	// Best ask price reference.
	PeggedReferenceBestAsk PeggedReference = zeta.PeggedReference_PEGGED_REFERENCE_BEST_ASK
)

type OrderStatus = zeta.Order_Status

const (
	// Default value, always invalid.
	OrderStatusUnspecified OrderStatus = zeta.Order_STATUS_UNSPECIFIED
	// Used for active unfilled or partially filled orders.
	OrderStatusActive OrderStatus = zeta.Order_STATUS_ACTIVE
	// Used for expired GTT orders.
	OrderStatusExpired OrderStatus = zeta.Order_STATUS_EXPIRED
	// Used for orders cancelled by the party that created the order.
	OrderStatusCancelled OrderStatus = zeta.Order_STATUS_CANCELLED
	// Used for unfilled FOK or IOC orders, and for orders that were stopped by the network.
	OrderStatusStopped OrderStatus = zeta.Order_STATUS_STOPPED
	// Used for closed fully filled orders.
	OrderStatusFilled OrderStatus = zeta.Order_STATUS_FILLED
	// Used for orders when not enough collateral was available to fill the margin requirements.
	OrderStatusRejected OrderStatus = zeta.Order_STATUS_REJECTED
	// Used for closed partially filled IOC orders.
	OrderStatusPartiallyFilled OrderStatus = zeta.Order_STATUS_PARTIALLY_FILLED
	// Order has been removed from the order book and has been parked, this applies to pegged orders only.
	OrderStatusParked OrderStatus = zeta.Order_STATUS_PARKED
)

type OrderType = zeta.Order_Type

const (
	// Default value, always invalid.
	OrderTypeUnspecified OrderType = zeta.Order_TYPE_UNSPECIFIED
	// Used for Limit orders.
	OrderTypeLimit OrderType = zeta.Order_TYPE_LIMIT
	// Used for Market orders.
	OrderTypeMarket OrderType = zeta.Order_TYPE_MARKET
	// Used for orders where the initiating party is the network (with distressed traders).
	OrderTypeNetwork OrderType = zeta.Order_TYPE_NETWORK
)

type OrderTimeInForce = zeta.Order_TimeInForce

const (
	// Default value for TimeInForce, can be valid for an amend.
	OrderTimeInForceUnspecified OrderTimeInForce = zeta.Order_TIME_IN_FORCE_UNSPECIFIED
	// Good until cancelled.
	OrderTimeInForceGTC OrderTimeInForce = zeta.Order_TIME_IN_FORCE_GTC
	// Good until specified time.
	OrderTimeInForceGTT OrderTimeInForce = zeta.Order_TIME_IN_FORCE_GTT
	// Immediate or cancel.
	OrderTimeInForceIOC OrderTimeInForce = zeta.Order_TIME_IN_FORCE_IOC
	// Fill or kill.
	OrderTimeInForceFOK OrderTimeInForce = zeta.Order_TIME_IN_FORCE_FOK
	// Good for auction.
	OrderTimeInForceGFA OrderTimeInForce = zeta.Order_TIME_IN_FORCE_GFA
	// Good for normal.
	OrderTimeInForceGFN OrderTimeInForce = zeta.Order_TIME_IN_FORCE_GFN
)

type OrderError = zeta.OrderError

const (
	// Default value, no error reported.
	OrderErrorUnspecified OrderError = zeta.OrderError_ORDER_ERROR_UNSPECIFIED
	// Order was submitted for a market that does not exist.
	OrderErrorInvalidMarketID OrderError = zeta.OrderError_ORDER_ERROR_INVALID_MARKET_ID
	// Order was submitted with an invalid identifier.
	OrderErrorInvalidOrderID OrderError = zeta.OrderError_ORDER_ERROR_INVALID_ORDER_ID
	// Order was amended with a sequence number that was not previous version + 1.
	OrderErrorOutOfSequence OrderError = zeta.OrderError_ORDER_ERROR_OUT_OF_SEQUENCE
	// Order was amended with an invalid remaining size (e.g. remaining greater than total size).
	OrderErrorInvalidRemainingSize OrderError = zeta.OrderError_ORDER_ERROR_INVALID_REMAINING_SIZE
	// Node was unable to get Zeta (blockchain) time.
	OrderErrorTimeFailure OrderError = zeta.OrderError_ORDER_ERROR_TIME_FAILURE
	// Failed to remove an order from the book.
	OrderErrorRemovalFailure OrderError = zeta.OrderError_ORDER_ERROR_REMOVAL_FAILURE
	// An order with `TimeInForce.TIME_IN_FORCE_GTT` was submitted or amended
	// with an expiration that was badly formatted or otherwise invalid.
	OrderErrorInvalidExpirationDatetime OrderError = zeta.OrderError_ORDER_ERROR_INVALID_EXPIRATION_DATETIME
	// Order was submitted or amended with an invalid reference field.
	OrderErrorInvalidOrderReference OrderError = zeta.OrderError_ORDER_ERROR_INVALID_ORDER_REFERENCE
	// Order amend was submitted for an order field that cannot not be amended (e.g. order identifier).
	OrderErrorEditNotAllowed OrderError = zeta.OrderError_ORDER_ERROR_EDIT_NOT_ALLOWED
	// Amend failure because amend details do not match original order.
	OrderErrorAmendFailure OrderError = zeta.OrderError_ORDER_ERROR_AMEND_FAILURE
	// Order not found in an order book or store.
	OrderErrorNotFound OrderError = zeta.OrderError_ORDER_ERROR_NOT_FOUND
	// Order was submitted with an invalid or missing party identifier.
	OrderErrorInvalidParty OrderError = zeta.OrderError_ORDER_ERROR_INVALID_PARTY_ID
	// Order was submitted for a market that has closed.
	OrderErrorMarketClosed OrderError = zeta.OrderError_ORDER_ERROR_MARKET_CLOSED
	// Order was submitted, but the party did not have enough collateral to cover the order.
	OrderErrorMarginCheckFailed OrderError = zeta.OrderError_ORDER_ERROR_MARGIN_CHECK_FAILED
	// Order was submitted, but the party did not have an account for this asset.
	OrderErrorMissingGeneralAccount OrderError = zeta.OrderError_ORDER_ERROR_MISSING_GENERAL_ACCOUNT
	// Unspecified internal error.
	OrderErrorInternalError OrderError = zeta.OrderError_ORDER_ERROR_INTERNAL_ERROR
	// Order was submitted with an invalid or missing size (e.g. 0).
	OrderErrorInvalidSize OrderError = zeta.OrderError_ORDER_ERROR_INVALID_SIZE
	// Order was submitted with an invalid persistence for its type.
	OrderErrorInvalidPersistance OrderError = zeta.OrderError_ORDER_ERROR_INVALID_PERSISTENCE
	// Order was submitted with an invalid type field.
	OrderErrorInvalidType OrderError = zeta.OrderError_ORDER_ERROR_INVALID_TYPE
	// Order was stopped as it would have traded with another order submitted from the same party.
	OrderErrorSelfTrading OrderError = zeta.OrderError_ORDER_ERROR_SELF_TRADING
	// Order was submitted, but the party did not have enough collateral to cover the fees for the order.
	OrderErrorInsufficientFundsToPayFees OrderError = zeta.OrderError_ORDER_ERROR_INSUFFICIENT_FUNDS_TO_PAY_FEES
	// Order was submitted with an incorrect or invalid market type.
	OrderErrorIncorrectMarketType OrderError = zeta.OrderError_ORDER_ERROR_INCORRECT_MARKET_TYPE
	// Order was submitted with invalid time in force.
	OrderErrorInvalidTimeInForce OrderError = zeta.OrderError_ORDER_ERROR_INVALID_TIME_IN_FORCE
	// A GFN order has got to the market when it is in auction mode.
	OrderErrorCannotSendGFNOrderDuringAnAuction OrderError = zeta.OrderError_ORDER_ERROR_CANNOT_SEND_GFN_ORDER_DURING_AN_AUCTION
	// A GFA order has got to the market when it is in continuous trading mode.
	OrderErrorCannotSendGFAOrderDuringContinuousTrading OrderError = zeta.OrderError_ORDER_ERROR_CANNOT_SEND_GFA_ORDER_DURING_CONTINUOUS_TRADING
	// Attempt to amend order to GTT without ExpiryAt.
	OrderErrorCannotAmendToGTTWithoutExpiryAt OrderError = zeta.OrderError_ORDER_ERROR_CANNOT_AMEND_TO_GTT_WITHOUT_EXPIRYAT
	// Attempt to amend ExpiryAt to a value before CreatedAt.
	OrderErrorExpiryAtBeforeCreatedAt OrderError = zeta.OrderError_ORDER_ERROR_EXPIRYAT_BEFORE_CREATEDAT
	// Attempt to amend to GTC without an ExpiryAt value.
	OrderErrorCannotHaveGTCAndExpiryAt OrderError = zeta.OrderError_ORDER_ERROR_CANNOT_HAVE_GTC_AND_EXPIRYAT
	// Amending to FOK or IOC is invalid.
	OrderErrorCannotAmendToFOKOrIOC OrderError = zeta.OrderError_ORDER_ERROR_CANNOT_AMEND_TO_FOK_OR_IOC
	// Amending to GFA or GFN is invalid.
	OrderErrorCannotAmendToGFAOrGFN OrderError = zeta.OrderError_ORDER_ERROR_CANNOT_AMEND_TO_GFA_OR_GFN
	// Amending from GFA or GFN is invalid.
	OrderErrorCannotAmendFromGFAOrGFN OrderError = zeta.OrderError_ORDER_ERROR_CANNOT_AMEND_FROM_GFA_OR_GFN
	// IOC orders are not allowed during auction.
	OrderErrorCannotSendIOCOrderDuringAuction OrderError = zeta.OrderError_ORDER_ERROR_CANNOT_SEND_IOC_ORDER_DURING_AUCTION
	// FOK orders are not allowed during auction.
	OrderErrorCannotSendFOKOrderDurinAuction OrderError = zeta.OrderError_ORDER_ERROR_CANNOT_SEND_FOK_ORDER_DURING_AUCTION
	// Pegged orders must be LIMIT orders.
	OrderErrorMustBeLimitOrder OrderError = zeta.OrderError_ORDER_ERROR_MUST_BE_LIMIT_ORDER
	// Pegged orders can only have TIF GTC or GTT.
	OrderErrorMustBeGTTOrGTC OrderError = zeta.OrderError_ORDER_ERROR_MUST_BE_GTT_OR_GTC
	// Pegged order must have a reference price.
	OrderErrorWithoutReferencePrice OrderError = zeta.OrderError_ORDER_ERROR_WITHOUT_REFERENCE_PRICE
	// Buy pegged order cannot reference best ask price.
	OrderErrorBuyCannotReferenceBestAskPrice OrderError = zeta.OrderError_ORDER_ERROR_BUY_CANNOT_REFERENCE_BEST_ASK_PRICE
	// Pegged order offset must be >= 0.
	OrderErrorOffsetMustBeGreaterOrEqualToZero OrderError = zeta.OrderError_ORDER_ERROR_OFFSET_MUST_BE_GREATER_OR_EQUAL_TO_ZERO
	// Sell pegged order cannot reference best bid price.
	OrderErrorSellCannotReferenceBestBidPrice OrderError = zeta.OrderError_ORDER_ERROR_SELL_CANNOT_REFERENCE_BEST_BID_PRICE
	// Pegged order offset must be > zero.
	OrderErrorOffsetMustBeGreaterThanZero OrderError = zeta.OrderError_ORDER_ERROR_OFFSET_MUST_BE_GREATER_THAN_ZERO
	// The party has an insufficient balance, or does not have
	// a general account to submit the order (no deposits made
	// for the required asset).
	OrderErrorInsufficientAssetBalance OrderError = zeta.OrderError_ORDER_ERROR_INSUFFICIENT_ASSET_BALANCE
	// Cannot amend a non pegged orders details.
	OrderErrorCannotAmendPeggedOrderDetailsOnNonPeggedOrder OrderError = zeta.OrderError_ORDER_ERROR_CANNOT_AMEND_PEGGED_ORDER_DETAILS_ON_NON_PEGGED_ORDER
	// We are unable to re-price a pegged order because a market price is unavailable.
	OrderErrorUnableToRepricePeggedOrder OrderError = zeta.OrderError_ORDER_ERROR_UNABLE_TO_REPRICE_PEGGED_ORDER
	// It is not possible to amend the price of an existing pegged order.
	OrderErrorUnableToAmendPriceOnPeggedOrder OrderError = zeta.OrderError_ORDER_ERROR_UNABLE_TO_AMEND_PRICE_ON_PEGGED_ORDER
	// An FOK, IOC, or GFN order was rejected because it resulted in trades outside the price bounds.
	OrderErrorNonPersistentOrderOutOfPriceBounds OrderError = zeta.OrderError_ORDER_ERROR_NON_PERSISTENT_ORDER_OUT_OF_PRICE_BOUNDS
)

type PositionStatus int32

const (
	PositionStatusUnspecified  = PositionStatus(zeta.PositionStatus_POSITION_STATUS_UNSPECIFIED)
	PositionStatusOrdersClosed = PositionStatus(zeta.PositionStatus_POSITION_STATUS_ORDERS_CLOSED)
	PositionStatusClosedOut    = PositionStatus(zeta.PositionStatus_POSITION_STATUS_CLOSED_OUT)
	PositionStatusDistressed   = PositionStatus(zeta.PositionStatus_POSITION_STATUS_DISTRESSED)
)

type TransferType int

const (
	Unknown TransferType = iota
	OneOff
	Recurring
)

const (
	OneOffStr    = "OneOff"
	RecurringStr = "Recurring"
	UnknownStr   = "Unknown"
)

func (m TransferType) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	mode := UnknownStr
	switch m {
	case OneOff:
		mode = OneOffStr
	case Recurring:
		mode = RecurringStr
	}

	return append(buf, []byte(mode)...), nil
}

func (m *TransferType) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val := Unknown
	switch string(src) {
	case OneOffStr:
		val = OneOff
	case RecurringStr:
		val = Recurring
	}

	*m = val
	return nil
}

type TransferStatus eventspb.Transfer_Status

const (
	TransferStatusUnspecified = TransferStatus(eventspb.Transfer_STATUS_UNSPECIFIED)
	TransferStatusPending     = TransferStatus(eventspb.Transfer_STATUS_PENDING)
	TransferStatusDone        = TransferStatus(eventspb.Transfer_STATUS_DONE)
	TransferStatusRejected    = TransferStatus(eventspb.Transfer_STATUS_REJECTED)
	TransferStatusStopped     = TransferStatus(eventspb.Transfer_STATUS_STOPPED)
	TransferStatusCancelled   = TransferStatus(eventspb.Transfer_STATUS_CANCELLED)
)

func (m TransferStatus) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	mode, ok := eventspb.Transfer_Status_name[int32(m)]
	if !ok {
		return buf, fmt.Errorf("unknown transfer status: %s", mode)
	}
	return append(buf, []byte(mode)...), nil
}

func (m *TransferStatus) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := eventspb.Transfer_Status_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown transfer status: %s", src)
	}

	*m = TransferStatus(val)
	return nil
}

type AssetStatus zeta.Asset_Status

const (
	AssetStatusUnspecified    = AssetStatus(zeta.Asset_STATUS_UNSPECIFIED)
	AssetStatusProposed       = AssetStatus(zeta.Asset_STATUS_PROPOSED)
	AssetStatusRejected       = AssetStatus(zeta.Asset_STATUS_REJECTED)
	AssetStatusPendingListing = AssetStatus(zeta.Asset_STATUS_PENDING_LISTING)
	AssetStatusEnabled        = AssetStatus(zeta.Asset_STATUS_ENABLED)
)

func (m AssetStatus) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	mode, ok := zeta.Asset_Status_name[int32(m)]
	if !ok {
		return buf, fmt.Errorf("unknown asset status: %s", mode)
	}
	return append(buf, []byte(mode)...), nil
}

func (m *AssetStatus) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zeta.Asset_Status_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown asset status: %s", src)
	}

	*m = AssetStatus(val)
	return nil
}

type MarketTradingMode zeta.Market_TradingMode

const (
	MarketTradingModeUnspecified       = MarketTradingMode(zeta.Market_TRADING_MODE_UNSPECIFIED)
	MarketTradingModeContinuous        = MarketTradingMode(zeta.Market_TRADING_MODE_CONTINUOUS)
	MarketTradingModeBatchAuction      = MarketTradingMode(zeta.Market_TRADING_MODE_BATCH_AUCTION)
	MarketTradingModeOpeningAuction    = MarketTradingMode(zeta.Market_TRADING_MODE_OPENING_AUCTION)
	MarketTradingModeMonitoringAuction = MarketTradingMode(zeta.Market_TRADING_MODE_MONITORING_AUCTION)
	MarketTradingModeNoTrading         = MarketTradingMode(zeta.Market_TRADING_MODE_NO_TRADING)
)

func (m MarketTradingMode) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	mode, ok := zeta.Market_TradingMode_name[int32(m)]
	if !ok {
		return buf, fmt.Errorf("unknown trading mode: %s", mode)
	}
	return append(buf, []byte(mode)...), nil
}

func (m *MarketTradingMode) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zeta.Market_TradingMode_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown trading mode: %s", src)
	}

	*m = MarketTradingMode(val)
	return nil
}

type MarketState zeta.Market_State

const (
	MarketStateUnspecified       = MarketState(zeta.Market_STATE_UNSPECIFIED)
	MarketStateProposed          = MarketState(zeta.Market_STATE_PROPOSED)
	MarketStateRejected          = MarketState(zeta.Market_STATE_REJECTED)
	MarketStatePending           = MarketState(zeta.Market_STATE_PENDING)
	MarketStateCancelled         = MarketState(zeta.Market_STATE_CANCELLED)
	MarketStateActive            = MarketState(zeta.Market_STATE_ACTIVE)
	MarketStateSuspended         = MarketState(zeta.Market_STATE_SUSPENDED)
	MarketStateClosed            = MarketState(zeta.Market_STATE_CLOSED)
	MarketStateTradingTerminated = MarketState(zeta.Market_STATE_TRADING_TERMINATED)
	MarketStateSettled           = MarketState(zeta.Market_STATE_SETTLED)
)

func (s MarketState) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	state, ok := zeta.Market_State_name[int32(s)]
	if !ok {
		return buf, fmt.Errorf("unknown market state: %s", state)
	}
	return append(buf, []byte(state)...), nil
}

func (s *MarketState) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zeta.Market_State_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown market state: %s", src)
	}

	*s = MarketState(val)

	return nil
}

type DepositStatus zeta.Deposit_Status

const (
	DepositStatusUnspecified = DepositStatus(zeta.Deposit_STATUS_UNSPECIFIED)
	DepositStatusOpen        = DepositStatus(zeta.Deposit_STATUS_OPEN)
	DepositStatusCancelled   = DepositStatus(zeta.Deposit_STATUS_CANCELLED)
	DepositStatusFinalized   = DepositStatus(zeta.Deposit_STATUS_FINALIZED)
)

func (s DepositStatus) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	status, ok := zeta.Deposit_Status_name[int32(s)]
	if !ok {
		return buf, fmt.Errorf("unknown deposit state, %s", status)
	}
	return append(buf, []byte(status)...), nil
}

func (s *DepositStatus) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zeta.Deposit_Status_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown deposit state: %s", src)
	}

	*s = DepositStatus(val)

	return nil
}

type WithdrawalStatus zeta.Withdrawal_Status

const (
	WithdrawalStatusUnspecified = WithdrawalStatus(zeta.Withdrawal_STATUS_UNSPECIFIED)
	WithdrawalStatusOpen        = WithdrawalStatus(zeta.Withdrawal_STATUS_OPEN)
	WithdrawalStatusRejected    = WithdrawalStatus(zeta.Withdrawal_STATUS_REJECTED)
	WithdrawalStatusFinalized   = WithdrawalStatus(zeta.Withdrawal_STATUS_FINALIZED)
)

func (s WithdrawalStatus) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	status, ok := zeta.Withdrawal_Status_name[int32(s)]
	if !ok {
		return buf, fmt.Errorf("unknown withdrawal status: %s", status)
	}
	return append(buf, []byte(status)...), nil
}

func (s *WithdrawalStatus) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zeta.Withdrawal_Status_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown withdrawal status: %s", src)
	}
	*s = WithdrawalStatus(val)
	return nil
}

/************************* Proposal State *****************************/

type ProposalState zeta.Proposal_State

const (
	ProposalStateUnspecified        = ProposalState(zeta.Proposal_STATE_UNSPECIFIED)
	ProposalStateFailed             = ProposalState(zeta.Proposal_STATE_FAILED)
	ProposalStateOpen               = ProposalState(zeta.Proposal_STATE_OPEN)
	ProposalStatePassed             = ProposalState(zeta.Proposal_STATE_PASSED)
	ProposalStateRejected           = ProposalState(zeta.Proposal_STATE_REJECTED)
	ProposalStateDeclined           = ProposalState(zeta.Proposal_STATE_DECLINED)
	ProposalStateEnacted            = ProposalState(zeta.Proposal_STATE_ENACTED)
	ProposalStateWaitingForNodeVote = ProposalState(zeta.Proposal_STATE_WAITING_FOR_NODE_VOTE)
)

func (s ProposalState) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	str, ok := zeta.Proposal_State_name[int32(s)]
	if !ok {
		return buf, fmt.Errorf("unknown state: %v", s)
	}
	return append(buf, []byte(str)...), nil
}

func (s *ProposalState) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zeta.Proposal_State_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown state: %s", src)
	}
	*s = ProposalState(val)
	return nil
}

/************************* Proposal Error *****************************/

type ProposalError zeta.ProposalError

const (
	ProposalErrorUnspecified                      = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_UNSPECIFIED)
	ProposalErrorCloseTimeTooSoon                 = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_CLOSE_TIME_TOO_SOON)
	ProposalErrorCloseTimeTooLate                 = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_CLOSE_TIME_TOO_LATE)
	ProposalErrorEnactTimeTooSoon                 = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_ENACT_TIME_TOO_SOON)
	ProposalErrorEnactTimeTooLate                 = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_ENACT_TIME_TOO_LATE)
	ProposalErrorInsufficientTokens               = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_INSUFFICIENT_TOKENS)
	ProposalErrorInvalidInstrumentSecurity        = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_INVALID_INSTRUMENT_SECURITY)
	ProposalErrorNoProduct                        = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_NO_PRODUCT)
	ProposalErrorUnsupportedProduct               = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_UNSUPPORTED_PRODUCT)
	ProposalErrorNoTradingMode                    = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_NO_TRADING_MODE)
	ProposalErrorUnsupportedTradingMode           = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_UNSUPPORTED_TRADING_MODE)
	ProposalErrorNodeValidationFailed             = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_NODE_VALIDATION_FAILED)
	ProposalErrorMissingBuiltinAssetField         = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_MISSING_BUILTIN_ASSET_FIELD)
	ProposalErrorMissingErc20ContractAddress      = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_MISSING_ERC20_CONTRACT_ADDRESS)
	ProposalErrorInvalidAsset                     = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_INVALID_ASSET)
	ProposalErrorIncompatibleTimestamps           = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_INCOMPATIBLE_TIMESTAMPS)
	ProposalErrorNoRiskParameters                 = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_NO_RISK_PARAMETERS)
	ProposalErrorNetworkParameterInvalidKey       = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_NETWORK_PARAMETER_INVALID_KEY)
	ProposalErrorNetworkParameterInvalidValue     = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_NETWORK_PARAMETER_INVALID_VALUE)
	ProposalErrorNetworkParameterValidationFailed = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_NETWORK_PARAMETER_VALIDATION_FAILED)
	ProposalErrorOpeningAuctionDurationTooSmall   = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_OPENING_AUCTION_DURATION_TOO_SMALL)
	ProposalErrorOpeningAuctionDurationTooLarge   = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_OPENING_AUCTION_DURATION_TOO_LARGE)
	ProposalErrorCouldNotInstantiateMarket        = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_COULD_NOT_INSTANTIATE_MARKET)
	ProposalErrorInvalidFutureProduct             = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_INVALID_FUTURE_PRODUCT)
	ProposalErrorInvalidRiskParameter             = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_INVALID_RISK_PARAMETER)
	ProposalErrorMajorityThresholdNotReached      = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_MAJORITY_THRESHOLD_NOT_REACHED)
	ProposalErrorParticipationThresholdNotReached = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_PARTICIPATION_THRESHOLD_NOT_REACHED)
	ProposalErrorInvalidAssetDetails              = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_INVALID_ASSET_DETAILS)
	ProposalErrorUnknownType                      = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_UNKNOWN_TYPE)
	ProposalErrorUnknownRiskParameterType         = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_UNKNOWN_RISK_PARAMETER_TYPE)
	ProposalErrorInvalidFreeform                  = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_INVALID_FREEFORM)
	ProposalErrorInsufficientEquityLikeShare      = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_INSUFFICIENT_EQUITY_LIKE_SHARE)
	ProposalErrorInvalidMarket                    = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_INVALID_MARKET)
	ProposalErrorTooManyMarketDecimalPlaces       = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_TOO_MANY_MARKET_DECIMAL_PLACES)
	ProposalErrorTooManyPriceMonitoringTriggers   = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_TOO_MANY_PRICE_MONITORING_TRIGGERS)
	ProposalErrorERC20AddressAlreadyInUse         = ProposalError(zeta.ProposalError_PROPOSAL_ERROR_ERC20_ADDRESS_ALREADY_IN_USE)
)

func (s ProposalError) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	str, ok := zeta.ProposalError_name[int32(s)]
	if !ok {
		return buf, fmt.Errorf("unknown proposal error: %v", s)
	}
	return append(buf, []byte(str)...), nil
}

func (s *ProposalError) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zeta.ProposalError_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown proposal error: %s", src)
	}
	*s = ProposalError(val)
	return nil
}

/************************* VoteValue *****************************/

type VoteValue zeta.Vote_Value

const (
	VoteValueUnspecified = VoteValue(zeta.Vote_VALUE_UNSPECIFIED)
	VoteValueNo          = VoteValue(zeta.Vote_VALUE_NO)
	VoteValueYes         = VoteValue(zeta.Vote_VALUE_YES)
)

func (s VoteValue) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	str, ok := zeta.Vote_Value_name[int32(s)]
	if !ok {
		return buf, fmt.Errorf("unknown vote value: %v", s)
	}
	return append(buf, []byte(str)...), nil
}

func (s *VoteValue) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zeta.Vote_Value_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown vote value: %s", src)
	}
	*s = VoteValue(val)
	return nil
}

/************************* NodeSignature Kind *****************************/

type NodeSignatureKind commandspb.NodeSignatureKind

const (
	NodeSignatureKindUnspecified          = NodeSignatureKind(commandspb.NodeSignatureKind_NODE_SIGNATURE_KIND_UNSPECIFIED)
	NodeSignatureKindAsset                = NodeSignatureKind(commandspb.NodeSignatureKind_NODE_SIGNATURE_KIND_ASSET_NEW)
	NodeSignatureKindAssetUpdate          = NodeSignatureKind(commandspb.NodeSignatureKind_NODE_SIGNATURE_KIND_ASSET_UPDATE)
	NodeSignatureKindAssetWithdrawal      = NodeSignatureKind(commandspb.NodeSignatureKind_NODE_SIGNATURE_KIND_ASSET_WITHDRAWAL)
	NodeSignatureKindMultisigSignerAdded  = NodeSignatureKind(commandspb.NodeSignatureKind_NODE_SIGNATURE_KIND_ERC20_MULTISIG_SIGNER_ADDED)
	NodeSignatureKindMultisigSignerRemove = NodeSignatureKind(commandspb.NodeSignatureKind_NODE_SIGNATURE_KIND_ERC20_MULTISIG_SIGNER_REMOVED)
)

func (s NodeSignatureKind) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	str, ok := commandspb.NodeSignatureKind_name[int32(s)]
	if !ok {
		return buf, fmt.Errorf("unknown state: %v", s)
	}
	return append(buf, []byte(str)...), nil
}

func (s *NodeSignatureKind) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := commandspb.NodeSignatureKind_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown state: %s", src)
	}
	*s = NodeSignatureKind(val)
	return nil
}

type (
	DataSourceSpecStatus zetapb.DataSourceSpec_Status
	OracleSpecStatus     = DataSourceSpecStatus
)

const (
	OracleSpecUnspecified = DataSourceSpecStatus(zetapb.DataSourceSpec_STATUS_UNSPECIFIED)
	OracleSpecActive      = DataSourceSpecStatus(zetapb.DataSourceSpec_STATUS_ACTIVE)
	OracleSpecDeactivated = DataSourceSpecStatus(zetapb.DataSourceSpec_STATUS_DEACTIVATED)
)

func (s DataSourceSpecStatus) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	status, ok := zetapb.DataSourceSpec_Status_name[int32(s)]
	if !ok {
		return buf, fmt.Errorf("unknown oracle spec value: %v", s)
	}
	return append(buf, []byte(status)...), nil
}

func (s *DataSourceSpecStatus) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zetapb.DataSourceSpec_Status_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown oracle spec status: %s", src)
	}
	*s = DataSourceSpecStatus(val)
	return nil
}

type LiquidityProvisionStatus zeta.LiquidityProvision_Status

const (
	LiquidityProvisionStatusUnspecified = LiquidityProvisionStatus(zeta.LiquidityProvision_STATUS_UNSPECIFIED)
	LiquidityProvisionStatusActive      = LiquidityProvisionStatus(zeta.LiquidityProvision_STATUS_ACTIVE)
	LiquidityProvisionStatusStopped     = LiquidityProvisionStatus(zeta.LiquidityProvision_STATUS_STOPPED)
	LiquidityProvisionStatusCancelled   = LiquidityProvisionStatus(zeta.LiquidityProvision_STATUS_CANCELLED)
	LiquidityProvisionStatusRejected    = LiquidityProvisionStatus(zeta.LiquidityProvision_STATUS_REJECTED)
	LiquidityProvisionStatusUndeployed  = LiquidityProvisionStatus(zeta.LiquidityProvision_STATUS_UNDEPLOYED)
	LiquidityProvisionStatusPending     = LiquidityProvisionStatus(zeta.LiquidityProvision_STATUS_PENDING)
)

func (s LiquidityProvisionStatus) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	status, ok := zeta.LiquidityProvision_Status_name[int32(s)]
	if !ok {
		return buf, fmt.Errorf("unknown liquidity provision status: %v", s)
	}
	return append(buf, []byte(status)...), nil
}

func (s *LiquidityProvisionStatus) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zeta.LiquidityProvision_Status_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown liquidity provision status: %s", src)
	}
	*s = LiquidityProvisionStatus(val)
	return nil
}

type StakeLinkingStatus eventspb.StakeLinking_Status

const (
	StakeLinkingStatusUnspecified = StakeLinkingStatus(eventspb.StakeLinking_STATUS_UNSPECIFIED)
	StakeLinkingStatusPending     = StakeLinkingStatus(eventspb.StakeLinking_STATUS_PENDING)
	StakeLinkingStatusAccepted    = StakeLinkingStatus(eventspb.StakeLinking_STATUS_ACCEPTED)
	StakeLinkingStatusRejected    = StakeLinkingStatus(eventspb.StakeLinking_STATUS_REJECTED)
)

func (s StakeLinkingStatus) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	status, ok := eventspb.StakeLinking_Status_name[int32(s)]
	if !ok {
		return buf, fmt.Errorf("unknown stake linking status: %v", s)
	}
	return append(buf, []byte(status)...), nil
}

func (s *StakeLinkingStatus) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := eventspb.StakeLinking_Status_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown stake linking status: %s", src)
	}
	*s = StakeLinkingStatus(val)
	return nil
}

type StakeLinkingType eventspb.StakeLinking_Type

const (
	StakeLinkingTypeUnspecified = StakeLinkingType(eventspb.StakeLinking_TYPE_UNSPECIFIED)
	StakeLinkingTypeLink        = StakeLinkingType(eventspb.StakeLinking_TYPE_LINK)
	StakeLinkingTypeUnlink      = StakeLinkingType(eventspb.StakeLinking_TYPE_UNLINK)
)

func (s StakeLinkingType) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	status, ok := eventspb.StakeLinking_Type_name[int32(s)]
	if !ok {
		return buf, fmt.Errorf("unknown stake linking type: %v", s)
	}
	return append(buf, []byte(status)...), nil
}

func (s *StakeLinkingType) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := eventspb.StakeLinking_Type_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown stake linking type: %s", src)
	}
	*s = StakeLinkingType(val)

	return nil
}

/************************* Node *****************************/

type NodeStatus zeta.NodeStatus

const (
	NodeStatusUnspecified  = NodeStatus(zeta.NodeStatus_NODE_STATUS_UNSPECIFIED)
	NodeStatusValidator    = NodeStatus(zeta.NodeStatus_NODE_STATUS_VALIDATOR)
	NodeStatusNonValidator = NodeStatus(zeta.NodeStatus_NODE_STATUS_NON_VALIDATOR)
)

func (ns NodeStatus) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	str, ok := zeta.NodeStatus_name[int32(ns)]
	if !ok {
		return buf, fmt.Errorf("unknown node status: %v", ns)
	}
	return append(buf, []byte(str)...), nil
}

func (ns *NodeStatus) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zeta.NodeStatus_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown node status: %s", src)
	}
	*ns = NodeStatus(val)
	return nil
}

type ValidatorNodeStatus zeta.ValidatorNodeStatus

const (
	ValidatorNodeStatusUnspecified = ValidatorNodeStatus(zeta.ValidatorNodeStatus_VALIDATOR_NODE_STATUS_UNSPECIFIED)
	ValidatorNodeStatusTendermint  = ValidatorNodeStatus(zeta.ValidatorNodeStatus_VALIDATOR_NODE_STATUS_TENDERMINT)
	ValidatorNodeStatusErsatz      = ValidatorNodeStatus(zeta.ValidatorNodeStatus_VALIDATOR_NODE_STATUS_ERSATZ)
	ValidatorNodeStatusPending     = ValidatorNodeStatus(zeta.ValidatorNodeStatus_VALIDATOR_NODE_STATUS_PENDING)
)

// ValidatorStatusRanking so we know which direction was a promotion and which was a demotion.
var ValidatorStatusRanking = map[ValidatorNodeStatus]int{
	ValidatorNodeStatusUnspecified: 0,
	ValidatorNodeStatusPending:     1,
	ValidatorNodeStatusErsatz:      2,
	ValidatorNodeStatusTendermint:  3,
}

func (ns ValidatorNodeStatus) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	str, ok := zeta.ValidatorNodeStatus_name[int32(ns)]
	if !ok {
		return buf, fmt.Errorf("unknown validator node status: %v", ns)
	}
	return append(buf, []byte(str)...), nil
}

func (ns *ValidatorNodeStatus) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zeta.ValidatorNodeStatus_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown validator node status: %s", src)
	}
	*ns = ValidatorNodeStatus(val)
	return nil
}

func (ns *ValidatorNodeStatus) UnmarshalJSON(src []byte) error {
	val, ok := zeta.ValidatorNodeStatus_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown validator node status: %s", src)
	}
	*ns = ValidatorNodeStatus(val)
	return nil
}

/************************* Position status  *****************************/

func (p PositionStatus) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	str, ok := zeta.PositionStatus_name[int32(p)]
	if !ok {
		return buf, fmt.Errorf("unknown position status: %v", p)
	}
	return append(buf, []byte(str)...), nil
}

func (p *PositionStatus) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := zeta.PositionStatus_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown position status: %s", string(src))
	}
	*p = PositionStatus(val)
	return nil
}

/************************* Protocol Upgrade *****************************/

type ProtocolUpgradeProposalStatus eventspb.ProtocolUpgradeProposalStatus

func (ps ProtocolUpgradeProposalStatus) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	str, ok := eventspb.ProtocolUpgradeProposalStatus_name[int32(ps)]
	if !ok {
		return buf, fmt.Errorf("unknown protocol upgrade proposal status: %v", ps)
	}
	return append(buf, []byte(str)...), nil
}

func (ps *ProtocolUpgradeProposalStatus) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	val, ok := eventspb.ProtocolUpgradeProposalStatus_value[string(src)]
	if !ok {
		return fmt.Errorf("unknown protocol upgrade proposal status: %s", src)
	}
	*ps = ProtocolUpgradeProposalStatus(val)
	return nil
}
