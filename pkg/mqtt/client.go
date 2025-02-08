package mqtt

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RobertMNewton/bambu-golang-api/pkg/mqtt/report"
	"github.com/RobertMNewton/bambu-golang-api/pkg/mqtt/request"
	"github.com/RobertMNewton/bambu-golang-api/pkg/types/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	keepAlive      = 30 * time.Second
	connectTimeout = 10 * time.Second
)

type Client struct {
	config    config.PrinterConfig
	client    mqtt.Client
	mu        sync.RWMutex
	connected bool
}

type ReportHandler func(report.Report)

func NewClient(config config.PrinterConfig) *Client {
	return &Client{
		config: config,
	}
}

func (client *Client) Connect(ctx context.Context) error {
	options, err := client.createClientOptions()
	if err != nil {
		return nil
	}

	client.client = mqtt.NewClient(options)

	connectChan := make(chan error, 1)
	go func() {
		token := client.client.Connect()
		if token.Wait() && token.Error() != nil {
			connectChan <- token.Error()
			return
		}
		connectChan <- nil
	}()

	select {
	case err := <-connectChan:
		if err != nil {
			return fmt.Errorf("mqtt connection failed: %w", err)
		}
	case <-ctx.Done():
		return fmt.Errorf("connection timeout: %w", ctx.Err())
	}

	client.setConnected(true)

	return nil
}

func (client *Client) Disconnect() {
	if client.client != nil && client.client.IsConnected() {
		client.client.Disconnect(250)
	}

	client.setConnected(false)
}

func (client *Client) IsConnected() bool {
	client.mu.RLock()
	defer client.mu.RUnlock()

	return client.connected
}

func (client *Client) Subscribe(ctx context.Context, callback ReportHandler) error {
	if !client.IsConnected() {
		return fmt.Errorf("mqtt client not connected")
	}

	topic := fmt.Sprintf("device/%s/report", client.config.GetDeviceID())

	handler := func(client mqtt.Client, msg mqtt.Message) {
		report, err := report.FromMessage(msg)
		if err != nil {
			return
		}

		callback(report)
	}

	token := client.client.Subscribe(topic, 1, handler)
	if err := token.Error(); err != nil {
		return fmt.Errorf("subscription failed: %w", err)
	}

	return nil
}

func (client *Client) Publish(ctx context.Context, request request.Request) error {
	if !client.IsConnected() {
		return fmt.Errorf("mqtt client is not connected")
	}

	topic := fmt.Sprintf("device/%s/request", client.config.GetDeviceID())
	msg, err := request.ToMessage()
	if err != nil {
		return err
	}

	// Create token first
	token := client.client.Publish(topic, 1, false, msg)

	done := make(chan error, 1)
	go func() {
		token.Wait()
		done <- token.Error()
		close(done)
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("publish failed: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("publish timeout: %w", ctx.Err())
	case <-time.After(5 * time.Second):
		return nil
	}
}

func (client *Client) createClientOptions() (*mqtt.ClientOptions, error) {
	options := mqtt.NewClientOptions()
	options.AddBroker(client.config.GetBrokerUrl())
	options.SetUsername(client.config.GetUsername())
	options.SetPassword(client.config.GetPassword())
	options.SetKeepAlive(keepAlive)
	options.SetConnectTimeout(connectTimeout)

	tls_config, err := client.config.CreateTLSConfig()
	if err != nil {
		return nil, err
	}

	options.SetTLSConfig(tls_config)

	return options, nil
}

func (client *Client) setConnected(connected bool) {
	client.mu.Lock()
	defer client.mu.Unlock()

	client.connected = connected
}
