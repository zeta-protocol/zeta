// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.ZETA file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package nullchain

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/zeta-protocol/zeta/core/examples/nullchain/config"
	"github.com/zeta-protocol/zeta/core/types"
	vgrand "github.com/zeta-protocol/zeta/libs/rand"
	"github.com/zeta-protocol/zeta/protos/zeta"
	v1 "github.com/zeta-protocol/zeta/protos/zeta/commands/v1"
	datav1 "github.com/zeta-protocol/zeta/protos/zeta/data/v1"
	walletpb "github.com/zeta-protocol/zeta/protos/zeta/wallet/v1"
)

func MarketProposalTxn(now time.Time, oraclePubkey string) (*walletpb.SubmitTransactionRequest, string) {
	reference := "ref-" + vgrand.RandomStr(10)
	asset := config.NormalAsset

	pubKey := types.CreateSignerFromString(oraclePubkey, types.DataSignerTypePubKey)
	cmd := &walletpb.SubmitTransactionRequest_ProposalSubmission{
		ProposalSubmission: &v1.ProposalSubmission{
			Reference: reference,
			Terms: &zeta.ProposalTerms{
				ValidationTimestamp: now.Add(2 * time.Second).Unix(),
				ClosingTimestamp:    now.Add(10 * time.Second).Unix(),
				EnactmentTimestamp:  now.Add(15 * time.Second).Unix(),
				Change: &zeta.ProposalTerms_NewMarket{
					NewMarket: &zeta.NewMarket{
						Changes: &zeta.NewMarketConfiguration{
							Instrument: &zeta.InstrumentConfiguration{
								Code: "CRYPTO:BTCUSD/NOV21",
								Name: "NOV 2021 BTC vs USD future",
								Product: &zeta.InstrumentConfiguration_Future{
									Future: &zeta.FutureProduct{
										SettlementAsset: asset,
										QuoteName:       "BTCUSD",
										DataSourceSpecForSettlementData: &zeta.DataSourceDefinition{
											SourceType: &zeta.DataSourceDefinition_External{
												External: &zeta.DataSourceDefinitionExternal{
													SourceType: &zeta.DataSourceDefinitionExternal_Oracle{
														Oracle: &zeta.DataSourceSpecConfiguration{
															Signers: []*datav1.Signer{pubKey.IntoProto()},
															Filters: []*datav1.Filter{
																{
																	Key: &datav1.PropertyKey{
																		Name: "prices." + asset + ".value",
																		Type: datav1.PropertyKey_TYPE_INTEGER,
																	},
																	Conditions: []*datav1.Condition{},
																},
															},
														},
													},
												},
											},
										},
										DataSourceSpecForTradingTermination: &zeta.DataSourceDefinition{
											SourceType: &zeta.DataSourceDefinition_External{
												External: &zeta.DataSourceDefinitionExternal{
													SourceType: &zeta.DataSourceDefinitionExternal_Oracle{
														Oracle: &zeta.DataSourceSpecConfiguration{
															Signers: []*datav1.Signer{pubKey.IntoProto()},
															Filters: []*datav1.Filter{
																{
																	Key: &datav1.PropertyKey{
																		Name: "trading.termination",
																		Type: datav1.PropertyKey_TYPE_BOOLEAN,
																	},
																	Conditions: []*datav1.Condition{},
																},
															},
														},
													},
												},
											},
										},
										DataSourceSpecBinding: &zeta.DataSourceSpecToFutureBinding{
											SettlementDataProperty:     "prices." + asset + ".value",
											TradingTerminationProperty: "trading.termination",
										},
									},
								},
							},
							DecimalPlaces: 5,
							Metadata:      []string{"base:BTC", "quote:USD", "class:fx/crypto", "monthly", "sector:crypto"},
							RiskParameters: &zeta.NewMarketConfiguration_Simple{
								Simple: &zeta.SimpleModelParams{
									FactorLong:           0.15,
									FactorShort:          0.25,
									MaxMoveUp:            10,
									MinMoveDown:          -5,
									ProbabilityOfTrading: 0.1,
								},
							},
							LpPriceRange:            "0.95",
							LinearSlippageFactor:    "0.1",
							QuadraticSlippageFactor: "0.1",
						},
					},
				},
			},
		},
	}

	return &walletpb.SubmitTransactionRequest{
		Command: cmd,
	}, reference
}

func VoteTxn(proposalID string, vote zeta.Vote_Value) *walletpb.SubmitTransactionRequest {
	return &walletpb.SubmitTransactionRequest{
		Command: &walletpb.SubmitTransactionRequest_VoteSubmission{
			VoteSubmission: &v1.VoteSubmission{
				ProposalId: proposalID,
				Value:      vote,
			},
		},
	}
}

func OrderTxn(
	marketId string,
	price, size uint64,
	side zeta.Side,
	orderT zeta.Order_Type,
	expiresAt time.Time,
) *walletpb.SubmitTransactionRequest {
	cmd := &walletpb.SubmitTransactionRequest_OrderSubmission{
		OrderSubmission: &v1.OrderSubmission{
			MarketId:    marketId,
			Price:       strconv.FormatUint(price, 10),
			Size:        size,
			Side:        side,
			Type:        orderT,
			TimeInForce: zeta.Order_TIME_IN_FORCE_GTT,
			ExpiresAt:   expiresAt.UnixNano(),
		},
	}

	return &walletpb.SubmitTransactionRequest{
		Command: cmd,
	}
}

func OracleTxn(key, value string) *walletpb.SubmitTransactionRequest {
	data := map[string]string{
		key: value,
	}

	b, _ := json.Marshal(data)

	cmd := &walletpb.SubmitTransactionRequest_OracleDataSubmission{
		OracleDataSubmission: &v1.OracleDataSubmission{
			Source:  v1.OracleDataSubmission_ORACLE_SOURCE_JSON,
			Payload: b,
		},
	}

	return &walletpb.SubmitTransactionRequest{
		Command: cmd,
	}
}
