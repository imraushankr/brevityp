// // package services

// // import (
// // 	"context"
// // 	"errors"
// // 	"net/url"
// // 	"strings"
// // 	"time"

// // 	"github.com/google/uuid"
// // 	"github.com/imraushankr/bervity/server/src/internal/models"
// // 	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
// // 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// // )

// // type urlService struct {
// // 	urlRepo       interfaces.URLRepository
// // 	creditRepo    interfaces.CreditRepository
// // 	analyticsRepo interfaces.AnalyticsRepository
// // 	logger        logger.Logger
// // 	baseURL       string
// // 	freeURLLimit  int
// // }

// // func NewURLService(
// // 	urlRepo interfaces.URLRepository,
// // 	creditRepo interfaces.CreditRepository,
// // 	analyticsRepo interfaces.AnalyticsRepository,
// // 	logger logger.Logger,
// // 	baseURL string,
// // 	freeURLLimit int,
// // ) interfaces.URLService {
// // 	return &urlService{
// // 		urlRepo:       urlRepo,
// // 		creditRepo:    creditRepo,
// // 		analyticsRepo: analyticsRepo,
// // 		logger:        logger,
// // 		baseURL:       baseURL,
// // 		freeURLLimit:  freeURLLimit,
// // 	}
// // }

// // func (s *urlService) CreateURL(ctx context.Context, req *models.CreateURLRequest, userID string) (*models.URLResponse, error) {
// // 	// Validate original URL
// // 	if _, err := url.ParseRequestURI(req.OriginalURL); err != nil {
// // 		s.logger.Debug("invalid URL format",
// // 			logger.String("url", req.OriginalURL),
// // 			logger.ErrorField(err))
// // 		return nil, models.ErrInvalidInput
// // 	}

// // 	// Check if user has enough credits if authenticated
// // 	if userID != "" {
// // 		balance, err := s.creditRepo.GetUserCreditBalance(ctx, userID)
// // 		if err != nil {
// // 			s.logger.Error("failed to get user credit balance",
// // 				logger.ErrorField(err),
// // 				logger.String("userID", userID))
// // 			return nil, err
// // 		}

// // 		if !balance.CanCreate {
// // 			s.logger.Warn("insufficient credits",
// // 				logger.String("userID", userID),
// // 				logger.Any("balance", balance))
// // 			return nil, models.ErrInsufficientCredits
// // 		}
// // 	}

// // 	shortCode := req.CustomCode
// // 	if shortCode == "" {
// // 		// Generate a random short code if not provided
// // 		shortCode = generateShortCode(6)
// // 	} else {
// // 		// Validate custom code if provided
// // 		if len(shortCode) < 3 || len(shortCode) > 10 || !isAlphanumeric(shortCode) {
// // 			s.logger.Debug("invalid custom code format",
// // 				logger.String("code", shortCode))
// // 			return nil, models.ErrInvalidInput
// // 		}
// // 	}

// // 	// Check if short code is already taken
// // 	_, err := s.urlRepo.GetByShortCode(ctx, shortCode)
// // 	if err == nil {
// // 		s.logger.Debug("short code already exists",
// // 			logger.String("code", shortCode))
// // 		return nil, models.ErrShortCodeTaken
// // 	} else if !errors.Is(err, models.ErrURLNotFound) {
// // 		s.logger.Error("failed to check short code availability",
// // 			logger.ErrorField(err),
// // 			logger.String("code", shortCode))
// // 		return nil, err
// // 	}

// // 	// Create the URL
// // 	newURL := &models.URL{
// // 		OriginalURL: req.OriginalURL,
// // 		ShortCode:   shortCode,
// // 		UserID:      userID,
// // 		Title:       req.Title,
// // 		Description: req.Description,
// // 		ExpiresAt:   req.ExpiresAt,
// // 		IsActive:    true,
// // 	}

// // 	if err := s.urlRepo.Create(ctx, newURL); err != nil {
// // 		s.logger.Error("failed to create URL",
// // 			logger.ErrorField(err),
// // 			logger.Any("url", newURL))
// // 		return nil, err
// // 	}

