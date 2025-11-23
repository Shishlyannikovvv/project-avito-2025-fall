package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/domain"
	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/store"
)

// Service содержит бизнес-логику
type Service struct {
	store *store.Store
	cache sync.Map // Кэш активных юзеров по командам (изюминка для скорости)
}

// NewService создаёт сервис
func NewService(store *store.Store) *Service {
	return &Service{store: store}
}

// CreateTeam создаёт команду
func (s *Service) CreateTeam(ctx context.Context, name string) (*domain.Team, error) {
	return s.store.CreateTeam(ctx, name)
}

// CreateUser создаёт юзера
func (s *Service) CreateUser(ctx context.Context, username string, teamID int) (*domain.User, error) {
	user, err := s.store.CreateUser(ctx, username, teamID)
	if err == nil {
		s.invalidateCache(teamID)
	}
	return user, err
}

// CreatePR создаёт PR и назначает ревьюверов
func (s *Service) CreatePR(ctx context.Context, title string, authorID int) (*domain.PR, error) {
	author, err := s.store.GetUser(ctx, authorID)
	if err != nil {
		return nil, err
	}

	pr, err := s.store.CreatePR(ctx, title, authorID)
	if err != nil {
		return nil, err
	}

	reviewers, err := s.getActiveReviewers(ctx, author.TeamID, authorID)
	if err != nil {
		return nil, err
	}

	s.shuffleReviewers(reviewers) // Изюминка: Fisher-Yates с crypto/rand

	if len(reviewers) > 2 {
		reviewers = reviewers[:2]
	}

	if len(reviewers) >= 1 {
		pr.Reviewer1ID = &reviewers[0].ID
	}
	if len(reviewers) >= 2 {
		pr.Reviewer2ID = &reviewers[1].ID
	}

	return s.store.UpdatePRReviewers(ctx, pr.ID, pr.Reviewer1ID, pr.Reviewer2ID)
}

// ReassignReviewer переназначает ревьювера
func (s *Service) ReassignReviewer(ctx context.Context, prID, reviewerID int) (*domain.PR, error) {
	pr, err := s.store.GetPR(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == "MERGED" {
		return nil, errors.New("cannot reassign merged PR")
	}

	var oldReviewerID *int
	if pr.Reviewer1ID != nil && *pr.Reviewer1ID == reviewerID {
		oldReviewerID = pr.Reviewer1ID
		pr.Reviewer1ID = nil
	} else if pr.Reviewer2ID != nil && *pr.Reviewer2ID == reviewerID {
		oldReviewerID = pr.Reviewer2ID
		pr.Reviewer2ID = nil
	} else {
		return nil, errors.New("reviewer not found in PR")
	}

	oldReviewer, err := s.store.GetUser(ctx, *oldReviewerID)
	if err != nil {
		return nil, err
	}

	candidates, err := s.getActiveReviewers(ctx, oldReviewer.TeamID, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	s.shuffleReviewers(candidates)

	var newReviewerID *int
	for _, cand := range candidates {
		if (pr.Reviewer1ID == nil || *pr.Reviewer1ID != cand.ID) && (pr.Reviewer2ID == nil || *pr.Reviewer2ID != cand.ID) {
			newReviewerID = &cand.ID
			break
		}
	}

	if newReviewerID == nil {
		return nil, errors.New("no available reviewers")
	}

	if pr.Reviewer1ID == nil {
		pr.Reviewer1ID = newReviewerID
	} else {
		pr.Reviewer2ID = newReviewerID
	}

	return s.store.UpdatePRReviewers(ctx, prID, pr.Reviewer1ID, pr.Reviewer2ID)
}

// MergePR мержит PR идемпотентно
func (s *Service) MergePR(ctx context.Context, prID int) (*domain.PR, error) {
	pr, err := s.store.GetPR(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == "MERGED" {
		return pr, nil // Идемпотентность
	}

	return s.store.MergePR(ctx, prID)
}

// GetUserPRs возвращает PR юзера
func (s *Service) GetUserPRs(ctx context.Context, userID int) ([]*domain.PR, error) {
	return s.store.GetUserPRs(ctx, userID)
}

// GetReviewerStats статистика (бонус)
func (s *Service) GetReviewerStats(ctx context.Context) (map[int]int, error) {
	return s.store.GetReviewerStats(ctx)
}

// DeactivateTeam деактивирует команду и переназначает (бонус)
func (s *Service) DeactivateTeam(ctx context.Context, teamID int) error {
	users, err := s.store.GetTeamUsers(ctx, teamID)
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.IsActive {
			if err := s.store.DeactivateUser(ctx, u.ID); err != nil {
				return err
			}
			s.invalidateCache(teamID)
		}
	}

	prs, err := s.store.GetOpenPRsByTeam(ctx, teamID)
	if err != nil {
		return err
	}

	for _, pr := range prs {
		if pr.Reviewer1ID != nil {
			rev1, _ := s.store.GetUser(ctx, *pr.Reviewer1ID)
			if rev1 != nil && !rev1.IsActive {
				s.ReassignReviewer(ctx, pr.ID, *pr.Reviewer1ID)
			}
		}
		if pr.Reviewer2ID != nil {
			rev2, _ := s.store.GetUser(ctx, *pr.Reviewer2ID)
			if rev2 != nil && !rev2.IsActive {
				s.ReassignReviewer(ctx, pr.ID, *pr.Reviewer2ID)
			}
		}
	}
	return nil
}

// getActiveReviewers получает активных ревьюверов с кэшем
func (s *Service) getActiveReviewers(ctx context.Context, teamID, excludeID int) ([]*domain.User, error) {
	key := fmt.Sprintf("team:%d:active", teamID)
	if cached, ok := s.cache.Load(key); ok {
		return cached.([]*domain.User), nil
	}

	users, err := s.store.GetActiveTeamUsers(ctx, teamID)
	if err != nil {
		return nil, err
	}

	var filtered []*domain.User
	for _, u := range users {
		if u.ID != excludeID {
			filtered = append(filtered, u)
		}
	}

	s.cache.Store(key, filtered)
	return filtered, nil
}

// shuffleReviewers шаффлит с crypto/rand (изюминка)
func (s *Service) shuffleReviewers(users []*domain.User) {
	n := len(users)
	for i := n - 1; i > 0; i-- {
		jBig, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(jBig.Int64())
		users[i], users[j] = users[j], users[i]
	}
}

func (s *Service) invalidateCache(teamID int) {
	key := fmt.Sprintf("team:%d:active", teamID)
	s.cache.Delete(key)
}
