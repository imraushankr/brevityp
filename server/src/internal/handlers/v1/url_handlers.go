// package v1

// import (
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/imraushankr/bervity/server/src/internal/models"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// 	"github.com/imraushankr/bervity/server/src/internal/utils"
// )

// type URLHandler struct {
// 	urlService interfaces.URLService
// 	log        logger.Logger
// }

// func NewURLHandler(urlService interfaces.URLService, log logger.Logger) *URLHandler {
// 	return &URLHandler{
// 		urlService: urlService,
// 		log:        log,
// 	}
// }

// func (h *URLHandler) CreateURL(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	var req models.CreateURLRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		h.log.Debug("invalid request body", logger.ErrorField(err))
// 		utils.Error(c, http.StatusBadRequest, "Invalid request body", models.ErrInvalidInput)
// 		return
// 	}

// 	userID := c.GetString("user_id")

// 	resp, err := h.urlService.CreateURL(ctx, &req, userID)
// 	if err != nil {
// 		switch err {
// 		case models.ErrInvalidInput, models.ErrShortCodeTaken:
// 			utils.Error(c, http.StatusBadRequest, err.Error(), err)
// 		case models.ErrInsufficientCredits:
// 			utils.Error(c, http.StatusPaymentRequired, err.Error(), err)
// 		default:
// 			h.log.Error("failed to create URL", logger.ErrorField(err))
// 			utils.Error(c, http.StatusInternalServerError, "Failed to create URL", err)
// 		}
// 		return
// 	}

// 	utils.Success(c, http.StatusCreated, "URL created successfully", resp)
// }

package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
	"github.com/imraushankr/bervity/server/src/internal/utils"
)

type URLHandler struct {
	urlService interfaces.URLService
	log        logger.Logger
}

func NewURLHandler(urlService interfaces.URLService, log logger.Logger) *URLHandler {
	return &URLHandler{
		urlService: urlService,
		log:        log,
	}
}

func (h *URLHandler) CreateURL(c *gin.Context) {
	ctx := c.Request.Context()
	var req models.CreateURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Debug("invalid request body", logger.ErrorField(err))
		utils.Error(c, http.StatusBadRequest, "Invalid request body", models.ErrInvalidInput)
		return
	}

	// Get user ID if authenticated
	userID := c.GetString("user_id")
	ip := c.ClientIP()

	resp, err := h.urlService.CreateURL(ctx, &req, userID, ip)
	if err != nil {
		switch err {
		case models.ErrInvalidInput, models.ErrShortCodeTaken:
			utils.Error(c, http.StatusBadRequest, err.Error(), err)
		case models.ErrInsufficientCredits:
			utils.Error(c, http.StatusPaymentRequired, err.Error(), err)
		default:
			h.log.Error("failed to create URL", logger.ErrorField(err))
			utils.Error(c, http.StatusInternalServerError, "Failed to create URL", err)
		}
		return
	}

	utils.Success(c, http.StatusCreated, "URL created successfully", resp)
}


func (h *URLHandler) GetURL(c *gin.Context) {
	ctx := c.Request.Context()
	shortCode := c.Param("code")

	url, err := h.urlService.GetURL(ctx, shortCode)
	if err != nil {
		if err == models.ErrURLNotFound {
			utils.Error(c, http.StatusNotFound, err.Error(), err)
			return
		}
		h.log.Error("failed to get URL", logger.ErrorField(err))
		utils.Error(c, http.StatusInternalServerError, "Failed to get URL", err)
		return
	}

	utils.Success(c, http.StatusOK, "URL retrieved successfully", url.ToResponse(""))
}

func (h *URLHandler) GetUserURLs(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	// Convert limit and offset from string to int
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid limit parameter", err)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid offset parameter", err)
		return
	}

	urls, err := h.urlService.GetUserURLs(ctx, userID, limit, offset)
	if err != nil {
		h.log.Error("failed to get user URLs", logger.ErrorField(err))
		utils.Error(c, http.StatusInternalServerError, "Failed to get URLs", err)
		return
	}

	utils.Success(c, http.StatusOK, "User URLs retrieved successfully", urls)
}

func (h *URLHandler) UpdateURL(c *gin.Context) {
	ctx := c.Request.Context()
	// userID := c.GetString("user_id")
	var url models.URL

	if err := c.ShouldBindJSON(&url); err != nil {
		h.log.Debug("invalid request body", logger.ErrorField(err))
		utils.Error(c, http.StatusBadRequest, "Invalid request body", models.ErrInvalidInput)
		return
	}

	resp, err := h.urlService.UpdateURL(ctx, &url)
	if err != nil {
		switch err {
		case models.ErrURLNotFound:
			utils.Error(c, http.StatusNotFound, err.Error(), err)
		case models.ErrForbidden:
			utils.Error(c, http.StatusForbidden, err.Error(), err)
		default:
			h.log.Error("failed to update URL", logger.ErrorField(err))
			utils.Error(c, http.StatusInternalServerError, "Failed to update URL", err)
		}
		return
	}

	utils.Success(c, http.StatusOK, "URL updated successfully", resp)
}

func (h *URLHandler) DeleteURL(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")
	urlID := c.Param("id")

	if err := h.urlService.DeleteURL(ctx, urlID, userID); err != nil {
		switch err {
		case models.ErrURLNotFound:
			utils.Error(c, http.StatusNotFound, err.Error(), err)
		case models.ErrForbidden:
			utils.Error(c, http.StatusForbidden, err.Error(), err)
		default:
			h.log.Error("failed to delete URL", logger.ErrorField(err))
			utils.Error(c, http.StatusInternalServerError, "Failed to delete URL", err)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *URLHandler) Redirect(c *gin.Context) {
	ctx := c.Request.Context()
	shortCode := c.Param("code")

	clickData := &models.URLClick{
		IPAddress: c.ClientIP(),
		Referrer:  c.Request.Referer(),
		UserAgent: c.Request.UserAgent(),
	}

	originalURL, err := h.urlService.RedirectURL(ctx, shortCode, clickData)
	if err != nil {
		if err == models.ErrURLNotFound {
			utils.Error(c, http.StatusNotFound, err.Error(), err)
			return
		}
		h.log.Error("failed to redirect", logger.ErrorField(err))
		utils.Error(c, http.StatusInternalServerError, "Failed to redirect", err)
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}

func (h *URLHandler) GetAnalytics(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")
	urlID := c.Param("id")

	fromStr := c.DefaultQuery("from", time.Now().AddDate(0, 0, -7).Format(time.RFC3339))
	toStr := c.DefaultQuery("to", time.Now().Format(time.RFC3339))

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid from date", err)
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid to date", err)
		return
	}

	analytics, err := h.urlService.GetURLAnalytics(ctx, urlID, userID, from, to)
	if err != nil {
		switch err {
		case models.ErrURLNotFound:
			utils.Error(c, http.StatusNotFound, err.Error(), err)
		case models.ErrForbidden:
			utils.Error(c, http.StatusForbidden, err.Error(), err)
		default:
			h.log.Error("failed to get analytics", logger.ErrorField(err))
			utils.Error(c, http.StatusInternalServerError, "Failed to get analytics", err)
		}
		return
	}

	utils.Success(c, http.StatusOK, "Analytics retrieved successfully", analytics)
}