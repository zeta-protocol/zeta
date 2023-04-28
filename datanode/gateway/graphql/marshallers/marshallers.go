package marshallers

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	v2 "github.com/zeta-protocol/zeta/protos/data-node/api/v2"
	"github.com/zeta-protocol/zeta/protos/zeta"
	zetapb "github.com/zeta-protocol/zeta/protos/zeta"
	commandspb "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"
	datapb "github.com/zeta-protocol/zeta/protos/zeta/data/v1"
	eventspb "github.com/zeta-protocol/zeta/protos/zeta/events/v1"

	"github.com/99designs/gqlgen/graphql"
)

var ErrUnimplemented = errors.New("unmarshaller not implemented as this API is query only")

func MarshalAccountType(t zeta.AccountType) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(t.String())))
	})
}

func UnmarshalAccountType(v interface{}) (zeta.AccountType, error) {
	s, ok := v.(string)
	if !ok {
		return zeta.AccountType_ACCOUNT_TYPE_UNSPECIFIED, fmt.Errorf("expected account type to be a string")
	}

	t, ok := zeta.AccountType_value[s]
	if !ok {
		return zeta.AccountType_ACCOUNT_TYPE_UNSPECIFIED, fmt.Errorf("failed to convert AccountType from GraphQL to Proto: %v", s)
	}

	return zeta.AccountType(t), nil
}

func MarshalSide(s zeta.Side) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalSide(v interface{}) (zeta.Side, error) {
	s, ok := v.(string)
	if !ok {
		return zeta.Side_SIDE_UNSPECIFIED, fmt.Errorf("expected account type to be a string")
	}

	side, ok := zeta.Side_value[s]
	if !ok {
		return zeta.Side_SIDE_UNSPECIFIED, fmt.Errorf("failed to convert AccountType from GraphQL to Proto: %v", s)
	}

	return zeta.Side(side), nil
}

func MarshalProposalState(s zeta.Proposal_State) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalProposalState(v interface{}) (zeta.Proposal_State, error) {
	s, ok := v.(string)
	if !ok {
		return zeta.Proposal_STATE_UNSPECIFIED, fmt.Errorf("expected proposal state to be a string")
	}

	side, ok := zeta.Proposal_State_value[s]
	if !ok {
		return zeta.Proposal_STATE_UNSPECIFIED, fmt.Errorf("failed to convert ProposalState from GraphQL to Proto: %v", s)
	}

	return zeta.Proposal_State(side), nil
}

func MarshalTransferType(t zeta.TransferType) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(t.String())))
	})
}

func UnmarshalTransferType(v interface{}) (zeta.TransferType, error) {
	s, ok := v.(string)
	if !ok {
		return zeta.TransferType_TRANSFER_TYPE_UNSPECIFIED, fmt.Errorf("expected transfer type to be a string")
	}

	t, ok := zeta.TransferType_value[s]
	if !ok {
		return zeta.TransferType_TRANSFER_TYPE_UNSPECIFIED, fmt.Errorf("failed to convert TransferType from GraphQL to Proto: %v", s)
	}

	return zeta.TransferType(t), nil
}

