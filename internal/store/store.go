package store

import (
	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/domain"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

// --- Team ---
func (s *Store) CreateTeam(team *domain.Team) error {
	return s.db.Create(team).Error
}

// --- User ---
func (s *Store) CreateUser(user *domain.User) error {
	return s.db.Create(user).Error
}

func (s *Store) GetUser(id int) (*domain.User, error) {
	var user domain.User
	err := s.db.First(&user, id).Error
	return &user, err
}

func (s *Store) GetTeamMembers(teamID int) ([]domain.User, error) {
	var users []domain.User
	err := s.db.Where("team_id = ?", teamID).Find(&users).Error
	return users, err
}

// --- Pull Request ---
func (s *Store) CreatePR(pr *domain.PullRequest) error {
	return s.db.Create(pr).Error
}

func (s *Store) GetPR(id int) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	err := s.db.First(&pr, id).Error
	return &pr, err
}

func (s *Store) UpdatePR(pr *domain.PullRequest) error {
	return s.db.Save(pr).Error
}
