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

package gql

import (
	"context"
	"errors"
	"strconv"

	"github.com/zeta-protocol/zeta/protos/zeta"
)

type updateMarketConfigurationResolver ZetaResolverRoot

func (r *updateMarketConfigurationResolver) Instrument(ctx context.Context,
	obj *zeta.UpdateMarketConfiguration,
) (*UpdateInstrumentConfiguration, error) {
	if obj == nil {
		return nil, errors.New("no market configuration update provided")
	}
	protoInstrument := obj.Instrument

	var product *zeta.UpdateFutureProduct

	switch p := protoInstrument.Product.(type) {
	case *zeta.UpdateInstrumentConfiguration_Future:
		product = &zeta.UpdateFutureProduct{
			QuoteName:                           p.Future.QuoteName,
			DataSourceSpecForSettlementData:     p.Future.DataSourceSpecForSettlementData,
			DataSourceSpecForTradingTermination: p.Future.DataSourceSpecForTradingTermination,
			DataSourceSpecBinding:               p.Future.DataSourceSpecBinding,
		}
	default:
		return nil, ErrUnsupportedProduct
	}

	updateInstrumentConfiguration := &UpdateInstrumentConfiguration{
		Code:    protoInstrument.Code,
		Product: product,
	}

	return updateInstrumentConfiguration, nil
}

func (r *updateMarketConfigurationResolver) PriceMonitoringParameters(ctx context.Context,
	obj *zeta.UpdateMarketConfiguration,
) (*PriceMonitoringParameters, error) {
	if obj == nil {
		return nil, errors.New("no market configuration update provided")
	}

	if obj.PriceMonitoringParameters == nil {
		return nil, nil
	}

	triggers := make([]*PriceMonitoringTrigger, 0, len(obj.PriceMonitoringParameters.Triggers))

	for _, trigger := range obj.PriceMonitoringParameters.Triggers {
		probability, err := strconv.ParseFloat(trigger.Probability, 64)
		if err != nil {
			continue
		}
		triggers = append(triggers, &PriceMonitoringTrigger{
			HorizonSecs:          int(trigger.Horizon),
			Probability:          probability,
			AuctionExtensionSecs: int(trigger.AuctionExtension),
		})
	}

	params := &PriceMonitoringParameters{
		Triggers: triggers,
	}

	return params, nil
}

func (r *updateMarketConfigurationResolver) LiquidityMonitoringParameters(ctx context.Context,
	obj *zeta.UpdateMarketConfiguration,
) (*LiquidityMonitoringParameters, error) {
	if obj == nil {
		return nil, errors.New("no market configuration update provided")
	}

	if obj.LiquidityMonitoringParameters == nil {
		return nil, nil
	}

	return &LiquidityMonitoringParameters{
		TargetStakeParameters: &TargetStakeParameters{
			TimeWindow:    int(obj.LiquidityMonitoringParameters.TargetStakeParameters.TimeWindow),
			ScalingFactor: obj.LiquidityMonitoringParameters.TargetStakeParameters.ScalingFactor,
		},
		TriggeringRatio: obj.LiquidityMonitoringParameters.TriggeringRatio,
	}, nil
}

func (r *updateMarketConfigurationResolver) RiskParameters(ctx context.Context,
	obj *zeta.UpdateMarketConfiguration,
) (UpdateMarketRiskParameters, error) {
	if obj == nil {
		return nil, errors.New("no market configuration update provided")
	}

	if obj.RiskParameters == nil {
		return nil, errors.New("no risk configuration provided")
	}

	var params UpdateMarketRiskParameters

	switch rp := obj.RiskParameters.(type) {
	case *zeta.UpdateMarketConfiguration_Simple:
		params = rp
	case *zeta.UpdateMarketConfiguration_LogNormal:
		params = rp
	default:
		return nil, errors.New("invalid risk configuration provided")
	}

	return params, nil
}
