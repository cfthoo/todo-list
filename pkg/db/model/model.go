package model

import "time"

type Task struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	CreatedBy  string    `json:"created_by"`
	Complete   bool      `json:"complete"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}
