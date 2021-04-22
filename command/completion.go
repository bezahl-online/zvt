package command

import (
	"fmt"
	"net"
	"time"

	"github.com/bezahl-online/zvt/messages"
)

type CompletionResponse interface {
	Process(*Command) error
}

type TransactionResponse struct {
	Status  byte
	Message string
}

// Definition of <status-byte>
// first byte of (06 0F) completion
const (
	Status_initialisation_necessary = 1 << iota
	Status_diagnosis_necessary
	Status_OPT_action_necessary
	Status_fillingstation_mode
	Status_vendingmachine_mode
)

// Completion implements slave mode of ECR
// ECR can instruct the PT to abort execution of the command
func (p *PT) Completion(response CompletionResponse) error {
	var err error
	var result *Command
	if result, err = p.ReadResponseWithTimeout(30 * time.Second); err != nil {
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			p.Status()
			if result, err = p.ReadResponse(); err != nil {
				return err
			}
			if result != nil && result.Data.BMPOBJs != nil &&
				len(result.Data.BMPOBJs) > 0 &&
				result.Data.BMPOBJs[0].ID == 0x27 {
				errCode := result.Data.BMPOBJs[0].Data[0]
				message, ok := messages.ErrorMessage[errCode]
				if ok {
					Logger.Info(fmt.Sprintf("PT: '%s'", message))
				} else {
					Logger.Error(fmt.Sprintf("06 0F: unmapped error message code %0X", result.Data.Data[0]))
				}
			}
			p.SendACK()
			return nil
		} else {
			return err
		}
	}
	if err = response.Process(result); err != nil {
		return err
	}
	if err = p.SendACK(); err != nil {
		return err
	}
	return nil
}
