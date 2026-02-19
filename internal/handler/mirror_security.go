package handler

import (
	"context"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"time"

	"trykkeri-api/internal/errors"
)

const (
	mirrorDialTimeout  = 10 * time.Second
	mirrorMaxRedirects = 10
)

type mirrorLookupIPAddrFunc func(ctx context.Context, host string) ([]net.IPAddr, error)
type mirrorDialContextFunc func(ctx context.Context, network, address string) (net.Conn, error)

// mirrorBlockedPrefixes covers private, loopback, link-local, multicast, and reserved ranges.
var mirrorBlockedPrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),
	netip.MustParsePrefix("10.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("169.254.0.0/16"),
	netip.MustParsePrefix("172.16.0.0/12"),
	netip.MustParsePrefix("192.0.0.0/24"),
	netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("192.168.0.0/16"),
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"),
	netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("224.0.0.0/4"),
	netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("::/128"),
	netip.MustParsePrefix("::1/128"),
	netip.MustParsePrefix("2001:db8::/32"),
	netip.MustParsePrefix("fc00::/7"),
	netip.MustParsePrefix("fe80::/10"),
	netip.MustParsePrefix("ff00::/8"),
}

func validateMirrorRequestURL(targetURL *url.URL) error {
	if targetURL.Scheme != "http" && targetURL.Scheme != "https" {
		return errors.InvalidInput("url scheme must be http or https")
	}
	if targetURL.Host == "" || targetURL.Hostname() == "" {
		return errors.InvalidInput("url must have a host")
	}

	if ip := net.ParseIP(targetURL.Hostname()); ip != nil && isMirrorBlockedIP(ip) {
		return errors.InvalidInput("url must not target private/internal addresses")
	}

	return nil
}

func newMirrorHTTPClient() *http.Client {
	dialer := &net.Dialer{Timeout: mirrorDialTimeout}
	lookup := net.DefaultResolver.LookupIPAddr
	return &http.Client{
		Timeout:   mirrorFetchTimeout,
		Transport: newMirrorTransport(lookup, dialer.DialContext),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= mirrorMaxRedirects {
				return errors.InvalidInput("too many redirects")
			}
			return validateMirrorRequestURL(req.URL)
		},
	}
}

func newMirrorTransport(lookup mirrorLookupIPAddrFunc, dial mirrorDialContextFunc) *http.Transport {
	if lookup == nil {
		lookup = net.DefaultResolver.LookupIPAddr
	}
	if dial == nil {
		dialer := &net.Dialer{Timeout: mirrorDialTimeout}
		dial = dialer.DialContext
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	transport.DialContext = newMirrorDialContext(lookup, dial)
	return transport
}

func newMirrorDialContext(lookup mirrorLookupIPAddrFunc, dial mirrorDialContextFunc) mirrorDialContextFunc {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, errors.InvalidInput("invalid target address")
		}

		ipAddrs, err := resolveMirrorIPAddrs(ctx, host, lookup)
		if err != nil {
			return nil, err
		}

		for _, ipAddr := range ipAddrs {
			if isMirrorBlockedIP(ipAddr.IP) {
				return nil, errors.InvalidInput("url must not target private/internal addresses")
			}
		}

		var dialErr error
		for _, ipAddr := range ipAddrs {
			conn, err := dial(ctx, network, net.JoinHostPort(ipAddr.IP.String(), port))
			if err == nil {
				return conn, nil
			}
			dialErr = err
		}

		if dialErr != nil {
			return nil, dialErr
		}
		return nil, errors.InvalidInput("url host did not resolve to any addresses")
	}
}

func resolveMirrorIPAddrs(ctx context.Context, host string, lookup mirrorLookupIPAddrFunc) ([]net.IPAddr, error) {
	if ip := net.ParseIP(host); ip != nil {
		return []net.IPAddr{{IP: ip}}, nil
	}

	ipAddrs, err := lookup(ctx, host)
	if err != nil {
		return nil, errors.InvalidInput("failed to resolve host: %v", err)
	}
	if len(ipAddrs) == 0 {
		return nil, errors.InvalidInput("url host did not resolve to any addresses")
	}

	return ipAddrs, nil
}

func isMirrorBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}

	addr, ok := netip.AddrFromSlice(ip)
	if !ok {
		return true
	}
	addr = addr.Unmap()

	for _, blocked := range mirrorBlockedPrefixes {
		if blocked.Contains(addr) {
			return true
		}
	}
	return false
}
