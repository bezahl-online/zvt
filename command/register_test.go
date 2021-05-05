package command

import (
	"testing"

	"github.com/bezahl-online/zvt/instr"
	"github.com/bezahl-online/zvt/util"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	if skipShort(t) {
		return
	} // start
	i := instr.Map["ACK"]
	i.Length.Size = 1
	err := PaymentTerminal.Register()
	assert.NoError(t, err)
}

func TestRegisterCompletion(t *testing.T) {
	if skipShort(t) {
		return
	}
	TestRegister(t)
	for {
		got := RegisterResponse{}
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

func TestRegisterProcess1(t *testing.T) {
	testBytes, err := util.Load("testdata/1620218414428PT.hex")
	if !assert.NoError(t, err) {
		return
	}
	c := Command{}
	c.Unmarshal(&testBytes)
	want := RegisterResponse{
		TransactionResponse: TransactionResponse{Status: 0, Message: ""},
		Transaction: &RegisterResult{
			Error:  "",
			Result: "success",
		},
	}
	got := RegisterResponse{}
	got.Process(&c)
	assert.EqualValues(t, want, got)
}
