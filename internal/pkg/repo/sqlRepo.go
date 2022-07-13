package repo

import (
	"context"
	"database/sql"

	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/domain"
	"github.com/sirupsen/logrus"
)

const componentSQLRepo = "SQLRepo"

type SQlRepo struct {
	db     *sql.DB
	logger *logrus.Entry
}

func NewSQLRepo(db *sql.DB, logger *logrus.Logger) *SQlRepo {
	return &SQlRepo{
		db:     db,
		logger: logger.WithField("component", componentSQLRepo),
	}
}

func (r *SQlRepo) CreateUser(ctx context.Context, user domain.User) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)`,
		user.ID, user.Attributes.UserName, user.Attributes.Email, user.Password,
	)
	if err != nil {
		r.logger.WithError(err).WithField("method", "CreateUser").Error("failed to create user")
		return err
	}

	return nil
}

func (r *SQlRepo) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, email, password FROM users WHERE id = ?`,
		userID,
	).Scan(
		&user.ID, &user.Attributes.UserName, &user.Attributes.Email, &user.Password,
	)
	if err != nil {
		r.logger.WithError(err).WithField("method", "GetUser").Error("failed to get user")
		return nil, err
	}

	return &user, nil
}

func (r *SQlRepo) UpdateUser(ctx context.Context, user domain.User) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?`,
		user.Attributes.UserName, user.Attributes.Email, user.Password, user.ID,
	)
	if err != nil {
		r.logger.WithError(err).WithField("method", "UpdateUser").Error("failed to update user")
		return err
	}

	return nil
}
