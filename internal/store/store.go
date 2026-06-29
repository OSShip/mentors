package store

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/OSShip/mentors/internal/model"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) HasPendingApplication(ctx context.Context, userID string) (bool, error) {
	var id string
	err := s.pool.QueryRow(ctx,
		`SELECT id FROM mentor_applications WHERE user_id=$1 AND status='pending'`, userID).Scan(&id)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *Store) CreateApplication(ctx context.Context, userID string, githubData json.RawMessage) (model.Application, error) {
	id := uuid.New().String()
	_, err := s.pool.Exec(ctx,
		`INSERT INTO mentor_applications (id, user_id, github_data) VALUES ($1,$2,$3)`,
		id, userID, githubData)
	if err != nil {
		return model.Application{}, err
	}
	return model.Application{ID: id, UserID: userID, Status: "pending", GithubData: githubData}, nil
}

func (s *Store) ListApplications(ctx context.Context, status string) ([]model.Application, error) {
	q := `SELECT a.id, a.user_id, a.status, a.github_data, a.created_at,
		COALESCE(u.email, ''), COALESCE(u.display_name, ''), COALESCE(u.github_username, '')
		FROM mentor_applications a
		JOIN users u ON u.id = a.user_id`
	args := []interface{}{}
	if status != "" {
		q += ` WHERE a.status=$1`
		args = append(args, status)
	}
	q += ` ORDER BY a.created_at DESC`

	rows, err := s.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Application
	for rows.Next() {
		var a model.Application
		var createdAt interface{}
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.Status, &a.GithubData, &createdAt,
			&a.ApplicantEmail, &a.ApplicantDisplayName, &a.ApplicantGithubUsername,
		); err != nil {
			continue
		}
		list = append(list, a)
	}
	return list, nil
}

func (s *Store) ReviewApplication(ctx context.Context, id, status, adminID string) (string, error) {
	var userID string
	err := s.pool.QueryRow(ctx,
		`UPDATE mentor_applications SET status=$1, reviewed_by=$2, reviewed_at=NOW() WHERE id=$3 RETURNING user_id`,
		status, adminID, id).Scan(&userID)
	return userID, err
}

func (s *Store) PromoteToMentor(ctx context.Context, userID string) error {
	_, err := s.pool.Exec(ctx, `UPDATE users SET role='mentor' WHERE id=$1`, userID)
	return err
}

func (s *Store) GetUserEmail(ctx context.Context, userID string) (string, error) {
	var email string
	err := s.pool.QueryRow(ctx, `SELECT email FROM users WHERE id=$1`, userID).Scan(&email)
	return email, err
}
