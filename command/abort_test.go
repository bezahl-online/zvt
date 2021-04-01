package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbort(t *testing.T) {
	skipShort(t)
	err := PaymentTerminal.Abort()
	assert.NoError(t, err)
}

func TestPT_Abort(t *testing.T) {
	tests := []struct {
		name    string
		p       *PT
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Abort(); (err != nil) != tt.wantErr {
				t.Errorf("PT.Abort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
