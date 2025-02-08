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
	conf := config.NewLocalPrinterConfig("01P00C490700226", "192.168.0.33", "74462393", "")
	printer := printer.NewPrinter(&conf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	fmt.Print("Connecting to printer... ")
	if err := printer.Connect(ctx); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}

	fmt.Println("Success!")

	printer.Subscribe(ctx, func(r report.Report) {
		log.Printf("[UPDATE RECEIVED] %s: %v", r.Type, r.Payload)
	})

	// 1. Home the printer (this is correct)
	fmt.Print("Homing Printer... ")
	if err := printer.SendGCode("G28", ctx); err != nil {
		log.Fatalf("failed: %v", err)
	}
	fmt.Println("Success")

	// 2.1 Set absolute positioning and move
	fmt.Print("Moving Absolute Axes... ")
	if err := printer.SendGCode("G90\nM221 X1 Y1 Z1\n G0 X10 Y10 Z10 F3000", ctx); err != nil {
		log.Fatalf("failed: %v", err)
	}
	fmt.Println("Success")

	// 2.2 Set relative positioning and move
	fmt.Print("Moving Relative Axes... ")
	if err := printer.SendGCode("G91\nM221 X1 Y1 Z1\n G0 X50 Y50 Z50 F3000", ctx); err != nil {
		log.Fatalf("failed: %v", err)
	}
	fmt.Println("Success")
}
