package ssrf

import (
	"fmt"
	"net"
	"strings"
)

// ErrHostBlocked is returned when the host resolves to a private, link-local, or
// metadata address that must not be fetched (SSRF protection).
var ErrHostBlocked = fmt.Errorf("url host is not allowed (private or internal)")

// BlockPrivateOrInternal returns ErrHostBlocked if host resolves to any IP in
// private ranges, link-local, or cloud metadata. Host may include a port.
func BlockPrivateOrInternal(host string) error {
	hostname, _, err := net.SplitHostPort(host)
	if err != nil {
		if strings.Contains(err.Error(), "missing port") {
			hostname = host
		} else {
			return err
		}
	}
	hostname = strings.Trim(hostname, "[]")
	if hostname == "" {
		return fmt.Errorf("empty host")
	}
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return err
	}
	if len(ips) == 0 {
		return fmt.Errorf("no addresses for host")
	}
	for _, ip := range ips {
		if isBlockedIP(ip) {
			return ErrHostBlocked
		}
	}
	return nil
}

func isBlockedIP(ip net.IP) bool {
	ip4 := ip.To4()
	if ip4 != nil {
		return isBlockedIPv4(ip4)
	}
	return isBlockedIPv6(ip)
}

func isBlockedIPv4(ip net.IP) bool {
	if len(ip) != 4 {
		return false
	}
	// 127.0.0.0/8
	if ip[0] == 127 {
		return true
	}
	// 10.0.0.0/8
	if ip[0] == 10 {
		return true
	}
	// 172.16.0.0/12
	if ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31 {
		return true
	}
	// 192.168.0.0/16
	if ip[0] == 192 && ip[1] == 168 {
		return true
	}
	// 169.254.0.0/16 (link-local, includes 169.254.169.254 metadata)
	if ip[0] == 169 && ip[1] == 254 {
		return true
	}
	// 0.0.0.0/8 ("this" network)
	if ip[0] == 0 {
		return true
	}
	return false
}

func isBlockedIPv6(ip net.IP) bool {
	if len(ip) != 16 {
		return false
	}
	// ::1
	if ip.Equal(net.IPv6loopback) {
		return true
	}
	// fc00::/7 (unique local)
	if ip[0] == 0xfc || ip[0] == 0xfd {
		return true
	}
	// fe80::/10 (link-local)
	if ip[0] == 0xfe && (ip[1]&0xc0) == 0x80 {
		return true
	}
	return false
}
