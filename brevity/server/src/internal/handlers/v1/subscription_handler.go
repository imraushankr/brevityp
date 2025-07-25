// package v1

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/imraushankr/bervity/server/src/internal/models"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// 	"github.com/imraushankr/bervity/server/src/internal/utils"
// )

// type SubscriptionHandler struct {
// 	subService interfaces.SubscriptionService
// 	log        logger.Logger
// }

// func NewSubscriptionHandler(subService interfaces.SubscriptionService, log logger.Logger) *SubscriptionHandler {
// 	return &SubscriptionHandler{
// 		subService: subService,
// 		log:        log,
// 	}
// }

// func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	userID := c.GetString("user_id")
// 	var req models.CreateSubscriptionRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		h.log.Debug("invalid request body", logger.ErrorField(err))
// 		utils.Error(c, http.StatusBadRequest, "Invalid request body", models.ErrInvalidInput)
// 		return
// 	}

// 	sub, err := h.subService.CreateSubscription(ctx, userID, &req)
// 	if err != nil {
// 		switch err {
// 		case models.ErrInvalidInput, models.ErrInvalidPlan:
// 			utils.Error(c, http.StatusBadRequest, err.Error(), err)
// 		case models.ErrActiveSubscriptionExists:
// 			utils.Error(c, http.StatusConflict, err.Error(), err)
// 		case models.ErrPaymentFailed:
// 			utils.Error(c, http.StatusPaymentRequired, err.Error(), err)
// 		default:
// 			h.log.Error("failed to create subscription", logger.ErrorField(err))
// 			utils.Error(c, http.StatusInternalServerError, "Failed to create subscription", err)
// 		}
// 		return
// 	}

// 	utils.Success(c, http.StatusCreated, "Subscription created successfully", sub)
// }

// func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	userID := c.GetString("user_id")

// 	sub, err := h.subService.GetUserSubscription(ctx, userID)
// 	if err != nil {
// 		if err == models.ErrSubscriptionNotActive {
// 			utils.Error(c, http.StatusNotFound, err.Error(), err)
// 			return
// 		}
// 		h.log.Error("failed to get subscription", logger.ErrorField(err))
// 		utils.Error(c, http.StatusInternalServerError, "Failed to get subscription", err)
// 		return
// 	}

// 	utils.Success(c, http.StatusOK, "Subscription retrieved successfully", sub)
// }

// func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	userID := c.GetString("user_id")
// 	var req models.UpdateSubscriptionRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		h.log.Debug("invalid request body", logger.ErrorField(err))
// 		utils.Error(c, http.StatusBadRequest, "Invalid request body", models.ErrInvalidInput)
// 		return
// 	}

// 	sub, err := h.subService.UpdateSubscription(ctx, userID, &req)
// 	if err != nil {
// 		switch err {
// 		case models.ErrInvalidInput, models.ErrInvalidPlan:
// 			utils.Error(c, http.StatusBadRequest, err.Error(), err)
// 		case models.ErrSubscriptionNotActive:
// 			utils.Error(c, http.StatusNotFound, err.Error(), err)
// 		case models.ErrPaymentFailed:
// 			utils.Error(c, http.StatusPaymentRequired, err.Error(), err)
// 		default:
// 			h.log.Error("failed to update subscription", logger.ErrorField(err))
// 			utils.Error(c, http.StatusInternalServerError, "Failed to update subscription", err)
// 		}
// 		return
// 	}

// 	utils.Success(c, http.StatusOK, "Subscription updated successfully", sub)
// }

// func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	userID := c.GetString("user_id")
// 	var req models.CancelSubscriptionRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		h.log.Debug("invalid request body", logger.ErrorField(err))
// 		utils.Error(c, http.StatusBadRequest, "Invalid request body", models.ErrInvalidInput)
// 		return
// 	}

// 	if err := h.subService.CancelSubscription(ctx, userID, &req); err != nil {
// 		switch err {
// 		case models.ErrSubscriptionNotActive:
// 			utils.Error(c, http.StatusNotFound, err.Error(), err)
// 		default:
// 			h.log.Error("failed to cancel subscription", logger.ErrorField(err))
// 			utils.Error(c, http.StatusInternalServerError, "Failed to cancel subscription", err)
// 		}
// 		return
// 	}

// 	c.Status(http.StatusNoContent)
// }

// func (h *SubscriptionHandler) GetPlans(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	plans, err := h.subService.GetSubscriptionPlans(ctx)
// 	if err != nil {
// 		h.log.Error("failed to get subscription plans", logger.ErrorField(err))
// 		utils.Error(c, http.StatusInternalServerError, "Failed to get subscription plans", err)
// 		return
// 	}

// 	utils.Success(c, http.StatusOK, "Subscription plans retrieved successfully", plans)
// }

