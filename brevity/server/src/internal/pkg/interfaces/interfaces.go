package interfaces

// import (
// 	"context"
// )

// type AuthProvider interface {
// 	GenerateToken(userID string) (string, error)
// 	ValidateToken(token string) (string, error)
// 	GenerateRefreshToken() (string, error)
// }

// type EmailService interface {
// 	SendVerificationEmail(email, token string) error
// 	SendPasswordResetEmail(email, token string) error
// }

// type StorageProvider interface {
// 	UploadFile(ctx context.Context, file []byte, filename string) (string, error)
// 	DeleteFile(ctx context.Context, filename string) error
// }