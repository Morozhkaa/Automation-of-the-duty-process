package models

type Role string

var (
	RolePrimary     = "primary"
	RoleSecondary   = "secondary"
	RoleShadow      = "shadow"
	RoleManager     = "manager"
	RoleVacation    = "vacation"
	RoleUnavailable = "unavailable"
)

type Event struct {
	Start int64  `json:"start"`
	End   int64  `json:"end"`
	User  string `json:"user"`
	Team  string `json:"team"`
	Role  Role   `json:"role"`
}
