package client

import (
	"bytes"
	"context"
	"duty-schedule/config"
	"duty-schedule/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CreateEvents adds basic duty information for each user in the group.
func (c *Client) CreateEvents(ctx context.Context) {
	config := config.Get()
	urlPath := fmt.Sprintf("%s/api/v0/events", c.oncallURL)
	log := c.logger.With().Str("method", "CreateEvents").Logger()

	for _, t := range config.Teams {
		for _, u := range t.Users {

			for _, d := range u.Duty {
				layout := "02/01/2006"
				date, err := time.Parse(layout, d.Date)
				if err != nil {
					log.Fatal().Err(err)
					return
				}

				// setting the time for the beginning of the day (00:00:00)
				location, _ := time.LoadLocation("Local")
				start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, location)
				// setting the time for the beginning of the next day (00:00:00)
				end := start.Add(24 * time.Hour)

				event := models.Event{Start: start.Unix(), End: end.Unix(), User: u.Name, Team: t.Name, Role: models.Role(d.Role)}

				jsonData, err := json.Marshal(event)
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
				case resp.StatusCode == http.StatusUnprocessableEntity:
					log.Error().Msg("Event creation failed: nonexistent role/event/team")
				case resp.StatusCode == http.StatusBadRequest:
					log.Error().Msg("Event validation checks failed")
				default:
					log.Error().Msg(fmt.Sprintf("CreateEvents failed with status code: %d", resp.StatusCode))
				}
			}
		}
	}
}
