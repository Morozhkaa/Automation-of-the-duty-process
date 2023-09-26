package client

import (
	"bytes"
	"context"
	"duty-schedule/config"
	"duty-schedule/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/yaml.v2"
)

// AddUsers adds existing users to the previously created "Group-1" roster of the desired team.
func (c *Client) AddUsers(ctx context.Context) {
	config := config.Get()
	log := c.logger.With().Str("method", "AddUsers").Logger()

	for _, t := range config.Teams {
		urlPath := fmt.Sprintf("%s/api/v0/teams/%s/rosters/%s/users", c.oncallURL, t.Name, rosterName)

		for _, u := range t.Users {
			user := models.Username{
				Name: u.Name,
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

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}
			var userInfo models.User
			yaml.Unmarshal(body, &userInfo)
			if err != nil {
				return
			}

			switch {
			case resp.StatusCode == http.StatusCreated:
				log.Info().Msg(fmt.Sprintf("User '%s' added successfully", user.Name))
			case resp.StatusCode == http.StatusUnprocessableEntity:
				log.Error().Msg(fmt.Sprintf("Invalid team/user or user is already in roster: Team -'%s', User - '%s'", t.Name, u.Name))
			case resp.StatusCode == http.StatusBadRequest:
				log.Error().Msg("Missing 'name' parameter")
			default:
				log.Error().Msg(fmt.Sprintf("AddUsers failed with status code: %d", resp.StatusCode))
			}
		}
	}
}
