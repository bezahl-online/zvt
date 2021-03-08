package zvt

import (
	"testing"
	"time"

	"bezahl.online/zvt/src/apdu"
	"bezahl.online/zvt/src/apdu/bmp"
	"bezahl.online/zvt/src/instr"
	"bezahl.online/zvt/src/zvt/config"
	"bezahl.online/zvt/src/zvt/tlv"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	// start
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
	want := Command{
		Instr: i,
		Data: apdu.DataUnit{
			Data:    []byte{},
			BMPOBJs: []bmp.OBJ{},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{},
			},
		},
	}
	err := ZVT.Register(&Config{
		pwd:          fixedPassword,
		config:       byte(configByte),
		currency:     EUR,
		service:      byte(serviceByte),
		tlvContainer: tlvContainer,
	})
	got, err := ZVT.ReadResponse(time.Second * 5)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
		// completion
		i := instr.Map["Completion"]
		want := Command{
			Instr: i,
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
		got, err = ZVT.ReadResponse(5 * time.Second)
		if assert.NoError(t, err) {
			if assert.EqualValues(t, want, *got) {
				ZVT.SendACK()
			}
		}

	}

}

// func TestDisplayText(t *testing.T) {
// 	want := aprc.Response{
// 		CCRC:   0x80,
// 		APRC:   0x00,
// 		Length: 0x00,
// 		Data:   []byte{},
// 	}
// 	got, err := ZVT.DisplayText([]string{
// 		"Da steh ich nun,",
// 		"ich armer Tor,",
// 		"Und bin so klug",
// 		"als wie zuvor."})
// 	if assert.NoError(t, err) {
// 		assert.EqualValues(t, want, *got)
// 	}
// }

func TestAbort(t *testing.T) {
	i := instr.Map["ACK"]
	want := Command{
		Instr: i,
		Data: apdu.DataUnit{
			Data:    []byte{},
			BMPOBJs: []bmp.OBJ{},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{},
			},
		},
	}
	err := ZVT.Abort()
	got, err := ZVT.ReadResponse(time.Second * 5)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
	}
}

func TestLogOff(t *testing.T) {
	i := instr.Map["ACK"]
	want := Command{
		Instr: i,
		Data: apdu.DataUnit{
			Data:    []byte{},
			BMPOBJs: []bmp.OBJ{},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{},
			},
		},
	}
	err := ZVT.LogOff()
	got, err := ZVT.ReadResponse(time.Second * 5)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
	}
}

// FIXME: getting error "0x83 function not possible" from PT
// func TestChangPassword(t *testing.T) {
// 	// start
// 	want := Response{
// 		CCRC:   0x80,
// 		APRC:   0x00,
// 		Length: 0x00,
// 		Data:   []byte{},
// 	}
// 	var old, new *[3]byte = &[3]byte{0x12, 0x34, 0x56},
// 		&[3]byte{0x12, 0x34, 0x56}
// 	got, err := ZVT.ChangePassword(old, new)
// 	if assert.NoError(t, err) {
// 		assert.EqualValues(t, want, *got)

// 		// completion
// 		want = Response{
// 			CCRC:   0x06,
// 			APRC:   0x0F,
// 			Length: 0x01,
// 			Data:   []byte{0},
// 		}
// 		got, err = ZVT.readResponse(5 * time.Second)
// 		if assert.NoError(t, err) {
// 			assert.EqualValues(t, want, *got)
// 		}
// 	}
// }