// // 	// Deduct credit if user is authenticated
// // 	if userID != "" {
// // 		err = s.creditRepo.UseCredits(ctx, &models.CreditUsage{
// // 			UserID:    userID,
// // 			URLID:     newURL.ID,
// // 			Amount:    1,
// // 			Operation: "url_creation",
// // 		})
// // 		if err != nil {
// // 			s.logger.Error("failed to deduct credits",
// // 				logger.ErrorField(err),
// // 				logger.String("userID", userID),
// // 				logger.String("urlID", newURL.ID))
// // 			// Attempt to delete the URL if credit deduction fails
// // 			if delErr := s.urlRepo.Delete(ctx, newURL.ID); delErr != nil {
// // 				s.logger.Error("failed to rollback URL creation",
// // 					logger.ErrorField(delErr),
// // 					logger.String("urlID", newURL.ID))
// // 			}
// // 			return nil, err
// // 		}
// // 	}

// // 	s.logger.Info("URL created successfully",
// // 		logger.String("urlID", newURL.ID),
// // 		logger.String("shortCode", newURL.ShortCode))

// // 	return newURL.ToResponse(s.baseURL), nil
// // }

// package services

// import (
// 	"context"
// 	"errors"
// 	"net/url"
// 	"strings"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/imraushankr/bervity/server/src/internal/models"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
// 	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
// )

// type urlService struct {
// 	urlRepo       interfaces.URLRepository
// 	creditRepo    interfaces.CreditRepository
// 	analyticsRepo interfaces.AnalyticsRepository
// 	logger        logger.Logger
// 	baseURL       string
// 	anonURLLimit  int // 5 for anonymous users
// 	authURLLimit  int // 15 for authenticated users
// }

// func NewURLService(
// 	urlRepo interfaces.URLRepository,
// 	creditRepo interfaces.CreditRepository,
// 	analyticsRepo interfaces.AnalyticsRepository,
// 	logger logger.Logger,
// 	baseURL string,
// 	anonURLLimit int,
// 	authURLLimit int,
// ) interfaces.URLService {
// 	return &urlService{
// 		urlRepo:       urlRepo,
// 		creditRepo:    creditRepo,
// 		analyticsRepo: analyticsRepo,
// 		logger:        logger,
// 		baseURL:       baseURL,
// 		anonURLLimit:  anonURLLimit,
// 		authURLLimit:  authURLLimit,
// 	}
// }

// func (s *urlService) CreateURL(ctx context.Context, req *models.CreateURLRequest, userID string, ip string) (*models.URLResponse, error) {
// 	// Validate original URL
// 	if _, err := url.ParseRequestURI(req.OriginalURL); err != nil {
// 		s.logger.Debug("invalid URL format",
// 			logger.String("url", req.OriginalURL),
// 			logger.ErrorField(err))
// 		return nil, models.ErrInvalidInput
// 	}

// 	shortCode := req.CustomCode
// 	if shortCode == "" {
// 		// Generate a random short code if not provided
// 		shortCode = generateShortCode(6)
// 	} else {
// 		// Validate custom code if provided
// 		if len(shortCode) < 3 || len(shortCode) > 10 || !isAlphanumeric(shortCode) {
// 			s.logger.Debug("invalid custom code format",
// 				logger.String("code", shortCode))
// 			return nil, models.ErrInvalidInput
// 		}
// 	}

// 	// Check if short code is already taken
// 	_, err := s.urlRepo.GetByShortCode(ctx, shortCode)
// 	if err == nil {
// 		s.logger.Debug("short code already exists",
// 			logger.String("code", shortCode))
// 		return nil, models.ErrShortCodeTaken
// 	} else if !errors.Is(err, models.ErrURLNotFound) {
// 		s.logger.Error("failed to check short code availability",
// 			logger.ErrorField(err),
// 			logger.String("code", shortCode))
// 		return nil, err
// 	}