func MarshalTransferStatus(s eventspb.Transfer_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalTransferStatus(v interface{}) (eventspb.Transfer_Status, error) {
	return eventspb.Transfer_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalDispatchMetric(s zeta.DispatchMetric) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalDispatchMetric(v interface{}) (zeta.DispatchMetric, error) {
	return zeta.DispatchMetric_DISPATCH_METRIC_UNSPECIFIED, ErrUnimplemented
}

func MarshalNodeStatus(s zeta.NodeStatus) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalNodeStatus(v interface{}) (zeta.NodeStatus, error) {
	return zeta.NodeStatus_NODE_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalAssetStatus(s zeta.Asset_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalAssetStatus(v interface{}) (zeta.Asset_Status, error) {
	return zeta.Asset_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalNodeSignatureKind(s commandspb.NodeSignatureKind) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalNodeSignatureKind(v interface{}) (commandspb.NodeSignatureKind, error) {
	return commandspb.NodeSignatureKind_NODE_SIGNATURE_KIND_UNSPECIFIED, ErrUnimplemented
}

func MarshalOracleSpecStatus(s zetapb.DataSourceSpec_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalOracleSpecStatus(v interface{}) (zetapb.DataSourceSpec_Status, error) {
	return zetapb.DataSourceSpec_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalPropertyKeyType(s datapb.PropertyKey_Type) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalPropertyKeyType(v interface{}) (datapb.PropertyKey_Type, error) {
	return datapb.PropertyKey_TYPE_UNSPECIFIED, ErrUnimplemented
}

func MarshalConditionOperator(s datapb.Condition_Operator) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalConditionOperator(v interface{}) (datapb.Condition_Operator, error) {
	return datapb.Condition_OPERATOR_UNSPECIFIED, ErrUnimplemented
}

func MarshalVoteValue(s zeta.Vote_Value) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalVoteValue(v interface{}) (zeta.Vote_Value, error) {
	return zeta.Vote_VALUE_UNSPECIFIED, ErrUnimplemented
}

func MarshalAuctionTrigger(s zeta.AuctionTrigger) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalAuctionTrigger(v interface{}) (zeta.AuctionTrigger, error) {
	return zeta.AuctionTrigger_AUCTION_TRIGGER_UNSPECIFIED, ErrUnimplemented
}

func MarshalStakeLinkingStatus(s eventspb.StakeLinking_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalStakeLinkingStatus(v interface{}) (eventspb.StakeLinking_Status, error) {
	return eventspb.StakeLinking_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalStakeLinkingType(s eventspb.StakeLinking_Type) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalStakeLinkingType(v interface{}) (eventspb.StakeLinking_Type, error) {
	return eventspb.StakeLinking_TYPE_UNSPECIFIED, ErrUnimplemented
}

func MarshalWithdrawalStatus(s zeta.Withdrawal_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalWithdrawalStatus(v interface{}) (zeta.Withdrawal_Status, error) {
	return zeta.Withdrawal_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalDepositStatus(s zeta.Deposit_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalDepositStatus(v interface{}) (zeta.Deposit_Status, error) {
	return zeta.Deposit_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalOrderStatus(s zeta.Order_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalOrderStatus(v interface{}) (zeta.Order_Status, error) {
	s, ok := v.(string)
	if !ok {
		return zeta.Order_STATUS_UNSPECIFIED, fmt.Errorf("exoected order status to be a string")
	}

	t, ok := zeta.Order_Status_value[s]
	if !ok {
		return zeta.Order_STATUS_UNSPECIFIED, fmt.Errorf("failed to convert order status from GraphQL to Proto: %v", s)
	}

	return zeta.Order_Status(t), nil
}

func MarshalOrderTimeInForce(s zeta.Order_TimeInForce) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalOrderTimeInForce(v interface{}) (zeta.Order_TimeInForce, error) {
	s, ok := v.(string)
	if !ok {
		return zeta.Order_TIME_IN_FORCE_UNSPECIFIED, fmt.Errorf("expected order time in force to be a string")
	}

	t, ok := zeta.Order_TimeInForce_value[s]
	if !ok {
		return zeta.Order_TIME_IN_FORCE_UNSPECIFIED, fmt.Errorf("failed to convert TimeInForce from GraphQL to Proto: %v", s)
	}

	return zeta.Order_TimeInForce(t), nil
}

func MarshalPeggedReference(s zeta.PeggedReference) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalPeggedReference(v interface{}) (zeta.PeggedReference, error) {
	return zeta.PeggedReference_PEGGED_REFERENCE_UNSPECIFIED, ErrUnimplemented
}

func MarshalProposalRejectionReason(s zeta.ProposalError) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalProposalRejectionReason(v interface{}) (zeta.ProposalError, error) {
	return zeta.ProposalError_PROPOSAL_ERROR_UNSPECIFIED, ErrUnimplemented
}

func MarshalOrderRejectionReason(s zeta.OrderError) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalOrderRejectionReason(v interface{}) (zeta.OrderError, error) {
	return zeta.OrderError_ORDER_ERROR_UNSPECIFIED, ErrUnimplemented
}

func MarshalOrderType(s zeta.Order_Type) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalOrderType(v interface{}) (zeta.Order_Type, error) {
	s, ok := v.(string)
	if !ok {
		return zeta.Order_TYPE_UNSPECIFIED, fmt.Errorf("expected order type to be a string")
	}

	t, ok := zeta.Order_Type_value[s]
	if !ok {
		return zeta.Order_TYPE_UNSPECIFIED, fmt.Errorf("failed to convert OrderType from GraphQL to Proto: %v", s)
	}

	return zeta.Order_Type(t), nil
}

func MarshalMarketState(s zeta.Market_State) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalMarketState(v interface{}) (zeta.Market_State, error) {
	return zeta.Market_STATE_UNSPECIFIED, ErrUnimplemented
}

func MarshalMarketTradingMode(s zeta.Market_TradingMode) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalMarketTradingMode(v interface{}) (zeta.Market_TradingMode, error) {
	return zeta.Market_TRADING_MODE_UNSPECIFIED, ErrUnimplemented
}

func MarshalInterval(s zeta.Interval) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalInterval(v interface{}) (zeta.Interval, error) {
	s, ok := v.(string)
	if !ok {
		return zeta.Interval_INTERVAL_UNSPECIFIED, fmt.Errorf("expected interval in force to be a string")
	}

	t, ok := zeta.Interval_value[s]
	if !ok {
		return zeta.Interval_INTERVAL_UNSPECIFIED, fmt.Errorf("failed to convert Interval from GraphQL to Proto: %v", s)
	}

	return zeta.Interval(t), nil
}

func MarshalProposalType(s v2.ListGovernanceDataRequest_Type) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalProposalType(v interface{}) (v2.ListGovernanceDataRequest_Type, error) {
	s, ok := v.(string)
	if !ok {
		return v2.ListGovernanceDataRequest_TYPE_UNSPECIFIED, fmt.Errorf("expected proposal type in force to be a string")
	}

	t, ok := v2.ListGovernanceDataRequest_Type_value[s]
	if !ok {
		return v2.ListGovernanceDataRequest_TYPE_UNSPECIFIED, fmt.Errorf("failed to convert proposal type from GraphQL to Proto: %v", s)
	}

	return v2.ListGovernanceDataRequest_Type(t), nil
}

func MarshalLiquidityProvisionStatus(s zeta.LiquidityProvision_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalLiquidityProvisionStatus(v interface{}) (zeta.LiquidityProvision_Status, error) {
	return zeta.LiquidityProvision_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalTradeType(s zeta.Trade_Type) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalTradeType(v interface{}) (zeta.Trade_Type, error) {
	return zeta.Trade_TYPE_UNSPECIFIED, ErrUnimplemented
}

func MarshalValidatorStatus(s zeta.ValidatorNodeStatus) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalValidatorStatus(v interface{}) (zeta.ValidatorNodeStatus, error) {
	return zeta.ValidatorNodeStatus_VALIDATOR_NODE_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalProtocolUpgradeProposalStatus(s eventspb.ProtocolUpgradeProposalStatus) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalProtocolUpgradeProposalStatus(v interface{}) (eventspb.ProtocolUpgradeProposalStatus, error) {
	s, ok := v.(string)
	if !ok {
		return eventspb.ProtocolUpgradeProposalStatus_PROTOCOL_UPGRADE_PROPOSAL_STATUS_UNSPECIFIED, fmt.Errorf("expected proposal type in force to be a string")
	}

	t, ok := eventspb.ProtocolUpgradeProposalStatus_value[s] // v2.ListGovernanceDataRequest_Type_value[s]
	if !ok {
		return eventspb.ProtocolUpgradeProposalStatus_PROTOCOL_UPGRADE_PROPOSAL_STATUS_UNSPECIFIED, fmt.Errorf("failed to convert proposal type from GraphQL to Proto: %v", s)
	}

	return eventspb.ProtocolUpgradeProposalStatus(t), nil
}

func MarshalPositionStatus(s zeta.PositionStatus) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalPositionStatus(v interface{}) (zeta.PositionStatus, error) {
	s, ok := v.(string)
	if !ok {
		return zeta.PositionStatus_POSITION_STATUS_UNSPECIFIED, fmt.Errorf("expected position status to be a string")
	}
	t, ok := zeta.PositionStatus_value[s]
	if !ok {
		return zeta.PositionStatus_POSITION_STATUS_UNSPECIFIED, fmt.Errorf("failed to convert position status to Proto: %v", s)
	}
	return zeta.PositionStatus(t), nil
}
