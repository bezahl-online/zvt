package command

import (
	"testing"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/instr"
	"github.com/stretchr/testify/assert"
)

func TestCompletion(t *testing.T) {
	want := CompletionResponse{}
	got, err := PaymentTerminal.Completion()
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}

func TestProcess04ff(t *testing.T) {
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
	response := CompletionResponse{}
	response.process(result)
	assert.Equal(t, want, response.Message)
}
func TestProcess040f(t *testing.T) {
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
	want := CompletionResponse{
		Message: "",
		Transaction: &AuthResult{
			Error:  "",
			Result: "abort",
			Data:   &AuthResultData{},
		},
	}
	response := CompletionResponse{
		Status:      0,
		Message:     "",
		Transaction: &AuthResult{},
	}
	response.process(result)
	assert.Equal(t, want, response)
}
