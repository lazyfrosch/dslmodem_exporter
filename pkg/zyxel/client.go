package zyxel

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var (
	ErrHTTPUnauthorized = errors.New("HTTP responded unauthorized")
)

type Client struct {
	BaseURL  *url.URL
	Username string
	Password string

	http.Client
}

func NewClient(baseURL, username, password string) (*Client, error) {
	jar, _ := cookiejar.New(nil)

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		BaseURL:  parsedURL,
		Username: username,
		Password: password,
		Client: http.Client{
			Jar: jar,
			Transport: &Transport{
				next: http.DefaultTransport,
			},
		},
	}, nil
}

// Login into the DSL modems admin interface.
//
//  /login/login-page.cgi
func (c *Client) Login() error {
	data := url.Values{}
	data.Set("AuthName", c.Username)
	data.Set("AuthPassword", c.Password)

	response, err := c.Post("/login/login-page.cgi", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("could not authenticate: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("login responded with non-ok status: %s", response.Status)
	} else if response.Header.Get("Set-Cookie") == "" {
		return fmt.Errorf("login responded without setting a cookie")
	}

	return nil
}

// GetXDSLStatistics returns the raw
//
//  /pages/systemMonitoring/xdslStatistics/GetxdslStatistics.html
func (c *Client) GetXDSLStatistics() (*VDSLStatus, error) {
	response, err := c.Post("/pages/systemMonitoring/xdslStatistics/GetxdslStatistics.html", "", nil)
	if err != nil {
		return nil, fmt.Errorf("could not get DSL statistics: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("server responded with non-ok status: %s", response.Status)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	status := &VDSLStatus{}
	err = status.UnmarshalText(data)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (c *Client) requestWithBaseURL(r *http.Request) *http.Request {
	if r.URL.Host == "" && c.BaseURL != nil {
		// Inject Base URL
		r.URL.Scheme = c.BaseURL.Scheme
		r.URL.Host = c.BaseURL.Host
	}

	return r
}

func (c *Client) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return c.Do(c.requestWithBaseURL(req))
}

type Transport struct {
	next http.RoundTripper
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := t.next.RoundTrip(r)
	if err != nil {
		return resp, err
	}

	// Load body to be able to parse it
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, fmt.Errorf("could not read body: %w", err)
	}

	// Check if login is required
	if bytes.Contains(body, []byte("top.location='/login/login.html';")) {
		resp.Status = "401 Unauthorized (redirect)"
		resp.StatusCode = http.StatusUnauthorized
		err = ErrHTTPUnauthorized

		// TODO: try to (re)login here
	}

	// Rebuild body to pass
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	return resp, err
}
