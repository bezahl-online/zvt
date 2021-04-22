package command

import (
	"log"
	"testing"

	"github.com/bezahl-online/zvt/util"
	"github.com/stretchr/testify/assert"
)

func TestStatusCompletion(t *testing.T) {
	skipShort(t)
	for {
		TestStatus(t)
		got := StatusResponse{}
		err := PaymentTerminal.Completion(&got)
		if err != nil {
			log.Println(err.Error())
			assert.NoError(t, err)
			break
		}
		if got.Transaction != nil && got.Transaction.Result != Result_Pending {
			if got.Transaction.Result == Result_Success {
				// TODO assert result values
				_ = 0
			}
			//time.Sleep(defaultTimeout)
			break
		}
		// assert.EqualValues(t, want, got)
	}
}

func TestStatusProcess(t *testing.T) {
	testBytes, err := util.Load("testdata/1619070661833PT.hex")
	if !assert.NoError(t, err) {
		return
	}
	c := Command{}
	c.Unmarshal(&testBytes)
	want := StatusResponse{
		TransactionResponse: TransactionResponse{},
		Transaction: &StatusResult{
			Error:  "",
			Result: "success",
		},
	}
	got := StatusResponse{}
	got.Process(&c)
	assert.EqualValues(t, want, got)
}
