package event

import (
	"testing"
	"time"
)

func TestFormatAge(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{0, "0s"},
		{30 * time.Second, "30s"},
		{59 * time.Second, "59s"},
		{60 * time.Second, "1m"},
		{5 * time.Minute, "5m"},
		{59 * time.Minute, "59m"},
		{time.Hour, "1h"},
		{2 * time.Hour, "2h"},
		{23 * time.Hour, "23h"},
		{24 * time.Hour, "1d"},
		{48 * time.Hour, "2d"},
		{72 * time.Hour, "3d"},
	}

	for _, tt := range tests {
		got := FormatAge(tt.d)
		if got != tt.want {
			t.Errorf("FormatAge(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}
