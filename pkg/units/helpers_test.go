package units

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalFloatWithUnits(t *testing.T) {
	for test, expected := range map[string]float64{
		"10 Mbps":     10,
		"10":          10,
		"54.123 Mbps": 54.123,
		"54.123Mbps":  54.123,
		"54.123":      54.123,
		"-10.2 dBm":   -10.2,
		"- 10.2 dBm":  -10.2,
	} {
		f, err := UnmarshalFloatWithUnits([]byte(test), "Mbps", "dBm")
		assert.NoError(t, err)
		assert.Equal(t, expected, f)
	}

	for _, test := range []string{
		"",
		"x",
		"10 ",
	} {
		_, err := UnmarshalFloatWithUnits([]byte(test))
		assert.Error(t, err)
	}
}
