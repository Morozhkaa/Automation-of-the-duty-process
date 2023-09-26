package models

type Team struct {
	Name               string `json:"name"`
	SchedulingTimezone string `json:"scheduling_timezone"`
	Email              string `json:"email"`
	SlackChannel       string `json:"slack_channel"`
}
