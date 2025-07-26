package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// APIResponse represents a standardized API response structure
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
	Path      string      `json:"path"`
}

// ResponseOptions configures the API response
type ResponseOptions struct {
	Success bool
	Message string
	Data    interface{}
	Error   error
}

// GetValidationErrors converts validator errors to a map
func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = err.Tag()
	}
	return errors
}

// Respond sends a standardized JSON response
func Respond(c *gin.Context, statusCode int, opts ResponseOptions) {
	response := APIResponse{
		Success:   opts.Success,
		Message:   opts.Message,
		Data:      opts.Data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Path:      c.Request.URL.Path,
	}

	if opts.Error != nil {
		response.Error = opts.Error.Error()
	}

	c.JSON(statusCode, response)
}

// Success sends a successful response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	Respond(c, statusCode, ResponseOptions{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, message string, err error) {
	Respond(c, statusCode, ResponseOptions{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success:   false,
		Message:   "Validation failed",
		Error:     "validation_error",
		Data:      gin.H{"errors": errors},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Path:      c.Request.URL.Path,
	})
}

// PaginatedResponse sends a paginated response
func PaginatedResponse(c *gin.Context, statusCode int, data interface{}, pagination interface{}) {
	Respond(c, statusCode, ResponseOptions{
		Success: true,
		Message: "Data retrieved successfully",
		Data: gin.H{
			"items":      data,
			"pagination": pagination,
		},
	})
}
