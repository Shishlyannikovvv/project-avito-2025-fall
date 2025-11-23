package service

import (
	"time"

	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/domain"
	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/store"
	"github.com/lib/pq"
)

type Service struct {
	store *store.Store
}

func New(store *store.Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateTeam(name string) (*domain.Team, error) {
	team := &domain.Team{Name: name}
	return team, s.store.CreateTeam(team)
}

func (s *Service) CreateUser(name string, teamID int) (*domain.User, error) {
	user := &domain.User{Name: name, TeamID: teamID, IsActive: true}
	return user, s.store.CreateUser(user)
}

func (s *Service) CreatePR(title string, authorID int) (*domain.PullRequest, error) {
	// 1. Находим команду автора (в реальном коде нужен метод GetUser)
	// Для MVP упростим: считаем, что teamID передается или мы его знаем.
	// Тут нужно дописать получение User, чтобы узнать его TeamID.

	// Заглушка:
	reviewers := pq.Int64Array{}

	pr := &domain.PullRequest{
		Title:       title,
		AuthorID:    authorID,
		Status:      "OPEN",
		ReviewerIDs: reviewers,
		CreatedAt:   time.Now(),
	}
	return pr, s.store.CreatePR(pr)
}
