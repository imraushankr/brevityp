// package v1

// import (
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// 	"github.com/imraushankr/bervity/server/src/internal/services"
// )

// type UserHandler struct {
// 	userService services.UserService
// 	log         logger.Logger
// }

// func NewUserHandler(userService services.UserService, log logger.Logger) *UserHandler {
// 	return &UserHandler{
// 		userService: userService,
// 		log:         log,
// 	}
// }

package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"github.com/imraushankr/bervity/server/src/internal/services"
	"github.com/imraushankr/bervity/server/src/internal/utils"
)

type UserHandler struct {
	userService services.UserService
	log         logger.Logger
}

func NewUserHandler(userService services.UserService, log logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		log:         log,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.log.Error("Failed to bind user data", logger.NamedError("error", err))
		utils.ValidationError(c, utils.GetValidationErrors(err))
		return
	}

	if err := h.userService.Register(c.Request.Context(), &user); err != nil {
		switch err {
		case models.ErrEmailAlreadyExists, models.ErrUsernameAlreadyExists:
			utils.Error(c, http.StatusConflict, "Registration failed", err)
		default:
			utils.Error(c, http.StatusInternalServerError, "Registration failed", err)
		}
		return
	}

	// Clear sensitive data before sending response
	user.Password = ""
	utils.Success(c, http.StatusCreated, "User registered successfully. Please check your email for verification.", user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginRequest struct {
		UserID   string `json:"userId" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		h.log.Error("Failed to bind login data", logger.NamedError("error", err))
		utils.ValidationError(c, utils.GetValidationErrors(err))
		return
	}

	user, token, err := h.userService.Login(c.Request.Context(), loginRequest.UserID, loginRequest.Password)
	if err != nil {
		switch err {
		case models.ErrInvalidCredentials:
			utils.Error(c, http.StatusUnauthorized, "Login failed", err)
		case models.ErrUserNotVerified:
			utils.Error(c, http.StatusForbidden, "Login failed", err)
		default:
			utils.Error(c, http.StatusInternalServerError, "Login failed", err)
		}
		return
	}

	// Clear sensitive data before sending response
	user.Password = ""
	response := gin.H{
		"user":  user,
		"token": token,
	}
	utils.Success(c, http.StatusOK, "Login successful", response)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	identifier := c.Param("identifier")

	user, err := h.userService.FindUser(c.Request.Context(), identifier)
	if err != nil {
		if err == models.ErrUserNotFound {
			utils.Error(c, http.StatusNotFound, "User not found", err)
		} else {
			utils.Error(c, http.StatusInternalServerError, "Failed to get user", err)
		}
		return
	}

	// Clear sensitive data before sending response
	user.Password = ""
	utils.Success(c, http.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.log.Error("Failed to bind user data", logger.NamedError("error", err))
		utils.ValidationError(c, utils.GetValidationErrors(err))
		return
	}

	user.ID = c.Param("id")
	if err := h.userService.UpdateUser(c.Request.Context(), &user); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	// Clear sensitive data before sending response
	user.Password = ""
	utils.Success(c, http.StatusOK, "User updated successfully", user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	utils.Success(c, http.StatusOK, "User deleted successfully", nil)
}

func (h *UserHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		utils.Error(c, http.StatusBadRequest, "Verification failed", models.ErrInvalidToken)
		return
	}

	if err := h.userService.VerifyEmail(c.Request.Context(), token); err != nil {
		switch err {
		case models.ErrInvalidVerificationToken:
			utils.Error(c, http.StatusBadRequest, "Verification failed", err)
		default:
			utils.Error(c, http.StatusInternalServerError, "Verification failed", err)
		}
		return
	}

	utils.Success(c, http.StatusOK, "Email verified successfully", nil)
}

func (h *UserHandler) InitiatePasswordReset(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Error("Failed to bind email data", logger.NamedError("error", err))
		utils.ValidationError(c, utils.GetValidationErrors(err))
		return
	}

	if err := h.userService.InitiatePasswordReset(c.Request.Context(), request.Email); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to initiate password reset", err)
		return
	}

	utils.Success(c, http.StatusOK, "Password reset email sent if account exists", nil)
}

func (h *UserHandler) CompletePasswordReset(c *gin.Context) {
	var request struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Error("Failed to bind reset data", logger.NamedError("error", err))
		utils.ValidationError(c, utils.GetValidationErrors(err))
		return
	}

	if err := h.userService.CompletePasswordReset(c.Request.Context(), request.Token, request.NewPassword); err != nil {
		switch err {
		case models.ErrInvalidResetToken:
			utils.Error(c, http.StatusBadRequest, "Password reset failed", err)
		default:
			utils.Error(c, http.StatusInternalServerError, "Password reset failed", err)
		}
		return
	}

	utils.Success(c, http.StatusOK, "Password reset successfully", nil)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Error("Failed to bind refresh token", logger.NamedError("error", err))
		utils.ValidationError(c, utils.GetValidationErrors(err))
		return
	}

	newToken, err := h.userService.RefreshToken(c.Request.Context(), request.RefreshToken)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "Token refresh failed", err)
		return
	}

	response := gin.H{
		"token": newToken,
	}
	utils.Success(c, http.StatusOK, "Token refreshed successfully", response)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := c.Param("id")
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Failed to get avatar file", err)
		return
	}
	defer file.Close()

	avatarURL, err := h.userService.UploadAvatar(c.Request.Context(), userID, file, header)
	if err != nil {
		switch err {
		case models.ErrFileTooLarge, models.ErrInvalidFileType:
			utils.Error(c, http.StatusBadRequest, "Avatar upload failed", err)
		default:
			utils.Error(c, http.StatusInternalServerError, "Avatar upload failed", err)
		}
		return
	}

	response := gin.H{
		"avatarUrl": avatarURL,
	}
	utils.Success(c, http.StatusOK, "Avatar uploaded successfully", response)
}
