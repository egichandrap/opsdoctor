package netscan

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

func normalizeHost(h string) string {
	h = strings.TrimSpace(h)
	h = strings.TrimPrefix(h, "https://")
	h = strings.TrimPrefix(h, "http://")
	if !strings.Contains(h, ":") {
		h = h + ":443"
	}
	return h
}

func CheckTLS(rawHost string) (string, error) {
	host := normalizeHost(rawHost)
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", host, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return "", err
	}
	defer conn.Close()

	state := conn.ConnectionState()
	cert := state.PeerCertificates[0]

	out := fmt.Sprintf("Host: %s\nTLS Version: %x\nCipherSuite: %#x\nIssuer: %s\nValid Until: %s\n",
		rawHost, state.Version, state.CipherSuite, cert.Issuer.CommonName, cert.NotAfter.Format(time.RFC1123))
	return out, nil
}

func CheckConnectivity(host string) string {
	// if host includes port use it; else try :443 then :80
	if !strings.Contains(host, ":") {
		// try 443
		if conn, err := net.DialTimeout("tcp", host+":443", 2*time.Second); err == nil {
			conn.Close()
			return fmt.Sprintf("Reachable: %s:443", host)
		}
		if conn, err := net.DialTimeout("tcp", host+":80", 2*time.Second); err == nil {
			conn.Close()
			return fmt.Sprintf("Reachable: %s:80", host)
		}
		return fmt.Sprintf("Unreachable: %s", host)
	}
	conn, err := net.DialTimeout("tcp", host, 2*time.Second)
	if err != nil {
		return fmt.Sprintf("Unreachable: %s (%v)", host, err)
	}
	conn.Close()
	return fmt.Sprintf("Reachable: %s", host)
}
