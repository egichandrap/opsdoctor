package tlscheck

import (
	"crypto/tls"
	"fmt"
)

func CheckTLSInfo(host string) {
	conn, err := tls.Dial("tcp", host, &tls.Config{})
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]

	fmt.Println("=== TLS Certificate ===")
	fmt.Println("Common Name:", cert.Subject.CommonName)
	fmt.Println("Valid From:", cert.NotBefore)
	fmt.Println("Valid Until:", cert.NotAfter)
	fmt.Println("=======================")
}
