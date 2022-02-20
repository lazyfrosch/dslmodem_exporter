package zyxel

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func testClient(t *testing.T) *Client {
	var (
		url      = os.Getenv("DSLEXPORTER_MODEM_URL")
		username = os.Getenv("DSLEXPORTER_MODEM_USERNAME")
		password = os.Getenv("DSLEXPORTER_MODEM_PASSWORD")
	)

	if url == "" || username == "" || password == "" {
		t.Skip("Needs environment variables: DSLEXPORTER_MODEM_URL, DSLEXPORTER_MODEM_USERNAME, DSLEXPORTER_MODEM_PASSWORD")
		return nil
	}

	c, err := NewClient(url, username, password)
	if err != nil {
		t.Fatal(err)
	}

	return c
}

func TestClient_Login(t *testing.T) {
	c := testClient(t)

	err := c.Login()
	assert.NoError(t, err)
}

func TestClient_GetXDSLStatistics(t *testing.T) {
	c := testClient(t)

	err := c.Login()
	assert.NoError(t, err)

	data, err := c.GetXDSLStatistics()
	assert.NoError(t, err)
	assert.NotNil(t, data)
}
