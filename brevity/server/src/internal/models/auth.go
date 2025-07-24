// package models

// // Request/Response structs
// type RegisterRequest struct {
// 	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
// 	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
// 	Username  string `json:"username" validate:"required,min=3,max=30,alphanum"`
// 	Email     string `json:"email" validate:"required,email"`
// 	Password  string `json:"password" validate:"required,min=8"`
// }

// type LoginRequest struct {
// 	UserID   string `json:"userId" validate:"required"`
// 	Password string `json:"password" validate:"required,min=8"`
// }

// type LoginResponse struct {
// 	User        User   `json:"user"`
// 	AccessToken string `json:"access_token"`
// 	TokenType   string `json:"token_type"`
// 	ExpiresIn   int    `json:"expires_in"`
// }

// type PasswordResetRequest struct {
// 	Email string `json:"email" validate:"required,email"`
// }

// type CompletePasswordResetRequest struct {
// 	Token       string `json:"-"`
// 	NewPassword string `json:"new_password" validate:"required,min=8"`
// }

// type RefreshTokenRequest struct {
// 	RefreshToken string `json:"refresh_token" validate:"required"`
// }

// type RefreshTokenResponse struct {
// 	AccessToken string `json:"access_token"`
// 	TokenType   string `json:"token_type"`
// 	ExpiresIn   int    `json:"expires_in"`
// }

// type ChangePasswordRequest struct {
// 	CurrentPassword string `json:"current_password" validate:"required,min=8"`
// 	NewPassword     string `json:"new_password" validate:"required,min=8"`
// }

// type MessageResponse struct {
// 	Message string `json:"message"`
// }

package models