package numbers

import "testing"

func TestParseBigFloat(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		returnErr bool
	}{
		{"valid float string", "4.23", false},
		{"invalid float string", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ParseBigFloat(tt.value); err != nil && !tt.returnErr {
				t.Error("Parsing test failed")
			}
		})
	}
}
