package command

import (
	"testing"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/bmp/blen"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/instr"
	"github.com/bezahl-online/zvt/util"
	"github.com/stretchr/testify/assert"
)

func TestEndOfDay(t *testing.T) {
	skipShort(t)
	// t.Cleanup(func() {
	// 	PaymentTerminal.Abort()
	// })
	err := PaymentTerminal.EndOfDay()
	assert.NoError(t, err)

}

func TestEndOfDayCompletion(t *testing.T) {
	skipShort(t)
	TestEndOfDay(t)
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
				Totals:  &testTotals,
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

func TestEndOfDayProcess3(t *testing.T) {
	testBytes, err := util.Load("testdata/1617777384987PT.hex")
	if !assert.NoError(t, err) {
		return
	}
	want := Command{
		CtrlField: instr.CtrlField{
			Class: 0x06,
			Instr: 0xD3,
			Length: blen.Length{
				Kind:  blen.BINARY,
				Size:  3,
				Value: 1069,
			},
			RawDataLength: 0,
		},
		Data: apdu.DataUnit{
			Data:    []byte{},
			BMPOBJs: []bmp.OBJ{},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{
					{TAG: []byte{0xD3}, Data: []byte{0x00}},
				},
			},
		},
	}
	got := Command{}
	if err = got.Unmarshal(&testBytes); assert.NoError(t, err) {
		if assert.EqualValues(t, want.CtrlField, got.CtrlField) &&
			assert.Equal(t, 2, len(got.Data.TLVContainer.Objects)) {
			got2 := EndOfDayResponse{}
			got2.Process(&got)
			assert.EqualValues(t, 3, got2.Transaction.Data.PrintOut.Type)
			assert.EqualValues(t, 924, len(got2.Transaction.Data.PrintOut.Text))
		}
	}

}
