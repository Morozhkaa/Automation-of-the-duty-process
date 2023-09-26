package models

type User struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Contacts struct {
		Email string `json:"email"`
		Call  string `json:"call"`
	} `json:"contacts"`
}

type Username struct {
	Name string `json:"name"`
}
