# Bambu Golang API

A Golang package for seamlessly interfacing with Bambu Lab Printers. With this API, you can:

- Send GCODE commands directly to your printer
- Automatically initiate and terminate prints through your code
- Monitor and track your printer's status in real-time

## Features

- Send GCODE commands to your Bambu Lab printer
- Monitor printer status and receive live updates
- Control printer movement with both absolute and relative positioning
- Optional: Secure connection with PEM certificate support

## Example Usage

Here's a simple example to get you started with connecting and controlling your printer:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RobertMNewton/bambu-golang-api/pkg/mqtt/report"
	"github.com/RobertMNewton/bambu-golang-api/pkg/printer"
	"github.com/RobertMNewton/bambu-golang-api/pkg/types/config"
)

func main() {
	// Configure the printer with your device's details
	conf := config.NewLocalPrinterConfig("{YOUR DEVICE ID}", "{YOUR DEVICE IP ADDRESS}", "{YOUR DEVICE ACCESS CODE}", "{PATH TO PEM CERT (OPTIONAL)}")
	printer := printer.NewPrinter(&conf)

	// Set up a context with a 5-minute timeout for operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Attempt to connect to the printer
	fmt.Print("Connecting to printer... ")
	if err := printer.Connect(ctx); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}

	fmt.Println("Connection successful!")

	// Subscribe to printer status updates
	printer.Subscribe(ctx, func(r report.Report) {
		log.Printf("[UPDATE RECEIVED] %s: %v", r.Type, r.Payload)
	})

	// Home the printer
	fmt.Print("Homing printer... ")
	if err := printer.SendGCode("G28", ctx); err != nil {
		log.Fatalf("Failed: %v", err)
	}
	fmt.Println("Success")

	// Absolute positioning movement
	fmt.Print("Moving to absolute position... ")
	if err := printer.SendGCode("G90\nM221 X1 Y1 Z1\n G0 X10 Y10 Z10 F3000", ctx); err != nil {
		log.Fatalf("Failed: %v", err)
	}
	fmt.Println("Success")

	// Relative positioning movement
	fmt.Print("Moving to relative position... ")
	if err := printer.SendGCode("G91\nM221 X1 Y1 Z1\n G0 X50 Y50 Z50 F3000", ctx); err != nil {
		log.Fatalf("Failed: %v", err)
	}
	fmt.Println("Success")
}
```

## Future Features
- [] Custom Error Types
- [] Custom Report and Status Types
- [] "Safe" GCODE Executor for enhanced safety during printing
- [] Live Video Stream from the printer
- [] Testing support for CloudPrinterConfig (Currently only LocalPrinterConfig has been tested)
## Acknowledgements

A special thanks to OpenBambuAPI for their incredible work, which served as the foundation for this project.
