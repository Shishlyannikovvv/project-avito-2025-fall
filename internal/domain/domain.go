package domain

import "time"

// Team представляет команду
type Team struct {
	ID   int
	Name string
}

// User представляет пользователя
type User struct {
	ID       int
	Username string
	TeamID   int
	IsActive bool
}

// PR представляет Pull Request
type PR struct {
	ID          int
	Title       string
	AuthorID    int
	Status      string
	Reviewer1ID *int
	Reviewer2ID *int
	CreatedAt   time.Time
	MergedAt    *time.Time
}

// CreateTeamRequest запрос на создание команды
type CreateTeamRequest struct {
	Name string `json:"name"`
}

// CreateUserRequest запрос на создание пользователя
type CreateUserRequest struct {
	Username string `json:"username"`
	TeamID   int    `json:"team_id"`
}

// CreatePRRequest запрос на создание PR
type CreatePRRequest struct {
	Title    string `json:"title"`
	AuthorID int    `json:"author_id"`
}

// ReassignRequest запрос на переназначение
type ReassignRequest struct {
	ReviewerID int `json:"reviewer_id"`
}
