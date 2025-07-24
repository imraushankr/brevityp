// package v1

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/imraushankr/bervity/server/src/configs"
// 	"github.com/imraushankr/bervity/server/src/internal/models"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// 	"github.com/imraushankr/bervity/server/src/internal/services"
// 	"github.com/imraushankr/bervity/server/src/internal/utils"
// )

// type AuthHandler struct {
// 	userService services.UserService
// 	log         logger.Logger
// }

// func NewAuthHandler(userService services.UserService, log logger.Logger) *AuthHandler {
// 	return &AuthHandler{
// 		userService: userService,
// 		log:         log,
// 	}
// }

// func (h *AuthHandler) Register(c *gin.Context) {
// 	var req models.RegisterRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		h.log.Error("Failed to bind user data", logger.NamedError("error", err))
// 		utils.ValidationError(c, utils.GetValidationErrors(err))
// 		return
// 	}

// 	user := &models.User{
// 		FirstName: req.FirstName,
// 		LastName:  req.LastName,
// 		Username:  req.Username,
// 		Email:     req.Email,
// 		Password:  req.Password,
// 	}

// 	if err := h.userService.Register(c.Request.Context(), user); err != nil {
// 		switch err {
// 		case models.ErrEmailAlreadyExists, models.ErrUsernameAlreadyExists:
// 			utils.Error(c, http.StatusConflict, "Registration failed", err)
// 		default:
// 			utils.Error(c, http.StatusInternalServerError, "Registration failed", err)
// 		}
// 		return
// 	}

// 	user.Sanitize()
// 	utils.Success(c, http.StatusCreated, "User registered successfully. Please check your email for verification.", user)
// }

// func (h *AuthHandler) Login(c *gin.Context) {
// 	var loginRequest models.LoginRequest
// 	if err := c.ShouldBindJSON(&loginRequest); err != nil {
// 		h.log.Error("Failed to bind login data", logger.NamedError("error", err))
// 		utils.ValidationError(c, utils.GetValidationErrors(err))
// 		return
// 	}

// 	user, token, err := h.userService.Login(c.Request.Context(), loginRequest.UserID, loginRequest.Password)
// 	if err != nil {
// 		switch err {
// 		case models.ErrInvalidCredentials:
// 			utils.Error(c, http.StatusUnauthorized, "Login failed", err)
// 		case models.ErrUserNotVerified:
// 			utils.Error(c, http.StatusForbidden, "Login failed", err)
// 		default:
// 			utils.Error(c, http.StatusInternalServerError, "Login failed", err)
// 		}
// 		return
// 	}

// 	// Set secure cookies using configuration
// 	c.SetSameSite(http.SameSiteLaxMode)
// 	c.SetCookie(
// 		"access_token",
// 		token,
// 		int(configs.BrevityApp.JWT.AccessTokenExpiry.Seconds()),
// 		"/",
// 		"",
// 		configs.BrevityApp.JWT.SecureCookie,
// 		true, // HttpOnly
// 	)

// 	c.SetCookie(
// 		"logged_in",
// 		"true",
// 		int(configs.BrevityApp.JWT.AccessTokenExpiry.Seconds()),
// 		"/",
// 		"",
// 		configs.BrevityApp.JWT.SecureCookie,
// 		false, // Not HttpOnly
// 	)

// 	// Set refresh token cookie if your system uses it
// 	/*
// 		c.SetCookie(
// 			"refresh_token",
// 			refreshToken,
// 			int(configs.BrevityApp.JWT.RefreshTokenExpiry.Seconds()),
// 			"/auth/refresh",
// 			"",
// 			configs.BrevityApp.JWT.SecureCookie,
// 			true,
// 		)
// 	*/

// 	user.Sanitize()
// 	response := models.LoginResponse{
// 		User:        *user,
// 		AccessToken: token,
// 		TokenType:   "Bearer",
// 		ExpiresIn:   int(configs.BrevityApp.JWT.AccessTokenExpiry.Seconds()),
// 	}
// 	utils.Success(c, http.StatusOK, "Login successful", response)
// }

// func (h *AuthHandler) Logout(c *gin.Context) {
// 	// Clear all auth cookies
// 	c.SetCookie("access_token", "", -1, "/", "", configs.BrevityApp.JWT.SecureCookie, true)
// 	c.SetCookie("logged_in", "", -1, "/", "", configs.BrevityApp.JWT.SecureCookie, false)

// 	// Clear refresh token cookie if used
// 	/*
// 		c.SetCookie(
// 			"refresh_token",
// 			"",
// 			-1,
// 			"/auth/refresh",
// 			"",
// 			configs.BrevityApp.JWT.SecureCookie,
// 			true,
// 		)
// 	*/

// 	utils.Success(c, http.StatusOK, "Logged out successfully", nil)
// }

// func (h *AuthHandler) VerifyEmail(c *gin.Context) {
// 	token := c.Query("token")
// 	if token == "" {
// 		utils.Error(c, http.StatusBadRequest, "Verification failed", models.ErrInvalidToken)
// 		return
// 	}

