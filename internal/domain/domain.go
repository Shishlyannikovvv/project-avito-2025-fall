package domain

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"unique;not null" json:"name"`
	TeamID   int    `json:"team_id"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

type Team struct {
	ID    int    `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"unique;not null" json:"name"`
	Users []User `gorm:"foreignKey:TeamID" json:"users,omitempty"`
}

type PullRequest struct {
	ID          int           `gorm:"primaryKey" json:"id"`
	Title       string        `json:"title"`
	AuthorID    int           `json:"author_id"`
	Status      string        `json:"status"` // OPEN, MERGED
	ReviewerIDs pq.Int64Array `gorm:"type:bigint[]" json:"reviewer_ids"`
	CreatedAt   time.Time     `json:"created_at"`
}
