package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"github.com/imraushankr/bervity/server/src/internal/utils"
)

type UserHandler struct {
	service interfaces.UserService
	log     logger.Logger
}

func NewUserHandler(service interfaces.UserService, log logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		log:     log,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.service.FindUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			utils.Error(c, http.StatusNotFound, "User not found", err)
		} else {
			utils.Error(c, http.StatusInternalServerError, "Failed to get profile", err)
		}
		return
	}

	user.Sanitize()
	utils.Success(c, http.StatusOK, "Profile retrieved successfully", user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.GetValidationErrors(err))
		return
	}

	user, err := h.service.FindUser(c.Request.Context(), userID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "User not found", err)
		return
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName

	if err := h.service.UpdateUser(c.Request.Context(), user); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update profile", err)
		return
	}

	user.Sanitize()
	utils.Success(c, http.StatusOK, "Profile updated successfully", user)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := c.GetString("user_id")
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid file upload", err)
		return
	}
	defer file.Close()

	avatarURL, err := h.service.UploadAvatar(c.Request.Context(), userID, file, header)
	if err != nil {
		switch err {
		case models.ErrFileTooLarge, models.ErrInvalidFileType:
			utils.Error(c, http.StatusBadRequest, "Invalid file", err)
		default:
			utils.Error(c, http.StatusInternalServerError, "Failed to upload avatar", err)
		}
		return
	}

	utils.Success(c, http.StatusOK, "Avatar uploaded successfully", models.UploadAvatarResponse{
		AvatarURL: avatarURL,
	})
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetString("user_id")
	if err := h.service.DeleteUser(c.Request.Context(), userID); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete account", err)
		return
	}

	utils.Success(c, http.StatusOK, "Account deleted successfully", nil)
}