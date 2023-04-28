package entities

import (
	"time"

	"github.com/zeta-protocol/zeta/libs/ptr"
	v2 "github.com/zeta-protocol/zeta/protos/data-node/api/v2"
)

type DateRange struct {
	Start *time.Time
	End   *time.Time
}

func DateRangeFromProto(dateRangeInput *v2.DateRange) (dateRange DateRange) {
	if dateRangeInput == nil {
		return
	}

	if dateRangeInput.StartTimestamp != nil {
		dateRange.Start = ptr.From(time.Unix(0, *dateRangeInput.StartTimestamp))
	}

	if dateRangeInput.EndTimestamp != nil {
		dateRange.End = ptr.From(time.Unix(0, *dateRangeInput.EndTimestamp))
	}

	return
}
