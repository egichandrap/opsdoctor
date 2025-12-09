package tlscheck

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

func normalizeHost(host string) string {
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "http://")
	if !strings.Contains(host, ":") {
		host += ":443"
	}
	return host
}

func CheckTLS(rawHost string) (string, error) {
	host := normalizeHost(rawHost)

	dialer := &net.Dialer{Timeout: 5 * time.Second}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         strings.Split(host, ":")[0],
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", host, tlsConfig)
	if err != nil {
		return "", fmt.Errorf("TLS handshake failed: %w", err)
	}
	defer conn.Close()

	cs := conn.ConnectionState()
	cert := cs.PeerCertificates[0]

	result := fmt.Sprintf(
		"Host: %s\nTLS Version: %x\nCipher: %x\nIssuer: %s\nValid Until: %s\n",
		rawHost,
		cs.Version,
		cs.CipherSuite,
		cert.Issuer.CommonName,
		cert.NotAfter.Format(time.RFC1123),
	)

	return result, nil
}
