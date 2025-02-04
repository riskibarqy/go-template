package datatransfers

import "github.com/riskibarqy/go-template/models"

// LoginParams represent the http request data for login user
type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response of login function
type LoginResponse struct {
	SessionID string       `json:"sessionId"`
	User      *models.User `json:"user"`
}

// ChangePasswordParams represent the http request data for change password
type ChangePasswordParams struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