// 	// Check URL creation limits
// 	if userID == "" {
// 		// Anonymous user - check if they've reached the limit
// 		// Note: In a real app, you might want to track this by IP or session
// 		// For simplicity, we're not tracking anonymous users here
// 		// You would need to implement IP-based tracking in the repository
// 		// For now, we'll just allow anonymous users to create URLs without tracking
// 		// In production, you should implement proper tracking
// 	} else {
// 		// Authenticated user - check free limit or credits
// 		freeCount, err := s.creditRepo.GetFreeURLCount(ctx, userID)
// 		if err != nil {
// 			s.logger.Error("failed to get free URL count",
// 				logger.ErrorField(err),
// 				logger.String("userID", userID))
// 			return nil, err
// 		}

// 		if freeCount < s.authURLLimit {
// 			// User is within free limit, proceed without deducting credits
// 		} else {
// 			// User exceeded free limit, check paid credits
// 			balance, err := s.creditRepo.GetUserCreditBalance(ctx, userID)
// 			if err != nil {
// 				s.logger.Error("failed to get user credit balance",
// 					logger.ErrorField(err),
// 					logger.String("userID", userID))
// 				return nil, err
// 			}

// 			if !balance.CanCreate {
// 				s.logger.Warn("insufficient credits",
// 					logger.String("userID", userID),
// 					logger.Any("balance", balance))
// 				return nil, models.ErrInsufficientCredits
// 			}
// 		}
// 	}

// 	var userIDPtr *string
// 	if userID != "" {
// 		userIDPtr = &userID
// 	}

// 	newURL := &models.URL{
// 		OriginalURL: req.OriginalURL,
// 		ShortCode:   shortCode,
// 		UserID:      userIDPtr, // Use pointer here
// 		CreatedByIP: ip,
// 		Title:       req.Title,
// 		Description: req.Description,
// 		ExpiresAt:   req.ExpiresAt,
// 		IsActive:    true,
// 	}

// 	if err := s.urlRepo.Create(ctx, newURL); err != nil {
// 		s.logger.Error("failed to create URL",
// 			logger.ErrorField(err),
// 			logger.Any("url", newURL))
// 		return nil, err
// 	}

// 	// Record the URL creation
// 	if userID != "" {
// 		freeCount, _ := s.creditRepo.GetFreeURLCount(ctx, userID)
// 		if freeCount < s.authURLLimit {
// 			// Record as free URL creation
// 			if err := s.creditRepo.RecordFreeURLCreation(ctx, userID, newURL.ID); err != nil {
// 				s.logger.Error("failed to record free URL creation",
// 					logger.ErrorField(err),
// 					logger.String("userID", userID),
// 					logger.String("urlID", newURL.ID))
// 				// Attempt to delete the URL if recording fails
// 				if delErr := s.urlRepo.Delete(ctx, newURL.ID); delErr != nil {
// 					s.logger.Error("failed to rollback URL creation",
// 						logger.ErrorField(delErr),
// 						logger.String("urlID", newURL.ID))
// 				}
// 				return nil, err
// 			}
// 		} else {
// 			// Deduct credit
// 			err = s.creditRepo.UseCredits(ctx, &models.CreditUsage{
// 				UserID:    userID,
// 				URLID:     newURL.ID,
// 				Amount:    1,
// 				Operation: "url_creation",
// 			})
// 			if err != nil {
// 				s.logger.Error("failed to deduct credits",
// 					logger.ErrorField(err),
// 					logger.String("userID", userID),
// 					logger.String("urlID", newURL.ID))
// 				// Attempt to delete the URL if credit deduction fails
// 				if delErr := s.urlRepo.Delete(ctx, newURL.ID); delErr != nil {
// 					s.logger.Error("failed to rollback URL creation",
// 						logger.ErrorField(delErr),
// 						logger.String("urlID", newURL.ID))
// 				}
// 				return nil, err
// 			}
// 		}
// 	}

// 	s.logger.Info("URL created successfully",
// 		logger.String("urlID", newURL.ID),
// 		logger.String("shortCode", newURL.ShortCode))

