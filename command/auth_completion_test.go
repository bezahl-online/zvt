package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndOfDayCompletion(t *testing.T) {
	skipShort(t)
	for {
		got := EndOfDayResponse{}
		err := PaymentTerminal.Completion(&got)
		if err != nil {
			assert.NoError(t, err)
			break
		}
		if got.Transaction != nil && got.Transaction.Result != Result_Pending {
			if got.Transaction.Result == Result_Success {
				// TODO assert result values
				_ = 0
			}
			break
		}
		// assert.EqualValues(t, want, got)
	}
}
