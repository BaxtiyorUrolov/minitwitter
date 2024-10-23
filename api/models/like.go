package models

import "time"

type Like struct {
	UserID    string    `json:"user_id"`
	TweetID   string    `json:"tweet_id"`
	CreatedAt time.Time `json:"created_at"`
}
