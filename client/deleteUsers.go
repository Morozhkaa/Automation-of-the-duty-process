package client

import (
	"context"
	"duty-schedule/config"
	"fmt"
	"net/http"
)

// DeleteUsers removes users from the system (written for convenience).
func (c *Client) DeleteUsers(ctx context.Context) {
	config := config.Get()
	log := c.logger.With().Str("method", "DeleteUsers").Logger()

	for _, t := range config.Teams {
		for _, u := range t.Users {
			urlPath := fmt.Sprintf("%s/api/v0/users/%s", c.oncallURL, u.Name)

			ctx, cancel := context.WithTimeout(ctx, c.defaultTimeout)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, "DELETE", urlPath, nil)
			if err != nil {
				log.Fatal().Err(err)
			}
			req.Header.Set("X-CSRF-Token", c.csrfToken)

			resp, err := c.client.Do(req)
			if err != nil {
				log.Fatal().Err(err)
			}
			defer resp.Body.Close()

			switch {
			case resp.StatusCode == http.StatusOK:
				log.Info().Msg(fmt.Sprintf("User '%s' deleted successfully", u.Name))
			default:
				log.Error().Msg(fmt.Sprintf("DeleteUsers failed with status code: %d", resp.StatusCode))
			}
		}
	}
}
