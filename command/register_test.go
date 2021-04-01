package command

import (
	"reflect"
	"testing"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/config"
	"github.com/bezahl-online/zvt/instr"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	skipShort(t) // start
	configByte := config.PaymentReceiptPrintedByECR +
		config.AdminReceiptPrintedByECR +
		config.PTSendsIntermediateStatus +
		config.ECRusingPrintLinesForPrintout
	serviceByte := config.ServiceMenuNOTAssignedToFunctionKey +
		config.DisplayTextsForCommandsAuthorisation
	var msgSquID *tlv.DataObject = &tlv.DataObject{
		TAG:  []byte{0x1F, 0x73},
		Data: []byte{0, 0, 0},
	}

	var listOfCommands *tlv.DataObject = &tlv.DataObject{
		TAG:  []byte{0x26},
		Data: []byte{0x0A, 0x02, 0x06, 0xD3},
	}
	var tlvContainer *tlv.Container = &tlv.Container{
		Objects: []tlv.DataObject{},
	}
	tlvContainer.Objects = append(tlvContainer.Objects, *listOfCommands, *msgSquID)
	i := instr.Map["ACK"]
	i.Length.Size = 1
	want := Command{
		CtrlField: i,
		Data:      apdu.DataUnit{},
	}
	err := PaymentTerminal.Register(&Config{
		pwd:          fixedPassword,
		config:       byte(configByte),
		currency:     EUR,
		service:      byte(serviceByte),
		tlvContainer: tlvContainer,
	})
	got, err := PaymentTerminal.ReadResponse()
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
		// completion
		i := instr.Map["Completion"]
		want := Command{
			CtrlField: i,
			Data: apdu.DataUnit{
				Data: []byte{},
				BMPOBJs: []bmp.OBJ{
					{ID: 0x19, Data: []byte{0}, Size: 2},
					{ID: 0x29, Data: []byte{0x29, 0x00, 0x10, 0x06}, Size: 5},
					{ID: 0x49, Data: []byte{0x09, 0x78}, Size: 3},
				},
				TLVContainer: tlv.Container{
					Objects: []tlv.DataObject{},
				},
			},
		}
		got, err = PaymentTerminal.ReadResponse()
		if assert.NoError(t, err) {
			if assert.EqualValues(t, want, *got) {
				PaymentTerminal.SendACK()
			}
		}

	}

}

func TestPT_Register(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		p       *PT
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Register(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("PT.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_CompileConfig(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want apdu.DataUnit
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.CompileConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.CompileConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
