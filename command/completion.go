package command

import "time"

type CompletionResponse interface {
	Process(*Command) error
}

type TransactionResponse struct {
	Status  byte
	Message string
}

// Completion implements slave mode of ECR
// ECR can instruct the PT to abort execution of the command
func (p *PT) Completion(response CompletionResponse) error {
	var err error
	var result *Command
	if result, err = p.ReadResponseWithTimeout(5 * time.Minute); err != nil {
		return err
	}
	if err = p.SendACK(); err != nil {
		return err
	}
	if err = response.Process(result); err != nil {
		return err
	}
	return nil
}
