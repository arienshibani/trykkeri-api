package handler

import (
	"context"
	stderrors "errors"
	"net"
	"net/url"
	"testing"

	"trykkeri-api/internal/errors"
)

func TestValidateMirrorRequestURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rawURL  string
		wantErr bool
	}{
		{
			name:    "rejects unsupported scheme",
			rawURL:  "ftp://example.com/file.txt",
			wantErr: true,
		},
		{
			name:    "rejects missing host",
			rawURL:  "http:///path",
			wantErr: true,
		},
		{
			name:    "rejects loopback literal",
			rawURL:  "http://127.0.0.1/admin",
			wantErr: true,
		},
		{
			name:    "rejects private IPv6 literal",
			rawURL:  "http://[::1]/admin",
			wantErr: true,
		},
		{
			name:    "accepts public host name",
			rawURL:  "https://example.com",
			wantErr: false,
		},
		{
			name:    "accepts public literal",
			rawURL:  "http://8.8.8.8",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			targetURL, err := url.Parse(tc.rawURL)
			if err != nil {
				t.Fatalf("url.Parse(%q) error = %v", tc.rawURL, err)
			}

			err = validateMirrorRequestURL(targetURL)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("validateMirrorRequestURL(%q) expected error, got nil", tc.rawURL)
				}
				if !stderrors.Is(err, errors.ErrInvalidInput) {
					t.Fatalf("validateMirrorRequestURL(%q) error = %v; want invalid input", tc.rawURL, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("validateMirrorRequestURL(%q) unexpected error: %v", tc.rawURL, err)
			}
		})
	}
}

func TestIsMirrorBlockedIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ip      string
		blocked bool
	}{
		{name: "loopback", ip: "127.0.0.1", blocked: true},
		{name: "private class A", ip: "10.1.2.3", blocked: true},
		{name: "private class B", ip: "172.16.1.1", blocked: true},
		{name: "private class C", ip: "192.168.1.10", blocked: true},
		{name: "metadata endpoint", ip: "169.254.169.254", blocked: true},
		{name: "ipv6 loopback", ip: "::1", blocked: true},
		{name: "public ipv4", ip: "8.8.8.8", blocked: false},
		{name: "public ipv6", ip: "2606:4700:4700::1111", blocked: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ip := net.ParseIP(tc.ip)
			if ip == nil {
				t.Fatalf("net.ParseIP(%q) returned nil", tc.ip)
			}

			got := isMirrorBlockedIP(ip)
			if got != tc.blocked {
				t.Fatalf("isMirrorBlockedIP(%s) = %v; want %v", tc.ip, got, tc.blocked)
			}
		})
	}
}

func TestNewMirrorDialContextRejectsPrivateResolvedIP(t *testing.T) {
	t.Parallel()

	lookup := func(context.Context, string) ([]net.IPAddr, error) {
		return []net.IPAddr{{IP: net.ParseIP("127.0.0.1")}}, nil
	}

	dialCalled := false
	dial := func(context.Context, string, string) (net.Conn, error) {
		dialCalled = true
		return nil, nil
	}

	dialContext := newMirrorDialContext(lookup, dial)
	_, err := dialContext(context.Background(), "tcp", "example.com:80")
	if err == nil {
		t.Fatal("expected error for private resolved address, got nil")
	}
	if !stderrors.Is(err, errors.ErrInvalidInput) {
		t.Fatalf("dial error = %v; want invalid input", err)
	}
	if dialCalled {
		t.Fatal("dial function should not be called when resolved IP is blocked")
	}
}

func TestNewMirrorDialContextRejectsMixedPublicAndPrivateIPs(t *testing.T) {
	t.Parallel()

	lookup := func(context.Context, string) ([]net.IPAddr, error) {
		return []net.IPAddr{
			{IP: net.ParseIP("8.8.8.8")},
			{IP: net.ParseIP("10.0.0.1")},
		}, nil
	}

	dialCalled := false
	dial := func(context.Context, string, string) (net.Conn, error) {
		dialCalled = true
		return nil, nil
	}

	dialContext := newMirrorDialContext(lookup, dial)
	_, err := dialContext(context.Background(), "tcp", "example.com:80")
	if err == nil {
		t.Fatal("expected error for mixed public/private resolution, got nil")
	}
	if !stderrors.Is(err, errors.ErrInvalidInput) {
		t.Fatalf("dial error = %v; want invalid input", err)
	}
	if dialCalled {
		t.Fatal("dial function should not be called when any resolved IP is blocked")
	}
}

func TestNewMirrorDialContextDialsResolvedPublicIP(t *testing.T) {
	t.Parallel()

	lookup := func(context.Context, string) ([]net.IPAddr, error) {
		return []net.IPAddr{{IP: net.ParseIP("8.8.8.8")}}, nil
	}

	var gotDialAddress string
	dial := func(_ context.Context, _ string, address string) (net.Conn, error) {
		gotDialAddress = address
		client, server := net.Pipe()
		_ = server.Close()
		return client, nil
	}

	dialContext := newMirrorDialContext(lookup, dial)
	conn, err := dialContext(context.Background(), "tcp", "example.com:443")
	if err != nil {
		t.Fatalf("dialContext returned unexpected error: %v", err)
	}
	_ = conn.Close()

	if gotDialAddress != "8.8.8.8:443" {
		t.Fatalf("dialed address = %q; want %q", gotDialAddress, "8.8.8.8:443")
	}
}
