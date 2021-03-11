package command

import (
	"fmt"
	"time"

	"github.com/albenik/bcd"
	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/instr"
)

// AuthConfig is the auth data struct
type AuthConfig struct {
	Amount      int
	PaymentType byte
	Currency    int
	// TLV         *tlv.Container
}

type CardData struct {
	Name string
	Type string
	PAN  string
	Tech int
}

type AuthResult struct {
	Success   bool
	Card      CardData
	ReceiptNr string
	TID       string
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
					if got.Data.Data[0] == 0x6C {
						fmt.Println("Transaction aborted")
						return result, nil
					}
				case 0x0F:
					fmt.Println("Transaction successfull")
					return result, nil
				}
			}
		}
	}
	return result, fmt.Errorf("timeout or connection lost")
}

func (a *AuthConfig) marshal() apdu.DataUnit {
	return apdu.DataUnit{
		BMPOBJs: []bmp.OBJ{
			{ID: 0x04, Data: bcd.FromUint(uint64(a.Amount), 6)},
		},
	}
}
