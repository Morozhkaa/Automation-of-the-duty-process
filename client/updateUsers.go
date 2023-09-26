package client

import (
	"bytes"
	"context"
	"duty-schedule/config"
	"duty-schedule/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// UpdateUsers saves detailed information about users (it's strange that the CreateUsers method doesn't do this).
func (c *Client) UpdateUsers(ctx context.Context) {
	config := config.Get()
	log := c.logger.With().Str("method", "UpdateUsers").Logger()

	for _, t := range config.Teams {
		for _, u := range t.Users {
			encodedUserName := url.PathEscape(u.Name)
			urlPath := fmt.Sprintf("%s/api/v0/users/%s", c.oncallURL, encodedUserName)
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
			req, err := http.NewRequestWithContext(ctx, "PUT", urlPath, bytes.NewBuffer(jsonData))
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
			case resp.StatusCode == http.StatusNoContent:
				log.Info().Msg(fmt.Sprintf("User information '%s' was successfully updated", user.Name))
			case resp.StatusCode == http.StatusNotFound:
				log.Error().Msg(fmt.Sprintf("User '%s' not found", u.Name))
			default:
				log.Error().Msg(fmt.Sprintf("UpdateUsers failed with status code: %d", resp.StatusCode))
			}
		}
	}
}
