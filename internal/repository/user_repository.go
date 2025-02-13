package repository

import (
	"context"
	"fmt"
	"go-scheduler/internal/repository/entity"
	"go-scheduler/logger"
	"go-scheduler/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UserRepositoryInterface interface {
	GetUsers(ctx context.Context) ([]*entity.GetUser, error)
	BatchUpdateUser(ctx context.Context, users []*model.UpdateUserRequest, workerID int, batchID int) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepositoryInterface {
	return &userRepository{db: db}
}

func (ur *userRepository) GetUsers(ctx context.Context) ([]*entity.GetUser, error) {
	users := []*entity.GetUser{}

	query := `
	SELECT username, full_name
	FROM c_user_copy cu
	WHERE 
		cu.username IS NOT NULL
		AND cu.full_name IS NOT NULL
	ORDER BY cu.full_name
	`

	err := ur.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *userRepository) BatchUpdateUser(ctx context.Context, users []*model.UpdateUserRequest, workerID int, batchID int) error {
	if len(users) == 0 {
		return nil
	}

	query := `
	UPDATE c_user_copy AS u
	SET full_name = updates.full_name
	FROM (VALUES %s) AS updates(username, full_name)
	WHERE 
		u.username = updates.username
		AND u.deleted_at IS NULL
	`
	values := []string{}
	args := []interface{}{}
	argID := 1
	for _, user := range users {
		values = append(values, fmt.Sprintf("($%d, $%d)", argID, argID+1))
		args = append(args, user.Username, user.FullName)
		argID += 2
	}
	finalQuery := fmt.Sprintf(query, strings.Join(values, ","))

	tx, err := ur.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.ErrorfCtx(ctx, "[Worker %d] Failed to begin transaction: %v", workerID, err)
		return err
	}

	_, err = tx.ExecContext(ctx, finalQuery, args...)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_ = tx.Commit()

	return nil
}
