package model

type UpdateUserRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
}
