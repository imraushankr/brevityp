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