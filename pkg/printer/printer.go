package printer

import (
	"context"
	"fmt"
	"sync"

	"github.com/RobertMNewton/bambu-golang-api/pkg/ftp"
	"github.com/RobertMNewton/bambu-golang-api/pkg/mqtt"
	"github.com/RobertMNewton/bambu-golang-api/pkg/mqtt/request"
	"github.com/RobertMNewton/bambu-golang-api/pkg/types/config"
)

type Printer struct {
	config     config.PrinterConfig
	mqttClient *mqtt.Client
	//httpClient *http.Client // currently not useful. Required to finish cloud implementation
	ftpClient *ftp.Client

	mu        sync.RWMutex
	connected bool

	sequence_id uint
}

func NewPrinter(config config.PrinterConfig) *Printer {
	return &Printer{
		config:     config,
		mqttClient: mqtt.NewClient(config),
		ftpClient:  ftp.NewClient(config),
	}
}

func (printer *Printer) Connect(ctx context.Context) error {
	if err := printer.mqttClient.Connect(ctx); err != nil {
		return fmt.Errorf("mqtt connection failed: %w", err)
	}

	printer.setConnected(true)

	return nil
}

func (printer *Printer) Subscribe(ctx context.Context, callback mqtt.ReportHandler) error {
	if err := printer.mqttClient.Subscribe(ctx, callback); err != nil {
		return fmt.Errorf("mqtt subscription failed: %w", err)
	}

	return nil
}

func (printer *Printer) Disconnect() {
	printer.mqttClient.Disconnect()
	printer.setConnected(false)
}

func (printer *Printer) IsConnected() bool {
	printer.mu.RLock()
	defer printer.mu.RUnlock()

	return printer.connected
}

func (printer *Printer) SendGCode(gcode string, ctx context.Context) error {
	request := request.CreateGCodeLineRequest(printer.GetNextSequenceId(), gcode)
	return printer.mqttClient.Publish(ctx, request)
}

func (printer *Printer) GetNextSequenceId() string {
	defer func() {
		printer.sequence_id += 1
	}()

	return fmt.Sprint(printer.sequence_id)
}

func (printer *Printer) setConnected(connected bool) {
	printer.mu.Lock()
	defer printer.mu.Unlock()
	printer.connected = connected
}
