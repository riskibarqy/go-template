package postgres

import (
	"context"
	"fmt"

	"github.com/riskibarqy/99backend-challenge/datatransfers"
	"github.com/riskibarqy/99backend-challenge/internal/data"
	"github.com/riskibarqy/99backend-challenge/internal/types"
	"github.com/riskibarqy/99backend-challenge/models"
)

// PostgresStorage implements the user storage service interface
type PostgresStorage struct {
	Storage data.GenericStorage
}

// FindAll find all users
func (s *PostgresStorage) FindAll(ctx context.Context, params *datatransfers.FindAllParams) ([]*models.User, *types.Error) {

	users := []*models.User{}
	where := `"deleted_at" IS NULL`

	if params.UserID != 0 {
		where += ` AND "id" = :userId`
	}
	if params.Email != "" {
		where += ` AND "email" ILIKE :email`
	}
	if params.Name != "" {
		where += ` AND "name" ILIKE :name`
	}
	if params.Search != "" {
		where += ` AND "name" ILIKE :search`
	}
	if params.Token != "" {
		where += ` AND "token" ILIKE :token`
	}
	if len(params.UserIDs) > 0 {
		where += ` AND "id" in (:userIds)`
	}
	if params.Page != 0 && params.Limit != 0 {
		where = fmt.Sprintf(`%s ORDER BY "id" DESC LIMIT :limit OFFSET :offset`, where)
	} else {
		where = fmt.Sprintf(`%s ORDER BY "id" DESC`, where)
	}

	err := s.Storage.Where(ctx, &users, where, map[string]interface{}{
		"userId":  params.UserID,
		"userIds": params.UserIDs,
		"limit":   params.Limit,
		"email":   params.Email,
		"name":    params.Name,
		"search":  "%" + params.Search + "%",
		"offset":  ((params.Page - 1) * params.Limit),
		"token":   params.Token,
	})
	if err != nil {
		return nil, &types.Error{
			Path:    ".UserPostgresStorage->FindAll()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return users, nil
}

// FindByID find user by its id
func (s *PostgresStorage) FindByID(ctx context.Context, userID int) (*models.User, *types.Error) {
	users, err := s.FindAll(ctx, &datatransfers.FindAllParams{
		UserID: userID,
	})
	if err != nil {
		err.Path = ".UserPostgresStorage->FindByID()" + err.Path
		return nil, err
	}

	if len(users) < 1 || users[0].ID != userID {
		return nil, &types.Error{
			Path:    ".UserPostgresStorage->FindByID()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "pq-error",
		}
	}

	return users[0], nil
}

// FindByEmail find user by its email
func (s *PostgresStorage) FindByEmail(ctx context.Context, email string) (*models.User, *types.Error) {
	users, err := s.FindAll(ctx, &datatransfers.FindAllParams{
		Email: email,
	})
	if err != nil {
		err.Path = ".UserPostgresStorage->FindByEmail()" + err.Path
		return nil, err
	}

	if len(users) < 1 || users[0].Email != email {
		return nil, &types.Error{
			Path:    ".UserPostgresStorage->FindByEmail()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "pq-error",
		}
	}

	return users[0], nil
}

// FindByToken find user by its token
func (s *PostgresStorage) FindByToken(ctx context.Context, token string) (*models.User, *types.Error) {
	users, err := s.FindAll(ctx, &datatransfers.FindAllParams{
		Token: token,
	})
	if err != nil {
		err.Path = ".UserPostgresStorage->FindByToken()" + err.Path
		return nil, err
	}

	if len(users) < 1 || (users[0].Token != nil && *users[0].Token != token) {
		return nil, &types.Error{
			Path:    ".UserPostgresStorage->FindByToken()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "pq-error",
		}
	}

	return users[0], nil
}

// Insert insert user
func (s *PostgresStorage) Insert(ctx context.Context, user *models.User) (*models.User, *types.Error) {
	err := s.Storage.Insert(ctx, user)
	if err != nil {
		return nil, &types.Error{
			Path:    ".UserPostgresStorage->Insert()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return user, nil
}

// Update update user
func (s *PostgresStorage) Update(ctx context.Context, user *models.User) (*models.User, *types.Error) {
	err := s.Storage.Update(ctx, user)
	if err != nil {
		return nil, &types.Error{
			Path:    ".UserPostgresStorage->Update()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return user, nil
}

// Delete delete a user
func (s *PostgresStorage) Delete(ctx context.Context, userID int) *types.Error {
	err := s.Storage.Delete(ctx, userID)
	if err != nil {
		return &types.Error{
			Path:    ".UserPostgresStorage->Delete()",
			Message: err.Error(),
			Error:   err,
			Type:    "pq-error",
		}
	}

	return nil
}

// NewPostgresStorage creates new user repository service
func NewPostgresStorage(
	storage data.GenericStorage,
) *PostgresStorage {
	return &PostgresStorage{
		Storage: storage,
	}
}
