package client

import (
	"context"
	"duty-schedule/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetUsers gets all users on the system (written for convenience).
func (c *Client) GetUsers(ctx context.Context) {
	urlPath := fmt.Sprintf("%s/api/v0/users", c.oncallURL)
	log := c.logger.With().Str("method", "GetUsers").Logger()

	ctx, cancel := context.WithTimeout(ctx, c.defaultTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", urlPath, nil)
	if err != nil {
		log.Fatal().Err(err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal().Err(err)
		}
		var responseBody []models.User
		err = json.Unmarshal(body, &responseBody)
		if err != nil {
			log.Fatal().Err(err)
		}
		for _, elem := range responseBody {
			log.Info().Msg(fmt.Sprintf("Name: '%s', Fullname: '%s', Email: '%s', Phone: '%s'", elem.Name, elem.FullName, elem.Contacts.Email, elem.Contacts.Call))
		}
	} else {
		log.Error().Msg(fmt.Sprintf("GetUsers:	failed with status code: %d", resp.StatusCode))
	}
}
