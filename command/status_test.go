package command

import (
	"log"
	"testing"

	"github.com/bezahl-online/zvt/util"
	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	if skipShort(t) {
		return
	}
	err := PaymentTerminal.Status()
	if !assert.NoError(t, err) {
		log.Fatal(err)
	}
}

func TestStatusCompletion(t *testing.T) {
	if skipShort(t) {
		return
	}
	TestStatus(t)
	for {
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
		TransactionResponse: TransactionResponse{
			Status:  0,
			Message: "card inserted",
		},
		Transaction: &StatusResult{
			Error:  "",
			Result: "success",
			Data: &StatusResultData{
				Date:   "",
				Time:   "",
				Status: 0xDC,
			},
		},
	}
	got := StatusResponse{}
	got.Process(&c)
	assert.EqualValues(t, want, got)
}