// 	return newURL.ToResponse(s.baseURL), nil
// }

// func (s *urlService) GetURL(ctx context.Context, shortCode string) (*models.URL, error) {
// 	url, err := s.urlRepo.GetByShortCode(ctx, shortCode)
// 	if err != nil {
// 		s.logger.Error("failed to get URL",
// 			logger.ErrorField(err),
// 			logger.String("shortCode", shortCode))
// 		return nil, err
// 	}
// 	return url, nil
// }

// func (s *urlService) GetUserURLs(ctx context.Context, userID string, limit, offset int) ([]*models.URLResponse, error) {
// 	urls, err := s.urlRepo.GetByUser(ctx, userID, limit, offset)
// 	if err != nil {
// 		s.logger.Error("failed to get user URLs",
// 			logger.ErrorField(err),
// 			logger.String("userID", userID))
// 		return nil, err
// 	}

// 	responses := make([]*models.URLResponse, len(urls))
// 	for i, u := range urls {
// 		responses[i] = u.ToResponse(s.baseURL)
// 	}

// 	return responses, nil
// }

// func (s *urlService) UpdateURL(ctx context.Context, url *models.URL) (*models.URLResponse, error) {
// 	// Verify the URL exists and belongs to the user
// 	existingURL, err := s.urlRepo.GetByID(ctx, url.ID)
// 	if err != nil {
// 		s.logger.Error("failed to get existing URL",
// 			logger.ErrorField(err),
// 			logger.String("urlID", url.ID))
// 		return nil, err
// 	}

// 	// Only allow certain fields to be updated
// 	existingURL.Title = url.Title
// 	existingURL.Description = url.Description
// 	existingURL.ExpiresAt = url.ExpiresAt
// 	existingURL.IsActive = url.IsActive

// 	if err := s.urlRepo.Update(ctx, existingURL); err != nil {
// 		s.logger.Error("failed to update URL",
// 			logger.ErrorField(err),
// 			logger.Any("url", existingURL))
// 		return nil, err
// 	}

// 	s.logger.Info("URL updated successfully",
// 		logger.String("urlID", existingURL.ID))

// 	return existingURL.ToResponse(s.baseURL), nil
// }

// func (s *urlService) DeleteURL(ctx context.Context, id, userID string) error {
// 	// Verify URL belongs to user
// 	url, err := s.urlRepo.GetByID(ctx, id)
// 	if err != nil {
// 		s.logger.Error("failed to get URL for deletion",
// 			logger.ErrorField(err),
// 			logger.String("urlID", id))
// 		return err
// 	}

// 	if url.UserID != userID {
// 		s.logger.Warn("unauthorized URL deletion attempt",
// 			logger.String("requestingUserID", userID),
// 			logger.String("urlOwnerID", url.UserID),
// 			logger.String("urlID", id))
// 		return models.ErrForbidden
// 	}

// 	if err := s.urlRepo.Delete(ctx, id); err != nil {
// 		s.logger.Error("failed to delete URL",
// 			logger.ErrorField(err),
// 			logger.String("urlID", id))
// 		return err
// 	}

// 	s.logger.Info("URL deleted successfully",
// 		logger.String("urlID", id))

// 	return nil
// }

// func (s *urlService) RedirectURL(ctx context.Context, shortCode string, clickData *models.URLClick) (string, error) {
// 	url, err := s.urlRepo.GetByShortCode(ctx, shortCode)
// 	if err != nil {
// 		s.logger.Error("failed to get URL for redirect",
// 			logger.ErrorField(err),
// 			logger.String("shortCode", shortCode))
// 		return "", err
// 	}

// 	// Record the click if analytics data is provided
// 	if clickData != nil {
// 		clickData.URLID = url.ID
// 		clickData.ID = uuid.New().String()
// 		clickData.CreatedAt = time.Now()

// 		if err := s.urlRepo.RecordClick(ctx, clickData); err != nil {
// 			s.logger.Error("failed to record URL click",
// 				logger.ErrorField(err),
// 				logger.Any("clickData", clickData))
// 			// Don't fail the redirect if click recording fails
// 		}
// 	}

