package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/riskibarqy/go-template/datatransfers"
	"github.com/riskibarqy/go-template/internal/appcontext"
	"github.com/riskibarqy/go-template/internal/data"
	"github.com/riskibarqy/go-template/internal/http/response"
	"github.com/riskibarqy/go-template/internal/types"
	"github.com/riskibarqy/go-template/internal/user"
	"github.com/riskibarqy/go-template/models"
)

// UserController represents the user controller
type UserController struct {
	userService user.ServiceInterface
	dataManager *data.Manager
}

// UserList user list and count
type UserList struct {
	Data  []*models.User `json:"data"`
	Count int            `json:"count"`
}

func (a *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)

	var params datatransfers.LoginParams
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".UserController->Login()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	var sess *datatransfers.LoginResponse
	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		sess, err = a.userService.Login(r.Context(), params.Email, params.Password)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".UserController->Login()" + err.Path
		if err.Error == user.ErrWrongPassword || err.Error == data.ErrNotFound {
			response.Error(w, "Email / password is wrong", http.StatusBadRequest, *err)
		} else {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "sessionId",
		Value: sess.SessionID,
	})

	response.JSON(w, http.StatusOK, sess)
}

func (a *UserController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)
	var params datatransfers.ChangePasswordParams
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".UserController->ChangePassword()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	userID := appcontext.UserID(r.Context())

	err = a.userService.ChangePassword(r.Context(), userID, params.OldPassword, params.NewPassword)
	if err != nil {
		err.Path = ".UserController->ChangePassword()" + err.Path
		if err.Error == user.ErrWrongPassword {
			response.Error(w, "Wrong old password", http.StatusBadRequest, *err)
		} else {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		}
		return
	}

	response.JSON(w, http.StatusNoContent, "")
}

func (a *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)

	var params *models.User
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".UserController->UpdateUser()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}
	var sUserID = chi.URLParam(r, "userId")
	userID, errConversion := strconv.Atoi(sUserID)
	if errConversion != nil {
		err = &types.Error{
			Path:    ".UserController->UpdateUser()",
			Message: errConversion.Error(),
			Error:   errConversion,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	var singleUser *models.User
	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		singleUser, err = a.userService.UpdateUser(ctx, userID, params)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".UserController->UpdateUser()" + err.Path
		if errTransaction == user.ErrEmailAlreadyExists {
			response.Error(w, "email has been registered", http.StatusUnprocessableEntity, *err)
		} else {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		}
		return
	}
	response.JSON(w, http.StatusOK, singleUser)

}

func (a *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	errParseForm := r.ParseForm()
	if errParseForm != nil {
		err = &types.Error{
			Path:    ".UserController->CreateUser()",
			Message: errParseForm.Error(),
			Error:   errParseForm,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	params := &models.User{
		Email:    r.FormValue("email"),
		Name:     r.FormValue("name"),
		Password: r.FormValue("password"),
	}

	requiredFields := map[string]string{
		"Name":     "user name is required",
		"Email":    "user email is required",
		"Password": "user password is required",
	}

	if !a.validateRequiredFields(w, params, requiredFields) {
		return
	}

	result := &models.User{}
	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		result, err = a.userService.CreateUser(ctx, params)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".UserController->CreateUser()" + err.Path
		if errTransaction == user.ErrEmailAlreadyExists {
			response.Error(w, "email has been registered", http.StatusUnprocessableEntity, *err)
		} else {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		}

		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (a *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var err *types.Error
	var sUserID = chi.URLParam(r, "userId")
	userID, errConversion := strconv.Atoi(sUserID)
	if errConversion != nil {
		err = &types.Error{
			Path:    ".UserController->DeleteUser()",
			Message: errConversion.Error(),
			Error:   errConversion,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		err = a.userService.DeleteUser(ctx, userID)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".USerController->DeleteUser()" + err.Path
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}
	response.JSON(w, http.StatusNoContent, "")

}

func (a *UserController) ListUser(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	queryValues := r.URL.Query()
	var pageSize = 10
	var errConversion error
	if queryValues.Get("page_size") != "" {
		pageSize, errConversion = strconv.Atoi(queryValues.Get("page_size"))
		if errConversion != nil {
			err = &types.Error{
				Path:    ".UserController->ListUser()",
				Message: errConversion.Error(),
				Error:   errConversion,
				Type:    "golang-error",
			}
			response.Error(w, "Bad Request", http.StatusBadRequest, *err)
			return
		}
	}

	var pageNum = 1
	if queryValues.Get("page_num") != "" {
		pageNum, errConversion = strconv.Atoi(queryValues.Get("page_num"))
		if errConversion != nil {
			err = &types.Error{
				Path:    ".UserController->ListUser()",
				Message: errConversion.Error(),
				Error:   errConversion,
				Type:    "golang-error",
			}
			response.Error(w, "Bad Request", http.StatusBadRequest, *err)
			return
		}
	}

	var search = queryValues.Get("search")

	if pageSize < 0 {
		pageSize = 10
	}
	if pageNum < 0 {
		pageNum = 1
	}
	userList, count, err := a.userService.ListUsers(r.Context(), &datatransfers.FindAllParams{
		Limit:  pageSize,
		Search: search,
		Page:   pageNum,
	})
	if err != nil {
		err.Path = ".UserController->ListUser()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}
	if userList == nil {
		userList = []*models.User{}
	}

	response.JSON(w, http.StatusOK, UserList{
		Data:  userList,
		Count: count,
	})
}

func (a *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	// get token from the context
	// log it out!
	loginToken, ok := r.Context().Value(appcontext.KeySessionID).(string)
	if !ok {
		errUserID := errors.New("failed to get user id from request context")
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, types.Error{
			Path:    ".UserController->Logout()",
			Message: errUserID.Error(),
			Error:   errUserID,
			Type:    "golang-error",
		})
		return
	}

	err = a.userService.Logout(r.Context(), loginToken)
	if err != nil {
		err.Path = ".UserController->Logout()" + err.Path
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	response.JSON(w, http.StatusNoContent, "")
}

func (a *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	var sUserID = chi.URLParam(r, "userId")
	userID, errConversion := strconv.Atoi(sUserID)
	if errConversion != nil {
		err = &types.Error{
			Path:    ".UserController->UpdateUser()",
			Message: errConversion.Error(),
			Error:   errConversion,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	user, err := a.userService.GetUser(r.Context(), userID)
	if err != nil {
		err.Path = ".UserController->GetUserByID()" + err.Path
		response.Error(w, "User Not Found", http.StatusNotFound, *err)
		if err.Error != data.ErrNotFound {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
		return
	}

	response.JSON(w, http.StatusOK, user)

}

func (a *UserController) validateRequiredFields(w http.ResponseWriter, params *models.User, fields map[string]string) bool {
	for field, message := range fields {
		value := reflect.ValueOf(params).Elem().FieldByName(field).String()
		if value == "" {
			errValidation := errors.New(message)
			err := &types.Error{
				Path:    ".UserController->CreateUser()",
				Message: errValidation.Error(),
				Error:   errValidation,
				Type:    "golang-error",
			}
			response.Error(w, message, http.StatusBadRequest, *err)
			return false
		}
	}
	return true
}

// NewUserController creates a new user controller
func NewUserController(
	userService user.ServiceInterface,
	dataManager *data.Manager,
) *UserController {
	return &UserController{
		userService: userService,
		dataManager: dataManager,
	}
}
