package service

import (
	"errors"
	"math/rand"
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

// CreatePR - основная логика задания
func (s *Service) CreatePR(title string, authorID int) (*domain.PullRequest, error) {
	// 1. Находим автора, чтобы узнать его команду
	author, err := s.store.GetUser(authorID)
	if err != nil {
		return nil, errors.New("автор не найден")
	}

	// 2. Получаем всех коллег
	teamMembers, err := s.store.GetTeamMembers(author.TeamID)
	if err != nil {
		return nil, err
	}

	// 3. Фильтруем: исключаем автора и неактивных
	var candidates []int64
	for _, user := range teamMembers {
		if user.IsActive && user.ID != author.ID {
			candidates = append(candidates, int64(user.ID))
		}
	}

	// 4. Выбираем случайных (максимум 2)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	limit := 2
	if len(candidates) < limit {
		limit = len(candidates)
	}
	reviewers := candidates[:limit]

	// 5. Создаем PR
	pr := &domain.PullRequest{
		Title:       title,
		AuthorID:    authorID,
		Status:      "OPEN",
		ReviewerIDs: pq.Int64Array(reviewers),
		CreatedAt:   time.Now(),
	}

	return pr, s.store.CreatePR(pr)
}
