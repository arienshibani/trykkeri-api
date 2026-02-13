package ssrf

import (
	"testing"
)

func TestBlockPrivateOrInternal(t *testing.T) {
	tests := []struct {
		host string
		want bool // true = must be blocked
	}{
		{"127.0.0.1", true},
		{"127.0.0.1:8080", true},
		{"localhost", true},
		{"localhost:80", true},
		{"192.168.1.1", true},
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"169.254.169.254", true},
		{"169.254.1.1", true},
		{"0.0.0.0", true},
		{"[::1]", true},
		{"[::1]:443", true},
	}
	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			err := BlockPrivateOrInternal(tt.host)
			blocked := err == ErrHostBlocked
			if blocked != tt.want {
				t.Errorf("BlockPrivateOrInternal(%q) err=%v, want blocked=%v", tt.host, err, tt.want)
			}
		})
	}
}

func TestBlockPrivateOrInternal_Public(t *testing.T) {
	// Only run if we expect public DNS to work
	hosts := []string{"example.com", "8.8.8.8"}
	for _, host := range hosts {
		err := BlockPrivateOrInternal(host)
		if err == ErrHostBlocked {
			t.Errorf("BlockPrivateOrInternal(%q) should allow public host, got blocked", host)
		}
		// Other errors (e.g. network) are ok in tests
	}
}
