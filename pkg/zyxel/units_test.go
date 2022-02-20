package zyxel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDuration_UnmarshalText(t *testing.T) {
	var d Uptime

	tests := map[string]float64{ // minutes
		"41 days: 20 hours: 46 minutes": 60286,
		"0 day: 0 hour: 1 minute":       1,
		"0 day: 0 hour: 2 minutes":      2,
		"0 day: 0 hour: 0 minute":       0,
	}

	for test, expected := range tests {
		err := d.UnmarshalText([]byte(test))
		assert.NoError(t, err)
		assert.Equal(t, expected, d.Duration.Minutes())
	}

	for _, test := range []string{
		"",
		"1 units",
		"x minutes",
	} {
		err := d.UnmarshalText([]byte(test))
		assert.Error(t, err)
	}

}
