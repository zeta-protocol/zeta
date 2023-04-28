package snapshot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractIntervalFromViewDefinition(t *testing.T) {
	viewDefinition := ` SELECT balances.account_id,
	time_bucket('01:00:00'::interval, balances.zeta_time) AS bucket,
		last(balances.balance, balances.zeta_time) AS balance,
		last(balances.tx_hash, balances.zeta_time) AS tx_hash,
		last(balances.zeta_time, balances.zeta_time) AS zeta_time
	FROM balances
	GROUP BY balances.account_id, (time_bucket('01:00:00'::interval, balances.zeta_time));`

	interval, err := extractIntervalFromViewDefinition(viewDefinition)
	require.NoError(t, err)
	assert.Equal(t, "01:00:00", interval)

	viewDefinition = ` SELECT balances.account_id,
	time_bucket('1 day'::interval, balances.zeta_time) AS bucket,
		last(balances.balance, balances.zeta_time) AS balance,
		last(balances.tx_hash, balances.zeta_time) AS tx_hash,
		last(balances.zeta_time, balances.zeta_time) AS zeta_time
	FROM balances
	GROUP BY balances.account_id, (time_bucket('1 day'::interval, balances.zeta_time));`

	interval, err = extractIntervalFromViewDefinition(viewDefinition)
	require.NoError(t, err)
	assert.Equal(t, "1 day", interval)
}
