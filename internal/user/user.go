package user

import (
	"context"
	"errors"

	"github.com/riskibarqy/go-template/config"
	"github.com/riskibarqy/go-template/datatransfers"
	"github.com/riskibarqy/go-template/internal/data"
	"github.com/riskibarqy/go-template/internal/types"
	"github.com/riskibarqy/go-template/models"
	"github.com/riskibarqy/go-template/utils"
	"golang.org/x/crypto/bcrypt"
)

// Errors
var (
	ErrWrongPassword      = errors.New("wrong password")
	ErrWrongEmail         = errors.New("wrong email")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNotFound           = errors.New("not found")
	ErrNoInput            = errors.New("no input")
	ErrLimitInput         = errors.New("name should be more than 5 char")
	ErrNameAlreadyExist   = errors.New(("name already exits"))
)

// Storage represents the user storage interface
type Storage interface {
	FindAll(ctx context.Context, params *datatransfers.FindAllParams) ([]*models.User, *types.Error)
	FindByID(ctx context.Context, userID int) (*models.User, *types.Error)
	FindByEmail(ctx context.Context, email string) (*models.User, *types.Error)
	FindByToken(ctx context.Context, token string) (*models.User, *types.Error)
	Insert(ctx context.Context, user *models.User) (*models.User, *types.Error)
	Update(ctx context.Context, user *models.User) (*models.User, *types.Error)
	Delete(ctx context.Context, userID int) *types.Error
}

// ServiceInterface represents the user service interface
type ServiceInterface interface {
	ListUsers(ctx context.Context, params *datatransfers.FindAllParams) ([]*models.User, int, *types.Error)
	GetUser(ctx context.Context, userID int) (*models.User, *types.Error)
	CreateUser(ctx context.Context, params *models.User) (*models.User, *types.Error)
	UpdateUser(ctx context.Context, userID int, params *models.User) (*models.User, *types.Error)
	DeleteUser(ctx context.Context, userID int) *types.Error
	ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) *types.Error
	Login(ctx context.Context, email string, password string) (*datatransfers.LoginResponse, *types.Error)
	Logout(ctx context.Context, token string) *types.Error
	GetByToken(ctx context.Context, token string) (*models.User, *types.Error)
}

// Service is the domain logic implementation of user Service interface
type Service struct {
	userStorage Storage
}

// ListUsers is listing users
func (s *Service) ListUsers(ctx context.Context, params *datatransfers.FindAllParams) ([]*models.User, int, *types.Error) {
	users, err := s.userStorage.FindAll(ctx, params)
	if err != nil {
		err.Path = ".UserService->ListUsers()" + err.Path
		return nil, 0, err
	}
	params.Page = 0
	params.Limit = 0
	allUsers, err := s.userStorage.FindAll(ctx, params)
	if err != nil {
		err.Path = ".UserService->ListUsers()" + err.Path
		return nil, 0, err
	}

	return users, len(allUsers), nil
}

// GetUser is get user
func (s *Service) GetUser(ctx context.Context, userID int) (*models.User, *types.Error) {
	user, err := s.userStorage.FindByID(ctx, userID)
	if err != nil {
		err.Path = ".UserService->GetUser()" + err.Path
		return nil, err
	}

	return user, nil
}

// CreateUser create user
func (s *Service) CreateUser(ctx context.Context, params *models.User) (*models.User, *types.Error) {
	users, _, errType := s.ListUsers(ctx, &datatransfers.FindAllParams{
		Email: params.Email,
	})
	if errType != nil {
		errType.Path = ".UserService->CreateUser()" + errType.Path
		return nil, errType
	}
	if len(users) > 0 {
		return nil, &types.Error{
			Path:    ".UserService->CreateUser()",
			Message: ErrEmailAlreadyExists.Error(),
			Error:   ErrEmailAlreadyExists,
			Type:    "validation-error",
		}
	}

	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &types.Error{
			Path:    ".UserService->CreateUser()",
			Message: err.Error(),
			Error:   err,
			Type:    "golang-error",
		}
	}

	now := utils.Now()

	user := &models.User{
		Name:           params.Name,
		Email:          params.Email,
		Password:       string(bcryptHash),
		Token:          nil,
		TokenExpiredAt: nil,
		CreatedAt:      now,
		UpdatedAt:      &now,
	}

	user, errType = s.userStorage.Insert(ctx, user)
	if errType != nil {
		errType.Path = ".UserService->CreateUser()" + errType.Path
		return nil, errType
	}

	return user, nil
}