// 	if err := h.userService.VerifyEmail(c.Request.Context(), token); err != nil {
// 		switch err {
// 		case models.ErrInvalidVerificationToken:
// 			utils.Error(c, http.StatusBadRequest, "Verification failed", err)
// 		default:
// 			utils.Error(c, http.StatusInternalServerError, "Verification failed", err)
// 		}
// 		return
// 	}

// 	utils.Success(c, http.StatusOK, "Email verified successfully", nil)
// }

// func (h *AuthHandler) InitiatePasswordReset(c *gin.Context) {
// 	var request models.PasswordResetRequest
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		h.log.Error("Failed to bind email data", logger.NamedError("error", err))
// 		utils.ValidationError(c, utils.GetValidationErrors(err))
// 		return
// 	}

// 	if err := h.userService.InitiatePasswordReset(c.Request.Context(), request.Email); err != nil {
// 		utils.Error(c, http.StatusInternalServerError, "Failed to initiate password reset", err)
// 		return
// 	}

// 	utils.Success(c, http.StatusOK, "Password reset email sent if account exists", nil)
// }

// /*
// func (h *AuthHandler) CompletePasswordReset(c *gin.Context) {
// 	var request models.CompletePasswordResetRequest
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		h.log.Error("Failed to bind reset data", logger.NamedError("error", err))
// 		utils.ValidationError(c, utils.GetValidationErrors(err))
// 		return
// 	}

// 	if err := h.userService.CompletePasswordReset(c.Request.Context(), request.Token, request.NewPassword); err != nil {
// 		switch err {
// 		case models.ErrInvalidResetToken:
// 			utils.Error(c, http.StatusBadRequest, "Password reset failed", err)
// 		default:
// 			utils.Error(c, http.StatusInternalServerError, "Password reset failed", err)
// 		}
// 		return
// 	}

// 	utils.Success(c, http.StatusOK, "Password reset successfully", nil)
// }
// *
// func (h *AuthHandler) CompletePasswordReset(c *gin.Context) {
// 	var request models.CompletePasswordResetRequest
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		h.log.Error("Failed to bind reset data",
// 			logger.NamedError("error", err),
// 			logger.Any("request_body", c.Request.Body))
// 		utils.ValidationError(c, utils.GetValidationErrors(err))
// 		return
// 	}

// 	h.log.Debug("Password reset request received",
// 		logger.String("token", request.Token),
// 		logger.Bool("has_new_password", request.NewPassword != ""))

// 	if err := h.userService.CompletePasswordReset(c.Request.Context(), request.Token, request.NewPassword); err != nil {
// 		switch err {
// 		case models.ErrInvalidResetToken:
// 			utils.Error(c, http.StatusBadRequest, "Password reset failed", err)
// 		default:
// 			utils.Error(c, http.StatusInternalServerError, "Password reset failed", err)
// 		}
// 		return
// 	}

// 	utils.Success(c, http.StatusOK, "Password reset successfully", nil)
// }
// */

// func (h *AuthHandler) CompletePasswordReset(c *gin.Context) {
//     // Get token from URL parameter
//     token := c.Param("token")
//     if token == "" {
//         utils.Error(c, http.StatusBadRequest, "Token is required", nil)
//         return
//     }

//     var request models.CompletePasswordResetRequest
//     if err := c.ShouldBindJSON(&request); err != nil {
//         h.log.Error("Failed to bind password reset data",
//             logger.NamedError("error", err),
//             logger.String("token", token))
//         utils.ValidationError(c, utils.GetValidationErrors(err))
//         return
//     }

//     // Set the token from URL param
//     request.Token = token

//     h.log.Debug("Processing password reset",
//         logger.String("token_prefix", token[:8]),
//         logger.Bool("has_password", request.NewPassword != ""))

//     if err := h.userService.CompletePasswordReset(c.Request.Context(), request.Token, request.NewPassword); err != nil {
//         h.log.Error("Password reset failed",
//             logger.NamedError("error", err),
//             logger.String("token_prefix", token[:8]))

//         switch err {
//         case models.ErrInvalidResetToken:
//             utils.Error(c, http.StatusBadRequest, "Invalid or expired reset token", err)
//         default:
//             utils.Error(c, http.StatusInternalServerError, "Password reset failed", err)
//         }
//         return
//     }

//     utils.Success(c, http.StatusOK, "Password reset successfully", nil)
// }

// func (h *AuthHandler) RefreshToken(c *gin.Context) {
// 	var request models.RefreshTokenRequest
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		h.log.Error("Failed to bind refresh token", logger.NamedError("error", err))
// 		utils.ValidationError(c, utils.GetValidationErrors(err))
// 		return
// 	}

// 	newToken, err := h.userService.RefreshToken(c.Request.Context(), request.RefreshToken)
// 	if err != nil {
// 		utils.Error(c, http.StatusUnauthorized, "Token refresh failed", err)
// 		return
// 	}

// 	response := models.RefreshTokenResponse{
// 		AccessToken: newToken,
// 		TokenType:   "Bearer",
// 		ExpiresIn:   3600, // 1 hour
// 	}
// 	utils.Success(c, http.StatusOK, "Token refreshed successfully", response)
// }

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