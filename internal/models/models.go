package models

import (
	"time"
)

type Post struct {
	ID        string
	Title     string
	URL       string
	Category  string
	Section   string
	PhotoURL  string
	PubDate   time.Time
	TotalPost int
	PostNum   int
}

type StateLog struct {
	Feed string            `json:"Feed"`
	Sent map[string]string `json:"Sent"`
}
