package command

import (
	"log"
	"testing"
	"time"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/instr"
	"github.com/bezahl-online/zvt/util"
	"github.com/stretchr/testify/assert"
)

func TestAuthorisation(t *testing.T) {
	skipShort(t)
	t.Cleanup(func() {
		time.Sleep(time.Second)
		// PaymentTerminal.Abort()
	})
	config := &AuthConfig{
		Amount: 1,
	}
	err := PaymentTerminal.Authorisation(config)
	assert.NoError(t, err)
}
func TestAuthorisationCompletion(t *testing.T) {
	skipShort(t)
	TestAuthorisation(t)
	if t.Failed() {
		return
	}
	for {
		got := AuthorisationResponse{}
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
			Status:  5,
			Message: "Karte abgelehnt",
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

func TestFormatPAN(t *testing.T) {
	want := "XXXX XXXX XXXX 5726"
	data := []uint8{0xee, 0xee, 0xee, 0xee, 0xee,
		0xee, 0x57, 0x26}
	got := formatPAN(data)
	assert.Equal(t, want, got)
	want = "XXXX XXXX XXXX 572"
	data = []uint8{0xee, 0xee, 0xee, 0xee, 0xee,
		0xee, 0x57, 0x2F}
	got = formatPAN(data)
	assert.Equal(t, want, got)
}

func TestAuthProcess04ff(t *testing.T) {
	result := &Command{
		CtrlField: instr.Map["Intermediate"],
		Data: apdu.DataUnit{
			Data: []byte{1},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{
					{
						TAG: []byte{0x24},
						Data: []byte{
							0x07, 0x26,
							0x4B, 0x61, 0x72, 0x74, 0x65, 0x20, 0x76, 0x6F,
							0x72, 0x68, 0x61, 0x6C, 0x74, 0x65, 0x6E, 0x0A,
							0x65, 0x69, 0x6E, 0x73, 0x74, 0x65, 0x63, 0x6B,
							0x65, 0x6E, 0x2F, 0x64, 0x75, 0x72, 0x63, 0x68,
							0x7A, 0x69, 0x65, 0x68, 0x65, 0x6E}},
				},
			},
		},
	}
	want := "Karte vorhalten\neinstecken/durchziehen"
	response := AuthorisationResponse{}
	response.Process(result)
	assert.Equal(t, want, response.Message)
}

func TestAuthProcess040f(t *testing.T) {
	result := &Command{
		CtrlField: instr.Map["StatusInformation"],
		Data: apdu.DataUnit{
			Data: []byte{1},
			BMPOBJs: []bmp.OBJ{
				{ID: 0x27, Data: []byte{0x6C}},
			},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{
					{
						TAG:  []byte{0x29},
						Data: []byte{0x29, 0x00, 0x10, 0x06},
					},
				},
			},
		},
	}
	want := AuthorisationResponse{
		Transaction: &AuthResult{
			Error:  "",
			Result: "pending",
			Data:   &AuthResultData{},
		},
	}
	response := AuthorisationResponse{
		Transaction: &AuthResult{},
	}
	response.Process(result)
	assert.Equal(t, want, response)
}

func TestAuthProcess040f_2(t *testing.T) {
	result := &Command{
		CtrlField: instr.CtrlField{
			Class: 0x04,
			Instr: 0x0F,
		},
		Data: apdu.DataUnit{
			Data:         []byte{},
			BMPOBJs:      Objects,
			TLVContainer: tlv.Container{},
		},
	}

	want := AuthorisationResponse{
		TransactionResponse: TransactionResponse{
			Status:  0,
			Message: "OK",
		},
		Transaction: &AuthResult{
			Error:  "",
			Result: "pending",
			Data: &AuthResultData{
				Amount:      1,
				Currency:    978,
				ReceiptNr:   22,
				TurnoverNr:  22,
				TraceNr:     22,
				Date:        "0308",
				Time:        "164923",
				TID:         "29001006",
				VU:          "100764992",
				AID:         "291675",
				Info:        "GEN.NR.:291675",
				PaymentType: 112,
				Card: CardData{
					Name:  "Debit Mastercard",
					Type:  46,
					PAN:   "XXXX XXXX XXXX 5726",
					Tech:  0,
					SeqNr: 1,
				},
			}},
	}
	response := AuthorisationResponse{
		Transaction: &AuthResult{},
	}
	response.Process(result)
	assert.Equal(t, want, response)
}
