package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/configs"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"github.com/imraushankr/bervity/server/src/internal/utils"
)

type AuthHandler struct {
	service interfaces.AuthService
	cfg     *configs.Config
	log     logger.Logger
}

func NewAuthHandler(service interfaces.AuthService, cfg *configs.Config, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		cfg:     cfg,
		log:     log,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.GetValidationErrors(err))
		return
	}

	user, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case models.ErrEmailAlreadyExists, models.ErrUsernameAlreadyExists:
			utils.Error(c, http.StatusConflict, "Registration failed", err)
		default:
			utils.Error(c, http.StatusInternalServerError, "Registration failed", err)
		}
		return
	}

	user.Sanitize()
	utils.Success(c, http.StatusCreated, "User registered successfully", user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.GetValidationErrors(err))
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &req)
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

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"access_token",
		resp.AccessToken,
		resp.ExpiresIn,
		"/",
		"",
		h.cfg.JWT.SecureCookie,
		true,
	)
	c.SetCookie(
		"logged_in",
		"true",
		resp.ExpiresIn,
		"/",
		"",
		h.cfg.JWT.SecureCookie,
		false,
	)

	utils.Success(c, http.StatusOK, "Login successful", resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", h.cfg.JWT.SecureCookie, true)
	c.SetCookie("logged_in", "", -1, "/", "", h.cfg.JWT.SecureCookie, false)
	utils.Success(c, http.StatusOK, "Logged out successfully", nil)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if err := h.service.VerifyEmail(c.Request.Context(), token); err != nil {
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

func (h *AuthHandler) InitiatePasswordReset(c *gin.Context) {
	var req models.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.GetValidationErrors(err))
		return
	}

	if err := h.service.InitiatePasswordReset(c.Request.Context(), req.Email); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to initiate password reset", err)
		return
	}

	utils.Success(c, http.StatusOK, "If an account exists, a reset link has been sent", nil)
}

func (h *AuthHandler) CompletePasswordReset(c *gin.Context) {
	token := c.Param("token")
	var req models.CompletePasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.GetValidationErrors(err))
		return
	}

	if err := h.service.CompletePasswordReset(c.Request.Context(), token, req.NewPassword); err != nil {
		switch err {
		case models.ErrInvalidResetToken:
			utils.Error(c, http.StatusBadRequest, "Invalid or expired token", err)
		default:
			utils.Error(c, http.StatusInternalServerError, "Password reset failed", err)
		}
		return
	}

	utils.Success(c, http.StatusOK, "Password reset successfully", nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.GetValidationErrors(err))
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "Token refresh failed", err)
		return
	}

	utils.Success(c, http.StatusOK, "Token refreshed successfully", resp)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.GetValidationErrors(err))
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		switch err {
		case models.ErrInvalidCredentials:
			utils.Error(c, http.StatusUnauthorized, "Current password is incorrect", err)
		default:
			utils.Error(c, http.StatusInternalServerError, "Password change failed", err)
		}
		return
	}

	utils.Success(c, http.StatusOK, "Password changed successfully", nil)
}