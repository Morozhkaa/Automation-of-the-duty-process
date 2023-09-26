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

var rosterName = "Group-1"

// CreateRosters creates a "Group-1" roster for each team.
func (c *Client) CreateRosters(ctx context.Context) {
	config := config.Get()
	log := c.logger.With().Str("method", "CreateRosters").Logger()

	for _, t := range config.Teams {
		r := models.Roster{
			Name: rosterName,
		}
		jsonData, err := json.Marshal(r)
		if err != nil {
			log.Fatal().Err(err)
		}
		urlPath := fmt.Sprintf("%s/api/v0/teams/%s/rosters", c.oncallURL, t.Name)

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
			log.Info().Msg(fmt.Sprintf("Roster in group '%s' added successfully", t.Name))
		case resp.StatusCode == http.StatusUnprocessableEntity:
			log.Error().Msg(fmt.Sprintf("Invalid character in roster name/Duplicate roster name: '%s'", rosterName))
		default:
			log.Error().Msg(fmt.Sprintf("CreateRosters failed with status code: %d", resp.StatusCode))
		}
	}
}
