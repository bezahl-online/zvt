package command

import (
	"testing"

	"github.com/bezahl-online/zvt/util"
	"github.com/stretchr/testify/assert"
)

func TestAuthorisationCompletion(t *testing.T) {
	skipShort(t)
	TestAuthorisation(t)
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
				_ = 0
			}
			break
		}
		// assert.EqualValues(t, want, got)
	}
}

func TestAuthorisationProcess(t *testing.T) {
	testBytes, err := util.Load("testdata/1618046758827PT.hex")
	if !assert.NoError(t, err) {
		return
	}
	c := Command{}
	c.Unmarshal(&testBytes)
	want := AuthorisationResponse{
		TransactionResponse: TransactionResponse{
			Status:  0,
			Message: "BZT sendet Autostorno",
		},
		Transaction: &AuthResult{
			Error:  "",
			Result: "pending",
			Data: &AuthResultData{
				Amount:      0,
				Currency:    0,
				ReceiptNr:   0,
				TurnoverNr:  0,
				TraceNr:     0,
				Date:        "",
				Time:        "",
				TID:         "29001006",
				VU:          "",
				AID:         "",
				Info:        "Karte abgelehnt",
				PaymentType: 0,
				Card: CardData{
					Name:  "",
					Type:  13,
					PAN:   "",
					Tech:  0,
					SeqNr: 0,
				},
			},
		},
	}
	got := AuthorisationResponse{}
	got.Process(&c)
	assert.EqualValues(t, want, got)
}
