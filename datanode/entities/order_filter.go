package entities

import "github.com/zeta-protocol/zeta/protos/zeta"

type OrderFilter struct {
	Statuses         []zeta.Order_Status
	Types            []zeta.Order_Type
	TimeInForces     []zeta.Order_TimeInForce
	Reference        *string
	DateRange        *DateRange
	ExcludeLiquidity bool
	LiveOnly         bool
	PartyIDs         []string
	MarketIDs        []string
}