// 	// Increment click count
// 	if err := s.urlRepo.IncrementClicks(ctx, url.ID); err != nil {
// 		s.logger.Error("failed to increment URL clicks",
// 			logger.ErrorField(err),
// 			logger.String("urlID", url.ID))
// 		// Don't fail the redirect if click count fails
// 	}

// 	return url.OriginalURL, nil
// }

// func (s *urlService) GetURLAnalytics(ctx context.Context, urlID, userID string, from, to time.Time) ([]*models.URLClick, error) {
// 	// Verify URL belongs to user
// 	url, err := s.urlRepo.GetByID(ctx, urlID)
// 	if err != nil {
// 		s.logger.Error("failed to get URL for analytics",
// 			logger.ErrorField(err),
// 			logger.String("urlID", urlID))
// 		return nil, err
// 	}

// 	if url.UserID != userID {
// 		s.logger.Warn("unauthorized analytics access attempt",
// 			logger.String("requestingUserID", userID),
// 			logger.String("urlOwnerID", url.UserID),
// 			logger.String("urlID", urlID))
// 		return nil, models.ErrForbidden
// 	}

// 	clicks, err := s.urlRepo.GetClicksAnalytics(ctx, urlID, from, to)
// 	if err != nil {
// 		s.logger.Error("failed to get URL analytics",
// 			logger.ErrorField(err),
// 			logger.String("urlID", urlID),
// 			logger.Time("from", from),
// 			logger.Time("to", to))
// 		return nil, err
// 	}

// 	return clicks, nil
// }

// // Helper functions
// func generateShortCode(length int) string {
// 	// This is a simplified version - in production, use a more robust generator
// 	uuidStr := uuid.New().String()
// 	clean := strings.ReplaceAll(uuidStr, "-", "")
// 	if length > len(clean) {
// 		length = len(clean)
// 	}
// 	return clean[:length]
// }

// func isAlphanumeric(s string) bool {
// 	for _, r := range s {
// 		if !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z') && !(r >= '0' && r <= '9') {
// 			return false
// 		}
// 	}
// 	return true
// }


package services

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/imraushankr/bervity/server/src/internal/models"
	"github.com/imraushankr/bervity/server/src/internal/pkg/interfaces"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

type urlService struct {
	urlRepo       interfaces.URLRepository
	creditRepo    interfaces.CreditRepository
	analyticsRepo interfaces.AnalyticsRepository
	logger        logger.Logger
	baseURL       string
	anonURLLimit  int // 5 for anonymous users
	authURLLimit  int // 15 for authenticated users
}

func NewURLService(
	urlRepo interfaces.URLRepository,
	creditRepo interfaces.CreditRepository,
	analyticsRepo interfaces.AnalyticsRepository,
	logger logger.Logger,
	baseURL string,
	anonURLLimit int,
	authURLLimit int,
) interfaces.URLService {
	return &urlService{
		urlRepo:       urlRepo,
		creditRepo:    creditRepo,
		analyticsRepo: analyticsRepo,
		logger:        logger,
		baseURL:       baseURL,
		anonURLLimit:  anonURLLimit,
		authURLLimit:  authURLLimit,
	}
}

