package units

import "testing"

func TestBitsHuman(t *testing.T) {
	tests := []struct {
		name string
		bits int64
		want string
	}{
		{name: "zero", bits: 0, want: "0.00 bps"},
		{name: "positive", bits: 8_000_000, want: "8.00 Mbps"},
		{name: "negative", bits: -8_000_000, want: "-8.00 Mbps"},
		{name: "minimum integer", bits: -1 << 63, want: "-9223372.04 Tbps"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BitsHuman(tt.bits); got != tt.want {
				t.Errorf("BitsHuman(%d) = %q, want %q", tt.bits, got, tt.want)
			}
		})
	}
}
