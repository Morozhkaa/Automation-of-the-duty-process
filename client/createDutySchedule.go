package client

import "context"

func (c *Client) CreateDutySchedule(ctx context.Context) {
	c.CreateTeams(ctx)
	c.CreateRosters(ctx)
	// c.DeleteUsers(ctx)
	c.CreateUsers(ctx)
	c.UpdateUsers(ctx)
	c.AddUsers(ctx)
	c.GetUsers(ctx)
	c.CreateEvents(ctx)
}
