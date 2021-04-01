package command

import (
	"reflect"
	"testing"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/instr"
	"github.com/stretchr/testify/assert"
)

func TestAuthorisation(t *testing.T) {
	skipShort(t)
	t.Cleanup(func() {
		PaymentTerminal.Abort()
	})
	config := &AuthConfig{
		Amount: 1,
	}
	err := PaymentTerminal.Authorisation(config)
	assert.NoError(t, err)
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
		Transaction: &AuthResult{
			Error:  "",
			Result: "pending",
			Data: &AuthResultData{
				Amount:     1,
				ReceiptNr:  22,
				TurnoverNr: 22,
				TraceNr:    22,
				Date:       "0308",
				Time:       "164923",
				TID:        "29001006",
				VU:         "100764992",
				AID:        "291675",
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

func TestPT_Authorisation(t *testing.T) {
	type args struct {
		config *AuthConfig
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
			if err := tt.p.Authorisation(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("PT.Authorisation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthConfig_marshal(t *testing.T) {
	tests := []struct {
		name string
		a    *AuthConfig
		want apdu.DataUnit
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.marshal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthConfig.marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorisationResponse_Process(t *testing.T) {
	type args struct {
		result *Command
	}
	tests := []struct {
		name    string
		r       *AuthorisationResponse
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Process(tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("AuthorisationResponse.Process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthResultData_FromOBJs(t *testing.T) {
	type args struct {
		objs []bmp.OBJ
	}
	tests := []struct {
		name       string
		r          *AuthResultData
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
				t.Errorf("AuthResultData.FromOBJs() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotError != tt.wantError {
				t.Errorf("AuthResultData.FromOBJs() gotError = %v, want %v", gotError, tt.wantError)
			}
		})
	}
}

func Test_formatPAN(t *testing.T) {
	type args struct {
		rawPAN []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatPAN(tt.args.rawPAN); got != tt.want {
				t.Errorf("formatPAN() = %v, want %v", got, tt.want)
			}
		})
	}
}
