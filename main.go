package main

import (
	"context"
	"duty-schedule/client"
)

func main() {
	c := client.New()
	c.CreateDutySchedule(context.Background())
}
