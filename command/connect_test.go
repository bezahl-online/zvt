package command

import "testing"

func TestPT_Connect(t *testing.T) {
	tests := []struct {
		name string
		p    *PT
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Connect()
		})
	}
}

func Test_delay_getSeconds(t *testing.T) {
	tests := []struct {
		name string
		w    *delay
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.getSeconds(); got != tt.want {
				t.Errorf("delay.getSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_delay_wait(t *testing.T) {
	tests := []struct {
		name string
		w    *delay
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.wait()
		})
	}
}

func Test_delay_double(t *testing.T) {
	tests := []struct {
		name string
		w    *delay
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.double()
		})
	}
}
