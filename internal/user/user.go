package user

import "context"

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Honor    int    `json:"honor"`
}

type Repository interface {
	GetUserByID(ctx context.Context, id string) (User, error)
}
