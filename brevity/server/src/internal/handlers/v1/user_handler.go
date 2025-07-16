package v1

import (
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"github.com/imraushankr/bervity/server/src/internal/services"
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
