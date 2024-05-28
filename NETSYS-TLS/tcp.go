package main

import (
	"crypto/tls"
	"fmt"
)

func main() {
	conn, err := tls.Dial("tcp", "www.example.com:443", nil)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	state := conn.ConnectionState()
	fmt.Println("TLS Version:", tlsVersionToString(state.Version))
	fmt.Println("Cipher Suite:", tls.CipherSuiteName(state.CipherSuite))

	for _, cert := range state.PeerCertificates {
		fmt.Println("Issuer Organization:", cert.Issuer.Organization)
	}
}

func tlsVersionToString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown"
	}
}