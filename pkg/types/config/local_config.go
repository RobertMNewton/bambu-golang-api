package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/RobertMNewton/bambu-golang-api/pkg/utils"
)

const (
	localUsername     = "bblp"
	localBrokerFormat = "tls://%s:8883"
)

type LocalPrinterConfig struct {
	deviceID   string
	ipAddress  string
	accessCode string
	caCertPath string
}

func NewLocalPrinterConfig(device_id, ip_address, access_code, ca_cert_path string) LocalPrinterConfig {
	return LocalPrinterConfig{
		deviceID:   device_id,
		ipAddress:  ip_address,
		accessCode: access_code,
		caCertPath: ca_cert_path,
	}
}

func (config *LocalPrinterConfig) GetDeviceID() string {
	return config.deviceID
}

func (config *LocalPrinterConfig) GetBrokerUrl() string {
	return fmt.Sprintf(localBrokerFormat, config.ipAddress)
}

func (config *LocalPrinterConfig) GetDeviceIPAddress() string {
	return config.ipAddress
}

func (config *LocalPrinterConfig) GetDeviceAccessCode() string {
	return config.accessCode
}

func (config *LocalPrinterConfig) GetUsername() string {
	return localUsername
}

func (config *LocalPrinterConfig) GetPassword() string {
	return config.accessCode
}

func (config *LocalPrinterConfig) CreateTLSConfig() (*tls.Config, error) {
	caCertPool := x509.NewCertPool()

	if config.caCertPath != "" {
		caCert, err := os.ReadFile(config.caCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %v", err)
		}

		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA certificate")
		}
	} else {
		caCert, err := utils.GetPrinterCert(config.GetDeviceIPAddress(), "8883")
		if err != nil {
			return nil, fmt.Errorf("failed to get CA certificate: %v", err)
		}

		caCertPool.AddCert(caCert)
	}

	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		ServerName:         config.GetDeviceID(),
		InsecureSkipVerify: true, // Disable default verification
		VerifyConnection: func(cs tls.ConnectionState) error {
			cert := cs.PeerCertificates[0]

			if len(cert.DNSNames) == 0 {
				if cert.Subject.CommonName == config.GetDeviceID() {
					return nil // Allow CN match
				}
				return fmt.Errorf("certificate does not contain SANs, and CN does not match host")
			}

			// Check if the cert chains properly
			opts := x509.VerifyOptions{
				Roots: caCertPool,
			}
			_, err := cert.Verify(opts)
			return err
		},
	}

	return tlsConfig, nil
}