// UpdateUser update a user
func (s *Service) UpdateUser(ctx context.Context, userID int, params *models.User) (*models.User, *types.Error) {
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		err.Path = ".UserService->UpdateUser()" + err.Path
		return nil, err
	}

	users, _, err := s.ListUsers(ctx, &datatransfers.FindAllParams{
		Email: params.Email,
	})
	if err != nil {
		err.Path = ".UserService->UpdateUser()" + err.Path
		return nil, err
	}
	if len(users) > 0 {
		return nil, &types.Error{
			Path:    ".UserService->CreateUser()",
			Message: data.ErrAlreadyExist.Error(),
			Error:   data.ErrAlreadyExist,
			Type:    "validation-error",
		}
	}

	user.Name = params.Name
	user.Email = params.Email

	user, err = s.userStorage.Update(ctx, user)
	if err != nil {
		err.Path = ".UserService->UpdateUser()" + err.Path
		return nil, err
	}

	return user, nil
}

// DeleteUser delete a user
func (s *Service) DeleteUser(ctx context.Context, userID int) *types.Error {
	err := s.userStorage.Delete(ctx, userID)
	if err != nil {
		err.Path = ".UserService->DeleteUser()" + err.Path
		return err
	}

	return nil
}

// ChangePassword change password
func (s *Service) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) *types.Error {
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		err.Path = ".UserService->ChangePassword()" + err.Path
		return err
	}

	errBcrypt := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if errBcrypt != nil {
		return &types.Error{
			Path:    ".UserService->ChangePassword()",
			Message: ErrWrongPassword.Error(),
			Error:   ErrWrongPassword,
			Type:    "golang-error",
		}
	}

	bcryptHash, errBcrypt := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if errBcrypt != nil {
		return &types.Error{
			Path:    ".UserService->ChangePassword()",
			Message: errBcrypt.Error(),
			Error:   errBcrypt,
			Type:    "golang-error",
		}
	}

	user.Password = string(bcryptHash)
	_, err = s.userStorage.Update(ctx, user)
	if err != nil {
		err.Path = ".UserService->ChangePassword()" + err.Path
		return err
	}

	return nil
}

// Login login
func (s *Service) Login(ctx context.Context, email string, password string) (*datatransfers.LoginResponse, *types.Error) {
	users, err := s.userStorage.FindAll(ctx, &datatransfers.FindAllParams{
		Email: email,
	})
	if err != nil {
		err.Path = ".UserService->Login()" + err.Path
		return nil, err
	}
	if len(users) < 1 {
		return nil, &types.Error{
			Path:    ".UserService->Login()",
			Message: ErrWrongEmail.Error(),
			Error:   ErrWrongEmail,
			Type:    "validation-error",
		}
	}

	user := users[0]
	errBcrypt := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errBcrypt != nil {
		return nil, &types.Error{
			Path:    ".UserService->ChangePassword()",
			Message: ErrWrongPassword.Error(),
			Error:   ErrWrongPassword,
			Type:    "golang-error",
		}
	}

	token, errToken := config.GenerateJWTToken(user)
	if errToken != nil {
		return nil, &types.Error{
			Path:    ".UserService->CreateUser()",
			Message: errToken.Error(),
			Error:   errToken,
			Type:    "golang-error",
		}
	}

	now := utils.Now()
	tokenExpiredAt := now + 72*3600

	user.Token = &token
	user.TokenExpiredAt = &tokenExpiredAt
	user.UpdatedAt = &now

	user, err = s.userStorage.Update(ctx, user)
	if err != nil {
		err.Path = ".UserService->CreateUser()" + err.Path
		return nil, err
	}

	return &datatransfers.LoginResponse{
		SessionID: token,
		User:      user,
	}, nil
}

// Logout logout
func (s *Service) Logout(ctx context.Context, token string) *types.Error {
	user, err := s.userStorage.FindByToken(ctx, token)
	if err != nil {
		err.Path = ".UserService->Logout()" + err.Path
		return err
	}

	user.Token = nil
	user.TokenExpiredAt = nil
	_, err = s.userStorage.Update(ctx, user)
	if err != nil {
		err.Path = ".UserService->Logout()" + err.Path
		return err
	}

	return nil
}

// GetByToken get user by its token
func (s *Service) GetByToken(ctx context.Context, token string) (*models.User, *types.Error) {
	user, err := s.userStorage.FindByToken(ctx, token)
	if err != nil {
		err.Path = ".UserService->GetByToken()" + err.Path
		return nil, err
	}

	return user, nil
}

// NewService creates a new user AppService
func NewService(
	userStorage Storage,
) *Service {
	return &Service{
		userStorage: userStorage,
	}
}
