package command

import (
	"reflect"
	"testing"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/instr"
	"github.com/stretchr/testify/assert"
)

func TestDisplayText(t *testing.T) {
	skipShort(t)
	i := instr.Map["ACK"]
	want := Command{
		CtrlField: i,
		Data:      apdu.DataUnit{},
	}
	want.CtrlField.Length.Size = 1
	err := PaymentTerminal.DisplayText([]string{
		"Da steh ich nun,",
		"ich armer Tor,",
		"Und bin so klug",
		"als wie zuvor."})
	if assert.NoError(t, err) {
		got, err := PaymentTerminal.ReadResponse()
		if assert.NoError(t, err) {
			assert.EqualValues(t, want, *got)
		}
	}
}

func TestPT_DisplayText(t *testing.T) {
	type args struct {
		text []string
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
			if err := tt.p.DisplayText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("PT.DisplayText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_compileText(t *testing.T) {
	type args struct {
		textarray []string
	}
	tests := []struct {
		name string
		args args
		want apdu.DataUnit
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compileText(tt.args.textarray); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compileText() = %v, want %v", got, tt.want)
			}
		})
	}
}
