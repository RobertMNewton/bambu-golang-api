package utils

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

func GetPrinterCert(host string, port string) (*x509.Certificate, error) {
	// Set up TLS config to accept any certificate
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Connect with timeout
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", net.JoinHostPort(host, port), conf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Get the peer certificates
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("no certificates found")
	}

	return certs[0], nil
}

func SaveCertToPEM(cert *x509.Certificate, fp string) error {
	// Create the directory if it doesn't exist
	if dir := filepath.Dir(fp); dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Create the file
	file, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write the certificate in PEM format
	err = pem.Encode(file, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	if err != nil {
		return fmt.Errorf("failed to encode certificate: %w", err)
	}

	return nil
}
