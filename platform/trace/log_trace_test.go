package trace

import "testing"

func TestTrace(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "trace",
			want: "log_trace_test.go -> TestTrace:17",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Trace(); got != tt.want {
				t.Errorf("Trace() = %v, want %v", got, tt.want)
			}
		})
	}
}
