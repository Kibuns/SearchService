package Models

import (
	"time"
)

type Twoot struct {
	UserID    string    `json:"userid"`
	UserName	string	`json:"username"`
	Content  string    `json:"content"`
	Hashtags []string	`json:"hashtags"`
	Created time.Time `json:"created"`
}
