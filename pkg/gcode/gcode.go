package gcode

import (
	"fmt"
	"strings"
)

type Builder struct {
	commands []string
}

func New() *Builder {
	return &Builder{
		commands: make([]string, 0),
	}
}

func (b *Builder) String() string {
	return strings.Join(b.commands, "\n")
}

func (b *Builder) Commands() []string {
	return b.commands
}

func (b *Builder) AddCommand(cmd string) *Builder {
	b.commands = append(b.commands, cmd)
	return b
}

func (b *Builder) HomeAll() *Builder {
	return b.AddCommand("G28")
}

func (b *Builder) HomeXY() *Builder {
	return b.AddCommand("G28 X")
}

func (b *Builder) HomeZ(precision bool, temp float64) *Builder {
	if temp > 0 {
		return b.AddCommand(fmt.Sprintf("G28 Z P0 T%.1f", temp))
	}
	return b.AddCommand("G28 Z P0")
}

func (b *Builder) BedMeshCalibration() *Builder {
	return b.AddCommand("G29")
}

func (b *Builder) SetZOffset(offset float64) *Builder {
	return b.AddCommand(fmt.Sprintf("G29.1 Z%.3f", offset))
}

func (b *Builder) MoveZ(z float64, speed float64) *Builder {
	return b.AddCommand(fmt.Sprintf("G380 S2 Z%.3f F%.1f", z, speed))
}

// Coordinate System Commands

// AbsolutePositioning adds G90 command for absolute coordinates
func (b *Builder) AbsolutePositioning() *Builder {
	return b.AddCommand("G90")
}

func (b *Builder) RelativePositioning() *Builder {
	return b.AddCommand("G91")
}

func (b *Builder) ResetExtruder() *Builder {
	return b.AddCommand("G92 E0")
}

func (b *Builder) SetExtruderRelative() *Builder {
	return b.AddCommand("M83")
}

func (b *Builder) SetHotendTemp(temp float64) *Builder {
	if temp < 0 || temp > 300 {
		panic("hotend temperature must be between 0 and 300")
	}
	return b.AddCommand(fmt.Sprintf("M104 S%.1f", temp))
}

func (b *Builder) WaitForHotend(temp float64) *Builder {
	if temp < 0 || temp > 300 {
		panic("hotend temperature must be between 0 and 300")
	}
	return b.AddCommand(fmt.Sprintf("M109 S%.1f", temp))
}

func (b *Builder) SetBedTemp(temp float64) *Builder {
	if temp < 0 || temp > 110 {
		panic("bed temperature must be between 0 and 110")
	}
	return b.AddCommand(fmt.Sprintf("M140 S%.1f", temp))
}

func (b *Builder) WaitForBed(temp float64) *Builder {
	if temp < 0 || temp > 110 {
		panic("bed temperature must be between 0 and 110")
	}
	return b.AddCommand(fmt.Sprintf("M190 S%.1f", temp))
}

func (b *Builder) SetPartFanSpeed(speed int) *Builder {
	if speed < 0 || speed > 255 {
		panic("fan speed must be between 0 and 255")
	}
	return b.AddCommand(fmt.Sprintf("M106 P1 S%d", speed))
}

func (b *Builder) SetAuxFanSpeed(speed int) *Builder {
	if speed < 0 || speed > 255 {
		panic("fan speed must be between 0 and 255")
	}
	return b.AddCommand(fmt.Sprintf("M106 P2 S%d", speed))
}

func (b *Builder) SetChamberFanSpeed(speed int) *Builder {
	if speed < 0 || speed > 255 {
		panic("fan speed must be between 0 and 255")
	}
	return b.AddCommand(fmt.Sprintf("M106 P3 S%d", speed))
}

func (b *Builder) LinearMove(x, y, z, feedRate float64) *Builder {
	cmd := strings.Builder{}
	cmd.WriteString("G0")
	if x != 0 {
		cmd.WriteString(fmt.Sprintf(" X%.3f", x))
	}
	if y != 0 {
		cmd.WriteString(fmt.Sprintf(" Y%.3f", y))
	}
	if z != 0 {
		cmd.WriteString(fmt.Sprintf(" Z%.3f", z))
	}
	if feedRate != 0 {
		cmd.WriteString(fmt.Sprintf(" F%.1f", feedRate))
	}
	return b.AddCommand(cmd.String())
}

func (b *Builder) LinearPrint(x, y, z, feedRate, extrusion float64) *Builder {
	cmd := strings.Builder{}
	cmd.WriteString("G1")
	if x != 0 {
		cmd.WriteString(fmt.Sprintf(" X%.3f", x))
	}
	if y != 0 {
		cmd.WriteString(fmt.Sprintf(" Y%.3f", y))
	}
	if z != 0 {
		cmd.WriteString(fmt.Sprintf(" Z%.3f", z))
	}
	if feedRate != 0 {
		cmd.WriteString(fmt.Sprintf(" F%.1f", feedRate))
	}
	if extrusion != 0 {
		cmd.WriteString(fmt.Sprintf(" E%.3f", extrusion))
	}
	return b.AddCommand(cmd.String())
}

func (b *Builder) SetHorizontalLaser(on bool) *Builder {
	return b.AddCommand(fmt.Sprintf("M960 S1 %d", boolToInt(on)))
}

func (b *Builder) SetVerticalLaser(on bool) *Builder {
	return b.AddCommand(fmt.Sprintf("M960 S2 %d", boolToInt(on)))
}

func (b *Builder) SetNozzleLED(on bool) *Builder {
	return b.AddCommand(fmt.Sprintf("M960 S4 %d", boolToInt(on)))
}

func (b *Builder) SetLogoLED(on bool) *Builder {
	return b.AddCommand(fmt.Sprintf("M960 S5 %d", boolToInt(on)))
}

func (b *Builder) SetAllLEDs(on bool) *Builder {
	return b.AddCommand(fmt.Sprintf("M960 S0 %d", boolToInt(on)))
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
