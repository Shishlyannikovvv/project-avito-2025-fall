package store

import (
	"context"
	"errors"

	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store управляет хранилищем
type Store struct {
	pool *pgxpool.Pool
}

// NewStore создаёт стор
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// CreateTeam создаёт команду
func (st *Store) CreateTeam(ctx context.Context, name string) (*domain.Team, error) {
	var team domain.Team
	err := st.pool.QueryRow(ctx, `INSERT INTO teams (name) VALUES ($1) RETURNING id, name`, name).Scan(&team.ID, &team.Name)
	return &team, err
}

// CreateUser создаёт юзера
func (st *Store) CreateUser(ctx context.Context, username string, teamID int) (*domain.User, error) {
	var user domain.User
	err := st.pool.QueryRow(ctx, `INSERT INTO users (username, team_id) VALUES ($1, $2) RETURNING id, username, team_id, is_active`, username, teamID).Scan(&user.ID, &user.Username, &user.TeamID, &user.IsActive)
	return &user, err
}

// GetUser получает юзера
func (st *Store) GetUser(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	err := st.pool.QueryRow(ctx, `SELECT id, username, team_id, is_active FROM users WHERE id = $1`, id).Scan(&user.ID, &user.Username, &user.TeamID, &user.IsActive)
	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return &user, err
}

// CreatePR создаёт PR
func (st *Store) CreatePR(ctx context.Context, title string, authorID int) (*domain.PR, error) {
	var pr domain.PR
	err := st.pool.QueryRow(ctx, `INSERT INTO pull_requests (title, author_id) VALUES ($1, $2) RETURNING id, title, author_id, status, reviewer1_id, reviewer2_id, created_at, merged_at`, title, authorID).Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, &pr.Reviewer1ID, &pr.Reviewer2ID, &pr.CreatedAt, &pr.MergedAt)
	return &pr, err
}

// UpdatePRReviewers обновляет ревьюверов
func (st *Store) UpdatePRReviewers(ctx context.Context, prID int, rev1, rev2 *int) (*domain.PR, error) {
	var pr domain.PR
	err := st.pool.QueryRow(ctx, `UPDATE pull_requests SET reviewer1_id = $1, reviewer2_id = $2 WHERE id = $3 RETURNING id, title, author_id, status, reviewer1_id, reviewer2_id, created_at, merged_at`, rev1, rev2, prID).Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, &pr.Reviewer1ID, &pr.Reviewer2ID, &pr.CreatedAt, &pr.MergedAt)
	return &pr, err
}

// GetPR получает PR
func (st *Store) GetPR(ctx context.Context, id int) (*domain.PR, error) {
	var pr domain.PR
	err := st.pool.QueryRow(ctx, `SELECT id, title, author_id, status, reviewer1_id, reviewer2_id, created_at, merged_at FROM pull_requests WHERE id = $1`, id).Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, &pr.Reviewer1ID, &pr.Reviewer2ID, &pr.CreatedAt, &pr.MergedAt)
	if err == pgx.ErrNoRows {
		return nil, errors.New("PR not found")
	}
	return &pr, err
}

// MergePR мержит PR
func (st *Store) MergePR(ctx context.Context, id int) (*domain.PR, error) {
	var pr domain.PR
	err := st.pool.QueryRow(ctx, `UPDATE pull_requests SET status = 'MERGED', merged_at = NOW() WHERE id = $1 AND status != 'MERGED' RETURNING id, title, author_id, status, reviewer1_id, reviewer2_id, created_at, merged_at`, id).Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, &pr.Reviewer1ID, &pr.Reviewer2ID, &pr.CreatedAt, &pr.MergedAt)
	if err == pgx.ErrNoRows {
		return st.GetPR(ctx, id) // Идемпотентно возвращаем текущий
	}
	return &pr, err
}

// GetUserPRs получает PR юзера
func (st *Store) GetUserPRs(ctx context.Context, userID int) ([]*domain.PR, error) {
	rows, err := st.pool.Query(ctx, `SELECT id, title, author_id, status, reviewer1_id, reviewer2_id, created_at, merged_at FROM pull_requests WHERE reviewer1_id = $1 OR reviewer2_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*domain.PR
	for rows.Next() {
		var pr domain.PR
		if err := rows.Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, &pr.Reviewer1ID, &pr.Reviewer2ID, &pr.CreatedAt, &pr.MergedAt); err != nil {
			return nil, err
		}
		prs = append(prs, &pr)
	}
	return prs, nil
}

// GetActiveTeamUsers получает активных юзеров команды
func (st *Store) GetActiveTeamUsers(ctx context.Context, teamID int) ([]*domain.User, error) {
	rows, err := st.pool.Query(ctx, `SELECT id, username, team_id, is_active FROM users WHERE team_id = $1 AND is_active = true`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.TeamID, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

// GetReviewerStats статистика
func (st *Store) GetReviewerStats(ctx context.Context) (map[int]int, error) {
	rows, err := st.pool.Query(ctx, `SELECT reviewer_id, COUNT(*) FROM (SELECT reviewer1_id AS reviewer_id FROM pull_requests WHERE reviewer1_id IS NOT NULL UNION ALL SELECT reviewer2_id FROM pull_requests WHERE reviewer2_id IS NOT NULL) AS reviews GROUP BY reviewer_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[int]int)
	for rows.Next() {
		var id, count int
		if err := rows.Scan(&id, &count); err != nil {
			return nil, err
		}
		stats[id] = count
	}
	return stats, nil
}

// DeactivateUser деактивирует юзера
func (st *Store) DeactivateUser(ctx context.Context, id int) error {
	_, err := st.pool.Exec(ctx, `UPDATE users SET is_active = false WHERE id = $1`, id)
	return err
}

// GetTeamUsers получает юзеров команды
func (st *Store) GetTeamUsers(ctx context.Context, teamID int) ([]*domain.User, error) {
	rows, err := st.pool.Query(ctx, `SELECT id, username, team_id, is_active FROM users WHERE team_id = $1`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.TeamID, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

// GetOpenPRsByTeam получает открытые PR команды
func (st *Store) GetOpenPRsByTeam(ctx context.Context, teamID int) ([]*domain.PR, error) {
	rows, err := st.pool.Query(ctx, `
		SELECT pr.id, pr.title, pr.author_id, pr.status, pr.reviewer1_id, pr.reviewer2_id, pr.created_at, pr.merged_at
		FROM pull_requests pr
		JOIN users u ON pr.author_id = u.id
		WHERE u.team_id = $1 AND pr.status = 'OPEN'
	`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*domain.PR
	for rows.Next() {
		var pr domain.PR
		if err := rows.Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, &pr.Reviewer1ID, &pr.Reviewer2ID, &pr.CreatedAt, &pr.MergedAt); err != nil {
			return nil, err
		}
		prs = append(prs, &pr)
	}
	return prs, nil
}
