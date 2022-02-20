package zyxel

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestVDSLStatus_UnmarshalText(t *testing.T) {
	var status VDSLStatus
	data, err := ioutil.ReadFile("testdata/example-vdsl.txt")
	assert.NoError(t, err)

	err = status.UnmarshalText(data)
	assert.NoError(t, err)
	assert.Equal(t, "Showtime", status.Status)
	assert.Equal(t, "VDSL2 Annex B", status.Mode)
	assert.Equal(t, "Profile 17a", status.Profile)
	assert.Equal(t, float64(60286), status.LinkUptime.Minutes())
	assert.Equal(t, 36.998, float64(status.LineRateUp))
	assert.Equal(t, 116.799, float64(status.LineRateDown))
}