func (s *urlService) CreateURL(ctx context.Context, req *models.CreateURLRequest, userID string, ip string) (*models.URLResponse, error) {
	// Validate original URL
	if _, err := url.ParseRequestURI(req.OriginalURL); err != nil {
		s.logger.Debug("invalid URL format",
			logger.String("url", req.OriginalURL),
			logger.ErrorField(err))
		return nil, models.ErrInvalidInput
	}

	shortCode := req.CustomCode
	if shortCode == "" {
		shortCode = generateShortCode(6)
	} else {
		if len(shortCode) < 3 || len(shortCode) > 10 || !isAlphanumeric(shortCode) {
			s.logger.Debug("invalid custom code format",
				logger.String("code", shortCode))
			return nil, models.ErrInvalidInput
		}
	}

	// Check if short code is already taken
	_, err := s.urlRepo.GetByShortCode(ctx, shortCode)
	if err == nil {
		s.logger.Debug("short code already exists",
			logger.String("code", shortCode))
		return nil, models.ErrShortCodeTaken
	} else if !errors.Is(err, models.ErrURLNotFound) {
		s.logger.Error("failed to check short code availability",
			logger.ErrorField(err),
			logger.String("code", shortCode))
		return nil, err
	}

	// Check URL creation limits
	if userID != "" {
		freeCount, err := s.creditRepo.GetFreeURLCount(ctx, userID)
		if err != nil {
			s.logger.Error("failed to get free URL count",
				logger.ErrorField(err),
				logger.String("userID", userID))
			return nil, err
		}

		if freeCount >= s.authURLLimit {
			balance, err := s.creditRepo.GetUserCreditBalance(ctx, userID)
			if err != nil {
				s.logger.Error("failed to get user credit balance",
					logger.ErrorField(err),
					logger.String("userID", userID))
				return nil, err
			}

			if !balance.CanCreate {
				s.logger.Warn("insufficient credits",
					logger.String("userID", userID),
					logger.Any("balance", balance))
				return nil, models.ErrInsufficientCredits
			}
		}
	}

	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}

	newURL := &models.URL{
		OriginalURL: req.OriginalURL,
		ShortCode:   shortCode,
		UserID:      userIDPtr,
		CreatedByIP: ip,
		Title:       req.Title,
		Description: req.Description,
		ExpiresAt:   req.ExpiresAt,
		IsActive:    true,
	}

	if err := s.urlRepo.Create(ctx, newURL); err != nil {
		s.logger.Error("failed to create URL",
			logger.ErrorField(err),
			logger.Any("url", newURL))
		return nil, err
	}

	// Record the URL creation for authenticated users
	if userID != "" {
		freeCount, _ := s.creditRepo.GetFreeURLCount(ctx, userID)
		if freeCount < s.authURLLimit {
			if err := s.creditRepo.RecordFreeURLCreation(ctx, userID, newURL.ID); err != nil {
				s.logger.Error("failed to record free URL creation",
					logger.ErrorField(err),
					logger.String("userID", userID),
					logger.String("urlID", newURL.ID))
				if delErr := s.urlRepo.Delete(ctx, newURL.ID); delErr != nil {
					s.logger.Error("failed to rollback URL creation",
						logger.ErrorField(delErr),
						logger.String("urlID", newURL.ID))
				}
				return nil, err
			}
		} else {
			err = s.creditRepo.UseCredits(ctx, &models.CreditUsage{
				UserID:    userID,
				URLID:     newURL.ID,
				Amount:    1,
				Operation: "url_creation",
			})
			if err != nil {
				s.logger.Error("failed to deduct credits",
					logger.ErrorField(err),
					logger.String("userID", userID),
					logger.String("urlID", newURL.ID))
				if delErr := s.urlRepo.Delete(ctx, newURL.ID); delErr != nil {
					s.logger.Error("failed to rollback URL creation",
						logger.ErrorField(delErr),
						logger.String("urlID", newURL.ID))
				}
				return nil, err
			}
		}
	}

	s.logger.Info("URL created successfully",
		logger.String("urlID", newURL.ID),
		logger.String("shortCode", newURL.ShortCode))

	return newURL.ToResponse(s.baseURL), nil
}

func (s *urlService) GetURL(ctx context.Context, shortCode string) (*models.URL, error) {
	url, err := s.urlRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		s.logger.Error("failed to get URL",
			logger.ErrorField(err),
			logger.String("shortCode", shortCode))
		return nil, err
	}
	return url, nil
}

func (s *urlService) GetUserURLs(ctx context.Context, userID string, limit, offset int) ([]*models.URLResponse, error) {
	urls, err := s.urlRepo.GetByUser(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Error("failed to get user URLs",
			logger.ErrorField(err),
			logger.String("userID", userID))
		return nil, err
	}

	responses := make([]*models.URLResponse, len(urls))
	for i, u := range urls {
		responses[i] = u.ToResponse(s.baseURL)
	}

	return responses, nil
}

