package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorisationCompletion(t *testing.T) {
	skipShort(t)
	for {
		got := AuthorisationResponse{}
		err := PaymentTerminal.Completion(&got)
		if err != nil {
			assert.NoError(t, err)
			break
		}
		if got.Transaction != nil && got.Transaction.Result != Result_Pending {
			if got.Transaction.Result == Result_Success {
				// TODO assert result values
			}
			break
		}
		// assert.EqualValues(t, want, got)
	}
}
