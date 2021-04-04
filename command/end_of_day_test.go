package command

import (
	"testing"

	"github.com/bezahl-online/zvt/util"
	"github.com/stretchr/testify/assert"
)

func TestEndOfDay(t *testing.T) {
	skipShort(t)
	t.Cleanup(func() {
		PaymentTerminal.Abort()
	})
	err := PaymentTerminal.EndOfDay()
	assert.NoError(t, err)

}

var testTotals SingleTotals = SingleTotals{ // TODO: check with printout
	ReceiptNrStart: 68,
	ReceiptNrEnd:   71,
	CountEC:        4,
	TotalEC:        2600,
	CountJCB:       0,
	TotalJCB:       0,
	CountEurocard:  0,
	TotalEurocard:  0,
	CountAmex:      0,
	TotalAmex:      0,
	CountVisa:      0,
	TotalVisa:      0,
	CountDiners:    0,
	TotalDiners:    0,
	CountOther:     0,
	TotalOther:     0,
}

func TestSingleTotalsUnmarshal(t *testing.T) {
	testBytes, err := util.Load("testdata/1617181236803PT.hex")
	want := testTotals
	if !assert.NoError(t, err) {
		return
	}
	got := SingleTotals{}
	got.Unmarshal(testBytes[22:])
	assert.EqualValues(t, want, got)
}

func TestEndOfDayProcess(t *testing.T) {
	want := EndOfDayResponse{
		TransactionResponse: TransactionResponse{},
		Transaction: &EoDResult{
			Error:  "",
			Result: "pending",
			Data: &EoDResultData{
				TraceNr: 0,
				Date:    "0331",
				Time:    "105937",
				Total:   2600,
				Totals:  testTotals,
			},
		},
	}
	testBytes, err := util.Load("testdata/1617181236803PT.hex")
	if !assert.NoError(t, err) {
		return
	}
	c := Command{}
	c.Unmarshal(&testBytes)
	got := EndOfDayResponse{}
	got.Process(&c)
	assert.EqualValues(t, want, got)
}

func TestEndOfDayProcess2(t *testing.T) {
	testBytes, err := util.Load("testdata/1617177992186PT.hex")
	if !assert.NoError(t, err) {
		return
	}
	c := Command{}
	err = c.Unmarshal(&testBytes)
	assert.Error(t, err)
}
