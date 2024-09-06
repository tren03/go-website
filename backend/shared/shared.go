package shared

import (
	"time"
)

type Post struct {
	ID    int       `json:"ID"`
	Title string    `json:"Title"`
	Url   string    `json:"Url"`
	Body  string    `json:"Body"`
	Date  time.Time `json:"Date"`
}
