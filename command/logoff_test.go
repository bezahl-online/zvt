package command

import (
	"testing"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/instr"
	"github.com/stretchr/testify/assert"
)

func TestLogOff(t *testing.T) {
	i := instr.Map["ACK"]
	want := Command{
		CtrlField: i,
		Data: apdu.DataUnit{
			Data:    []byte{},
			BMPOBJs: []bmp.OBJ{},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{},
			},
		},
	}
	err := PaymentTerminal.LogOff()
	got, err := PaymentTerminal.ReadResponse()
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
	}
}

func TestPT_LogOff(t *testing.T) {
	tests := []struct {
		name    string
		p       *PT
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.LogOff(); (err != nil) != tt.wantErr {
				t.Errorf("PT.LogOff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
