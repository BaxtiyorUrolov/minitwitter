package models

import "time"

type User struct {
	ID             string    `json:"id"`
	UserID         uint      `json:"user_id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	UserName       string    `json:"user_name"`
	Password       string    `json:"password"`
	Bio            string    `json:"bio"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      int       `json:"deleted_at"`
}

type CreateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type UpdateUser struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	UserName       string `json:"user_name"`
	Email          string `json:"email"`
	Bio            string `json:"bio"`
	ProfilePicture string `json:"profile_picture"`
}

type UsersResponse struct {
	Users []User `json:"user_names"`
	Count int    `json:"count"`
}

type UpdateUserPassword struct {
	ID          string `json:"id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
