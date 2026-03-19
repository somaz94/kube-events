package event

import "fmt"

// FormatAge formats a duration into a human-readable short form (e.g., 5s, 3m, 2h, 1d).
func FormatAge(d interface{ Seconds() float64 }) string {
	sec := d.Seconds()
	switch {
	case sec < 60:
		return fmt.Sprintf("%ds", int(sec))
	case sec < 3600:
		return fmt.Sprintf("%dm", int(sec/60))
	case sec < 86400:
		return fmt.Sprintf("%dh", int(sec/3600))
	default:
		return fmt.Sprintf("%dd", int(sec/86400))
	}
}
