package zyxel

import (
	"fmt"
	"github.com/lazyfrosch/dslmodem_exporter/pkg/units"
	"strconv"
	"strings"
	"time"
)

// Uptime is a Zyxel specific type derived from time.Duration.
type Uptime struct {
	time.Duration
}

// UnmarshalText parses the uptime format from Zyxel devices.
//
// Example:
//  41 days: 20 hours: 46 minutes
func (d *Uptime) UnmarshalText(data []byte) error {
	// Reset value
	d.Duration = 0

	for _, spec := range strings.Split(string(data), ": ") {
		numberUnit := strings.SplitN(spec, " ", 2) // number unit
		if len(numberUnit) != 2 {
			return fmt.Errorf("duration has unknown format: %s", spec)
		}

		var multiplier time.Duration

		switch u := numberUnit[1]; u {
		case "day", "days":
			multiplier = time.Hour * 24
		case "hour", "hours":
			multiplier = time.Hour
		case "minute", "minutes":
			multiplier = time.Minute
		default:
			return fmt.Errorf("invalid unit '%s' in duration format: %s", u, spec)
		}

		// Parse the numeric value
		value, err := strconv.ParseUint(numberUnit[0], 10, 64)
		if err != nil {
			return fmt.Errorf("could not parse number in spec: %s - %w", spec, err)
		}

		// increment duration for part
		d.Duration = d.Duration + time.Duration(value)*multiplier
	}

	return nil
}

// Delay is a Zyxel specific type derived from time.Duration.
type Delay struct {
	time.Duration
}

// UnmarshalText parses the uptime format from Zyxel devices.
//
// Example:
//  41 days: 20 hours: 46 minutes
func (d *Delay) UnmarshalText(data []byte) error {
	f, err := units.UnmarshalFloatWithUnits(data, "ms")
	if err != nil {
		return err
	}

	d.Duration = time.Duration(f) * time.Millisecond

	return nil
}
