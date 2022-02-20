package units

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reNumberUnit            = regexp.MustCompile(`^(?i)((?:-\s?)?\d+(?:\.\d+)?)(?:\s?(\w+))?$`)
	ErrInvalidFloatWithUnit = errors.New("invalid float value with optional unit")
)

func UnmarshalFloatWithUnits(data []byte, units ...string) (float64, error) {
	m := reNumberUnit.FindSubmatch(data)
	if m == nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidFloatWithUnit, string(data))
	}

	// Validate unit
	if len(m[2]) > 0 {
		var (
			allowed bool
			unit    = string(m[2])
		)

		for _, u := range units {
			if strings.EqualFold(u, unit) {
				allowed = true
				break
			}
		}

		if !allowed {
			return 0, ErrInvalidFloatWithUnit
		}
	}

	// Fix negative numbers for strconv.ParseFloat
	if len(m[1]) > 3 {
		if bytes.Equal(m[1][0:2], []byte("- ")) {
			m[1] = append(m[1][0:1], m[1][2:]...)
		}
	}

	// Parse number
	f, err := strconv.ParseFloat(string(m[1]), 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse float: %w", err)
	}

	return f, nil
}
