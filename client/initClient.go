package client

import (
	"context"
	"duty-schedule/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Client struct {
	oncallURL      string
	csrfToken      string
	logger         zerolog.Logger
	client         *http.Client
	defaultTimeout time.Duration
}

// New returns a new instance of an authorized Client.
func New() *Client {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil
	}
	client := &Client{
		oncallURL: "http://127.0.0.1:8080",
		logger:    zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger(),
		client: &http.Client{
			Jar: cookieJar,
		},
		defaultTimeout: 5 * time.Second,
	}
	client.Login(context.Background())
	return client

}

// Login performs authorization and initializes global variables csrfToken, cookie.
func (c *Client) Login(ctx context.Context) {
	urlPath := fmt.Sprintf("%s/login", c.oncallURL)
	log := c.logger.With().Str("method", "Login").Logger()

	formData := url.Values{}
	formData.Set("username", "root")
	formData.Set("password", "123")

	// Creating a POST request and set request headers
	ctx, cancel := context.WithTimeout(ctx, c.defaultTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", urlPath, strings.NewReader(formData.Encode()))
	if err != nil {
		log.Fatal().Err(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send a request
	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer resp.Body.Close()

	// Process the response
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal().Err(err)
		}
		var responseBody models.ResponseBody
		err = json.Unmarshal(body, &responseBody)
		if err != nil {
			log.Fatal().Err(err)
		}
		c.csrfToken = responseBody.CsrfToken
		log.Info().Msg("Login successful")
	} else {
		log.Error().Msg(fmt.Sprintf("Login failed: %d", resp.StatusCode))
	}
}