func (s *urlService) UpdateURL(ctx context.Context, url *models.URL) (*models.URLResponse, error) {
	existingURL, err := s.urlRepo.GetByID(ctx, url.ID)
	if err != nil {
		s.logger.Error("failed to get existing URL",
			logger.ErrorField(err),
			logger.String("urlID", url.ID))
		return nil, err
	}

	existingURL.Title = url.Title
	existingURL.Description = url.Description
	existingURL.ExpiresAt = url.ExpiresAt
	existingURL.IsActive = url.IsActive

	if err := s.urlRepo.Update(ctx, existingURL); err != nil {
		s.logger.Error("failed to update URL",
			logger.ErrorField(err),
			logger.Any("url", existingURL))
		return nil, err
	}

	s.logger.Info("URL updated successfully",
		logger.String("urlID", existingURL.ID))

	return existingURL.ToResponse(s.baseURL), nil
}

func (s *urlService) DeleteURL(ctx context.Context, id, userID string) error {
	url, err := s.urlRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get URL for deletion",
			logger.ErrorField(err),
			logger.String("urlID", id))
		return err
	}

	if url.UserID != nil && *url.UserID != userID {
		s.logger.Warn("unauthorized URL deletion attempt",
			logger.String("requestingUserID", userID),
			logger.String("urlOwnerID", *url.UserID),
			logger.String("urlID", id))
		return models.ErrForbidden
	}

	if err := s.urlRepo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete URL",
			logger.ErrorField(err),
			logger.String("urlID", id))
		return err
	}

	s.logger.Info("URL deleted successfully",
		logger.String("urlID", id))

	return nil
}

func (s *urlService) RedirectURL(ctx context.Context, shortCode string, clickData *models.URLClick) (string, error) {
	url, err := s.urlRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		s.logger.Error("failed to get URL for redirect",
			logger.ErrorField(err),
			logger.String("shortCode", shortCode))
		return "", err
	}

	if clickData != nil {
		clickData.URLID = url.ID
		clickData.ID = uuid.New().String()
		clickData.CreatedAt = time.Now()

		if err := s.urlRepo.RecordClick(ctx, clickData); err != nil {
			s.logger.Error("failed to record URL click",
				logger.ErrorField(err),
				logger.Any("clickData", clickData))
		}
	}

	if err := s.urlRepo.IncrementClicks(ctx, url.ID); err != nil {
		s.logger.Error("failed to increment URL clicks",
			logger.ErrorField(err),
			logger.String("urlID", url.ID))
	}

	return url.OriginalURL, nil
}

func (s *urlService) GetURLAnalytics(ctx context.Context, urlID, userID string, from, to time.Time) ([]*models.URLClick, error) {
	url, err := s.urlRepo.GetByID(ctx, urlID)
	if err != nil {
		s.logger.Error("failed to get URL for analytics",
			logger.ErrorField(err),
			logger.String("urlID", urlID))
		return nil, err
	}

	if url.UserID != nil && *url.UserID != userID {
		s.logger.Warn("unauthorized analytics access attempt",
			logger.String("requestingUserID", userID),
			logger.String("urlOwnerID", *url.UserID),
			logger.String("urlID", urlID))
		return nil, models.ErrForbidden
	}

	clicks, err := s.urlRepo.GetClicksAnalytics(ctx, urlID, from, to)
	if err != nil {
		s.logger.Error("failed to get URL analytics",
			logger.ErrorField(err),
			logger.String("urlID", urlID),
			logger.Time("from", from),
			logger.Time("to", to))
		return nil, err
	}

	return clicks, nil
}

func generateShortCode(length int) string {
	uuidStr := uuid.New().String()
	clean := strings.ReplaceAll(uuidStr, "-", "")
	if length > len(clean) {
		length = len(clean)
	}
	return clean[:length]
}

func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z') && !(r >= '0' && r <= '9') {
			return false
		}
	}
	return true
}