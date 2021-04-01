package command

import (
	"testing"

	"github.com/bezahl-online/zvt/apdu/bmp"
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
	CountJCB:       11,
	TotalJCB:       0,
	CountEurocard:  18,
	TotalEurocard:  0,
	CountAmex:      25,
	TotalAmex:      0,
	CountVisa:      32,
	TotalVisa:      0,
	CountDiners:    39,
	TotalDiners:    0,
	CountOther:     46,
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
	want := EndOfDayResponse{
		TransactionResponse: TransactionResponse{},
		Transaction: &EoDResult{
			Error:  "",
			Result: "pending",
			Data: &EoDResultData{
				TraceNr: 0,
				Date:    "",
				Time:    "",
				Total:   0,
				Totals:  SingleTotals{},
			},
		},
	}
	testBytes, err := util.Load("testdata/1617177992186PT.hex")
	if !assert.NoError(t, err) {
		return
	}
	c := Command{}
	if err := c.Unmarshal(&testBytes); assert.NoError(t, err) {
		got := EndOfDayResponse{}
		got.Process(&c)
		assert.EqualValues(t, want, got)
	}
}

func TestSingleTotals_Unmarshal(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		s    *SingleTotals
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Unmarshal(tt.args.data)
		})
	}
}

func TestPT_EndOfDay(t *testing.T) {
	tests := []struct {
		name    string
		p       *PT
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.EndOfDay(); (err != nil) != tt.wantErr {
				t.Errorf("PT.EndOfDay() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEndOfDayResponse_Process(t *testing.T) {
	type args struct {
		result *Command
	}
	tests := []struct {
		name    string
		r       *EndOfDayResponse
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Process(tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("EndOfDayResponse.Process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEoDResultData_FromOBJs(t *testing.T) {
	type args struct {
		objs []bmp.OBJ
	}
	tests := []struct {
		name       string
		r          *EoDResultData
		args       args
		wantResult string
		wantError  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotError := tt.r.FromOBJs(tt.args.objs)
			if gotResult != tt.wantResult {
				t.Errorf("EoDResultData.FromOBJs() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotError != tt.wantError {
				t.Errorf("EoDResultData.FromOBJs() gotError = %v, want %v", gotError, tt.wantError)
			}
		})
	}
}
