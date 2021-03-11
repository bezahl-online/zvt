package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/albenik/bcd"
	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/instr"
)

// AuthConfig is the auth data struct
type AuthConfig struct {
	Amount      int64
	PaymentType byte
	Currency    int
	// TLV         *tlv.Container
}

type CardData struct {
	Name  string
	Type  int
	PAN   string
	Tech  int
	SeqNr int
}

type AuthResultData struct {
	Amount     int64
	ReceiptNr  int
	TurnoverNr int
	TraceNr    int
	Date       string
	Time       string
	TID        string
	VU         string
	AID        string
	Card       CardData
}

const (
	Result_Success = "success"
	Result_Abort   = "abort"
	Result_Timeout = "timeout"
)

type AuthResult struct {
	Result string
	Data   *AuthResultData
}

// Authorisation implents 06 01
// initiates a payment process
func (p *PT) Authorisation(config *AuthConfig) (AuthResult, error) {
	ctrlField := instr.Map["Authorisation"]
	var result AuthResult
	err := p.send(Command{ctrlField, config.marshal()})
	if err != nil {
		return result, err
	}
	got, err := PaymentTerminal.ReadResponse()
	if err != nil {
		return result, err
	}
	if got.IsAck() {
		for x := 10; x > 0; x-- {
			got, err = PaymentTerminal.ReadResponseWithTimeout(20 * time.Second)
			if err != nil {
				return result, err
			}
			PaymentTerminal.SendACK()
			switch got.CtrlField.Class {
			case 0x06:
				switch got.CtrlField.Instr {
				case 0x1E:
					switch got.Data.Data[0] {
					case 0x6C:
						fmt.Println("Transaction aborted")
						result.Result = Result_Abort
					}
					return result, nil

				case 0x0F:
					fmt.Println("Transaction successfull")
					result.Result = Result_Success
					result.Data = &AuthResultData{
						Amount: 0,
						Card: CardData{
							Tech: 0,
						},
					}
					result.Data.FromOBJs(got.Data.BMPOBJs)
					return result, nil
				}
			}
		}
	}
	return result, fmt.Errorf("timeout or connection lost")
}

func formatPAN(rawPAN []byte) string {
	raw := fmt.Sprintf("%X", rawPAN)
	raw = strings.ReplaceAll(raw, "E", "X")
	pan := raw[0:4]
	l := int(len(raw) / 4)
	for i := 1; i < l; i++ {
		p := i * 4
		pan += " " + raw[p:p+4]
	}
	pan = strings.TrimRight(pan, "F")
	return pan
}

func (a *AuthConfig) marshal() apdu.DataUnit {
	return apdu.DataUnit{
		BMPOBJs: []bmp.OBJ{
			{ID: 0x04, Data: bcd.FromUint(uint64(a.Amount), 6)},
		},
	}
}

func (r *AuthResultData) FromOBJs(objs []bmp.OBJ) {
	for _, obj := range objs {
		switch obj.ID {
		case 0x04:
			amount := bcd.ToUint64(obj.Data)
			r.Amount = int64(amount)
		case 0x0B:
			r.TraceNr = int(bcd.ToUint32(obj.Data))
		case 0x0C:
			r.Time = fmt.Sprintf("%06X", obj.Data)
		case 0x0D:
			r.Date = fmt.Sprintf("%04X", obj.Data)
		case 0x17:
			r.Card.SeqNr = int(bcd.ToUint16(obj.Data))
		case 0x22:
			pan := formatPAN(obj.Data)
			r.Card.PAN = pan
		case 0x27:
			// FIXME: map error codes
		case 0x29:
			r.TID = fmt.Sprintf("%X", obj.Data)
		case 0x2A:
			r.VU = strings.TrimSpace(string(obj.Data))
		case 0x3B:
			r.AID = strings.Trim(string(obj.Data), string(byte(0x00)))
		case 0x87:
			r.ReceiptNr = int(bcd.ToUint16(obj.Data))
		case 0x88:
			r.TurnoverNr = int(bcd.ToUint64(obj.Data))
		case 0x8B:
			r.Card.Name = strings.Trim(string(obj.Data), string(byte(0x00)))
		case 0x8A:
			r.Card.Type = int(obj.Data[0])
		}
	}

}
