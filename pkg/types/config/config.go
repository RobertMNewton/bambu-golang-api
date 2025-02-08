package config

import (
	"crypto/tls"
)

type PrinterConfig interface {
	GetDeviceID() string
	GetBrokerUrl() string
	GetDeviceIPAddress() string
	GetDeviceAccessCode() string
	GetUsername() string
	GetPassword() string
	CreateTLSConfig() (*tls.Config, error)
}
