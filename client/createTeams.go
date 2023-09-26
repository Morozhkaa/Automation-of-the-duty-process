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

// CreateTeams creates the teams described in the config.
func (c *Client) CreateTeams(ctx context.Context) {
	config := config.Get()
	urlPath := fmt.Sprintf("%s/api/v0/teams", c.oncallURL)
	log := c.logger.With().Str("method", "CreateTeams").Logger()

	for _, t := range config.Teams {
		team := models.Team{
			Name:               t.Name,
			SchedulingTimezone: t.SchedulingTimezone,
			Email:              t.Email,
			SlackChannel:       t.SlackChannel,
		}

		jsonData, err := json.Marshal(team)
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
			log.Info().Msg(fmt.Sprintf("Team '%s' added successfully", team.Name))
		case resp.StatusCode == http.StatusUnprocessableEntity:
			log.Error().Msg(fmt.Sprintf("Duplicate team name '%s'", team.Name))
		case resp.StatusCode == http.StatusBadRequest:
			log.Error().Msg("Error in creating team. Possible errors: API key auth not allowed, invalid attributes, missing required attributes")
		default:
			log.Error().Msg(fmt.Sprintf("CreateTeams failed with status code: %d", resp.StatusCode))
		}
	}
}