// func (h *SubscriptionHandler) GetPaymentHistory(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	userID := c.GetString("user_id")

// 	payments, err := h.subService.GetPaymentHistory(ctx, userID)
// 	if err != nil {
// 		h.log.Error("failed to get payment history", logger.ErrorField(err))
// 		utils.Error(c, http.StatusInternalServerError, "Failed to get payment history", err)
// 		return
// 	}

// 	utils.Success(c, http.StatusOK, "Payment history retrieved successfully", payments)
// }


package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"github.com/imraushankr/bervity/server/src/internal/utils"
)

type SubscriptionHandler struct {
	subService interfaces.SubscriptionService
	log        logger.Logger
}

func NewSubscriptionHandler(subService interfaces.SubscriptionService, log logger.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		subService: subService,
		log:        log,
	}
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")
	var req models.CreateSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Debug("invalid request body", logger.ErrorField(err))
		utils.Error(c, http.StatusBadRequest, "Invalid request body", models.ErrInvalidInput)
		return
	}

	sub, err := h.subService.CreateSubscription(ctx, userID, &req)
	if err != nil {
		switch err {
		case models.ErrInvalidInput, models.ErrInvalidPlan:
			utils.Error(c, http.StatusBadRequest, err.Error(), err)
		case models.ErrActiveSubscriptionExists:
			utils.Error(c, http.StatusConflict, err.Error(), err)
		case models.ErrPaymentFailed:
			utils.Error(c, http.StatusPaymentRequired, err.Error(), err)
		default:
			h.log.Error("failed to create subscription", logger.ErrorField(err))
			utils.Error(c, http.StatusInternalServerError, "Failed to create subscription", err)
		}
		return
	}

	utils.Success(c, http.StatusCreated, "Subscription created successfully", sub)
}

func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	sub, err := h.subService.GetUserSubscription(ctx, userID)
	if err != nil {
		if err == models.ErrSubscriptionNotActive {
			utils.Error(c, http.StatusNotFound, err.Error(), err)
			return
		}
		h.log.Error("failed to get subscription", logger.ErrorField(err))
		utils.Error(c, http.StatusInternalServerError, "Failed to get subscription", err)
		return
	}

	utils.Success(c, http.StatusOK, "Subscription retrieved successfully", sub)
}

func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")
	var req models.UpdateSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Debug("invalid request body", logger.ErrorField(err))
		utils.Error(c, http.StatusBadRequest, "Invalid request body", models.ErrInvalidInput)
		return
	}

	sub, err := h.subService.UpdateSubscription(ctx, userID, &req)
	if err != nil {
		switch err {
		case models.ErrInvalidInput, models.ErrInvalidPlan:
			utils.Error(c, http.StatusBadRequest, err.Error(), err)
		case models.ErrSubscriptionNotActive:
			utils.Error(c, http.StatusNotFound, err.Error(), err)
		case models.ErrPaymentFailed:
			utils.Error(c, http.StatusPaymentRequired, err.Error(), err)
		default:
			h.log.Error("failed to update subscription", logger.ErrorField(err))
			utils.Error(c, http.StatusInternalServerError, "Failed to update subscription", err)
		}
		return
	}

	utils.Success(c, http.StatusOK, "Subscription updated successfully", sub)
}

func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")
	var req models.CancelSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Debug("invalid request body", logger.ErrorField(err))
		utils.Error(c, http.StatusBadRequest, "Invalid request body", models.ErrInvalidInput)
		return
	}

	if err := h.subService.CancelSubscription(ctx, userID, &req); err != nil {
		switch err {
		case models.ErrSubscriptionNotActive:
			utils.Error(c, http.StatusNotFound, err.Error(), err)
		default:
			h.log.Error("failed to cancel subscription", logger.ErrorField(err))
			utils.Error(c, http.StatusInternalServerError, "Failed to cancel subscription", err)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SubscriptionHandler) GetPlans(c *gin.Context) {
	ctx := c.Request.Context()

	plans, err := h.subService.GetSubscriptionPlans(ctx)
	if err != nil {
		h.log.Error("failed to get subscription plans", logger.ErrorField(err))
		utils.Error(c, http.StatusInternalServerError, "Failed to get subscription plans", err)
		return
	}

	utils.Success(c, http.StatusOK, "Subscription plans retrieved successfully", plans)
}

func (h *SubscriptionHandler) GetPaymentHistory(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	payments, err := h.subService.GetPaymentHistory(ctx, userID)
	if err != nil {
		h.log.Error("failed to get payment history", logger.ErrorField(err))
		utils.Error(c, http.StatusInternalServerError, "Failed to get payment history", err)
		return
	}

	utils.Success(c, http.StatusOK, "Payment history retrieved successfully", payments)
}