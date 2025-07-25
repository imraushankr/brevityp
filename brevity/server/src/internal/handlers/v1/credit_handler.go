// package v1

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/imraushankr/bervity/server/src/internal/models"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// 	"github.com/imraushankr/bervity/server/src/internal/utils"
// )

// type CreditHandler struct {
// 	creditService interfaces.CreditService
// 	log           logger.Logger
// }

// func NewCreditHandler(creditService interfaces.CreditService, log logger.Logger) *CreditHandler {
// 	return &CreditHandler{
// 		creditService: creditService,
// 		log:           log,
// 	}
// }

// func (h *CreditHandler) GetBalance(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	userID := c.GetString("user_id")

// 	balance, err := h.creditService.GetCreditBalance(ctx, userID)
// 	if err != nil {
// 		h.log.Error("failed to get credit balance", logger.ErrorField(err))
// 		utils.Error(c, http.StatusInternalServerError, "Failed to get credit balance", err)
// 		return
// 	}

// 	utils.Success(c, http.StatusOK, "Credit balance retrieved successfully", balance)
// }

// func (h *CreditHandler) ApplyPromoCode(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	userID := c.GetString("user_id")
// 	var req models.ApplyPromoCodeRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		h.log.Debug("invalid request body", logger.ErrorField(err))
// 		utils.Error(c, http.StatusBadRequest, "Invalid request body", models.ErrInvalidInput)
// 		return
// 	}

// 	credit, err := h.creditService.ApplyPromoCode(ctx, userID, req.Code)
// 	if err != nil {
// 		switch err {
// 		case models.ErrInvalidInput, models.ErrPromoCodeInvalid, models.ErrPromoCodeAlreadyUsed:
// 			utils.Error(c, http.StatusBadRequest, err.Error(), err)
// 		default:
// 			h.log.Error("failed to apply promo code", logger.ErrorField(err))
// 			utils.Error(c, http.StatusInternalServerError, "Failed to apply promo code", err)
// 		}
// 		return
// 	}

// 	utils.Success(c, http.StatusOK, "Promo code applied successfully", credit)
// }

// func (h *CreditHandler) GetUsage(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	userID := c.GetString("user_id")

// 	usage, err := h.creditService.GetCreditUsage(ctx, userID)
// 	if err != nil {
// 		h.log.Error("failed to get credit usage", logger.ErrorField(err))
// 		utils.Error(c, http.StatusInternalServerError, "Failed to get credit usage", err)
// 		return
// 	}

// 	utils.Success(c, http.StatusOK, "Credit usage retrieved successfully", usage)
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

type CreditHandler struct {
	creditService interfaces.CreditService
	log           logger.Logger
}

func NewCreditHandler(creditService interfaces.CreditService, log logger.Logger) *CreditHandler {
	return &CreditHandler{
		creditService: creditService,
		log:           log,
	}
}

func (h *CreditHandler) GetBalance(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	balance, err := h.creditService.GetCreditBalance(ctx, userID)
	if err != nil {
		h.log.Error("failed to get credit balance", logger.ErrorField(err))
		utils.Error(c, http.StatusInternalServerError, "Failed to get credit balance", err)
		return
	}

	utils.Success(c, http.StatusOK, "Credit balance retrieved successfully", balance)
}

func (h *CreditHandler) ApplyPromoCode(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")
	var req models.ApplyPromoCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Debug("invalid request body", logger.ErrorField(err))
		utils.Error(c, http.StatusBadRequest, "Invalid request body", models.ErrInvalidInput)
		return
	}

	credit, err := h.creditService.ApplyPromoCode(ctx, userID, req.Code)
	if err != nil {
		switch err {
		case models.ErrInvalidInput, models.ErrPromoCodeInvalid, models.ErrPromoCodeAlreadyUsed:
			utils.Error(c, http.StatusBadRequest, err.Error(), err)
		default:
			h.log.Error("failed to apply promo code", logger.ErrorField(err))
			utils.Error(c, http.StatusInternalServerError, "Failed to apply promo code", err)
		}
		return
	}

	utils.Success(c, http.StatusOK, "Promo code applied successfully", credit)
}

func (h *CreditHandler) GetUsage(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	usage, err := h.creditService.GetCreditUsage(ctx, userID)
	if err != nil {
		h.log.Error("failed to get credit usage", logger.ErrorField(err))
		utils.Error(c, http.StatusInternalServerError, "Failed to get credit usage", err)
		return
	}

	utils.Success(c, http.StatusOK, "Credit usage retrieved successfully", usage)
}