package test

type Command struct {
	text string
}

type CompletionResponse interface {
	Process(*Command) error
}
type AuthorisationResponse struct {
	Status  byte
	Message string
}
type B struct {
	AuthorisationResponse
	Test int
}

func (r *AuthorisationResponse) Process(result *Command) error {
	return nil
}

func Completion(response CompletionResponse) error {
	result := &Command{}
	(response).Process(result)
	return nil
}

func TestAuthorisationCompletion() {
	resp := AuthorisationResponse{}
	Completion(&resp)
	X := B{
		AuthorisationResponse: resp,
		Test:                  0,
	}
	_ = X
}
