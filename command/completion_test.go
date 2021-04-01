package command

import "testing"

func TestPT_Completion(t *testing.T) {
	type args struct {
		response CompletionResponse
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
			if err := tt.p.Completion(tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("PT.Completion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
