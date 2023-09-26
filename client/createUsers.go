package client

import (
	"bytes"
	"context"
	"duty-schedule/config"
	"duty-schedule/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateUsers creates the users described in the config.
func (c *Client) CreateUsers(ctx context.Context) {
	config := config.Get()
	urlPath := fmt.Sprintf("%s/api/v0/users", c.oncallURL)
	log := c.logger.With().Str("method", "CreateUsers").Logger()

	for _, t := range config.Teams {
		for _, u := range t.Users {
			user := models.User{
				Name:     u.Name,
				FullName: u.FullName,
				Contacts: struct {
					Email string `json:"email"`
					Call  string `json:"call"`
				}{
					Email: u.Email,
					Call:  u.PhoneNumber,
				},
			}
			jsonData, err := json.Marshal(user)
			if err != nil {
				log.Fatal().Err(err)
			}

			ctx, cancel := context.WithTimeout(ctx, c.defaultTimeout)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, "POST", urlPath, bytes.NewBuffer(jsonData))
			if err != nil {
				log.Fatal().Err(err)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-CSRF-Token", c.csrfToken)

			resp, err := c.client.Do(req)
			if err != nil {
				log.Fatal().Err(err)
			}
			defer resp.Body.Close()

			switch {
			case resp.StatusCode == http.StatusCreated:
				log.Info().Msg(fmt.Sprintf("User '%s' added successfully", user.Name))
			case resp.StatusCode == http.StatusUnprocessableEntity:
				log.Error().Msg(fmt.Sprintf("Duplicate user name '%s'", user.Name))
			case resp.StatusCode == http.StatusBadRequest:
				log.Error().Msg("Error in creating user. Possible errors: API key auth not allowed, invalid attributes, missing required attributes")
			default:
				log.Error().Msg(fmt.Sprintf("CreateUsers failed with status code: %d", resp.StatusCode))
			}
		}
	}
}
