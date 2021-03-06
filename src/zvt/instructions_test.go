package zvt

import (
	"testing"
	"time"

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
	serviceByte := ServiceMenuNOTAssignedToFunctionKey +
		DisplayTextsForCommandsAuthorisation
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
	want := Response{
		CCRC:   0x80,
		APRC:   0x00,
		Length: 0x00,
		Data:   []byte{},
	}
	got, err := ZVT.Register(&Config{
		pwd:          fixedPassword,
		config:       byte(configByte),
		currency:     EUR,
		service:      byte(serviceByte),
		tlvContainer: tlvContainer,
	})
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
		// completion
		want = Response{
			CCRC:   0x06,
			APRC:   0x0F,
			Length: 0x0a,
			Data:   []byte{0x19, 0x0, 0x29, 0x29, 0x0, 0x10, 0x6, 0x49, 0x9, 0x78},
		}
		got, err = ZVT.readResponse(5 * time.Second)
		if assert.NoError(t, err) {
			assert.EqualValues(t, want, *got)
		}

	}

}

func TestDisplayText(t *testing.T) {
	want := Response{
		CCRC:   0x80,
		APRC:   0x00,
		Length: 0x00,
		Data:   []byte{},
	}
	got, err := ZVT.DisplayText([]string{
		"Da steh ich nun,",
		"ich armer Tor,",
		"Und bin so klug",
		"als wie zuvor."})
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
	}
}

func TestAbort(t *testing.T) {
	want := Response{
		CCRC:   0x80,
		APRC:   0x00,
		Length: 0x00,
		Data:   []byte{},
	}
	got, err := ZVT.Abort()
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
	}
}

func TestLogOff(t *testing.T) {
	want := Response{
		CCRC:   0x80,
		APRC:   0x00,
		Length: 0x00,
		Data:   []byte{},
	}
	got, err := ZVT.LogOff()
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
