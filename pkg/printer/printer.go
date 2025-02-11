package printer

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

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

	sequence_id atomic.Uint32
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

func (printer *Printer) SendRequest(request request.Request, ctx context.Context) error {
	request.SetSequenceID(printer.getNextSequenceId())
	return printer.mqttClient.Publish(ctx, request)
}

func (printer *Printer) SendGCode(gcode string, ctx context.Context) error {
	return printer.SendRequest(request.CreateGCodeLineRequest("", gcode), ctx)
}

func (printer *Printer) StartPrint(filename string, ctx context.Context) error {
	return printer.SendRequest(request.CreateGCodeFileRequest("", filename), ctx)
}

func (printer *Printer) PausePrint(ctx context.Context) error {
	return printer.SendRequest(request.CreatePausePrintRequest(""), ctx)
}

func (printer *Printer) ResumePrint(ctx context.Context) error {
	return printer.SendRequest(request.CreateResumePrintRequest(""), ctx)
}

func (printer *Printer) StopPrint(ctx context.Context) error {
	return printer.SendRequest(request.CreateStopPrintRequest(""), ctx)
}

func (printer *Printer) UnloadFilament(ctx context.Context) error {
	return printer.SendRequest(request.CreateUnloadFilamentRequest(""), ctx)
}

func (printer *Printer) LoadFilament(ctx context.Context) error {
	return printer.SendRequest(request.CreateLoadFilamentRequest(""), ctx)
}

func (printer *Printer) getNextSequenceId() string {
	id := printer.sequence_id.Load()
	printer.sequence_id.Add(1)
	return fmt.Sprint(id)
}

func (printer *Printer) setConnected(connected bool) {
	printer.mu.Lock()
	defer printer.mu.Unlock()
	printer.connected = connected
}
