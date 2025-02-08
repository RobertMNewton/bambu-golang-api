package config

import (
	"crypto/tls"
	"fmt"
)

const (
	cloudBrokerUrl      = "tls://us.mqtt.bambulab.com:8883"
	cloudUsernameFormat = "u_%s"
)

type CloudPrinterConfig struct {
	deviceID string

	userID      string
	accessToken string

	// needed for FTP
	ipAddress  string
	accessCode string
}

func NewCloudPrinterConfig(device_id, user_id, access_token, ipAddress, accessCode string) *CloudPrinterConfig {
	return &CloudPrinterConfig{
		deviceID: device_id,

		userID:      user_id,
		accessToken: access_token,

		ipAddress:  ipAddress,
		accessCode: accessCode,
	}
}

func (config *CloudPrinterConfig) GetDeviceID() string {
	return config.deviceID
}

func (config *CloudPrinterConfig) GetBrokerUrl() string {
	return cloudBrokerUrl
}

func (config *CloudPrinterConfig) GetDeviceAccessCode() string {
	return config.accessCode
}

func (config *CloudPrinterConfig) GetUsername() string {
	return fmt.Sprintf(cloudUsernameFormat, config.userID)
}

func (config *CloudPrinterConfig) GetPassword() string {
	return config.accessToken
}

func (config *CloudPrinterConfig) CreateTLSConfig() (*tls.Config, error) {
	return &tls.Config{}, nil
}
