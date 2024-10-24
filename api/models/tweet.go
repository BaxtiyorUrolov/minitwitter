package models

import "time"

type Tweet struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Content    string    `json:"content"`
	Media      string    `json:"media"`
	ViewsCount int       `json:"views_count"`
	LikesCount int       `json:"likes_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateTweet struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
	Media   string `json:"media"`
}

type UpdateTweet struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
	Media   string `json:"media"`
}

type TweetsResponse struct {
	Count  int     `json:"count"`
	Tweets []Tweet `json:"tweets"`
}
