package email

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/imraushankr/bervity/server/src/configs"
	"github.com/imraushankr/bervity/server/src/internal/pkg/logger"
)

const (
	emailTimeout = 15 * time.Second
)

type EmailService struct {
	cfg    *configs.EmailConfig
	logger logger.Logger
}

func NewEmailService(cfg *configs.EmailConfig, log logger.Logger) *EmailService {
	return &EmailService{
		cfg:    cfg,
		logger: log,
	}
}

func (e *EmailService) SendVerificationEmail(to, verificationLink string) error {
	const subject = "Verify Your Email Address for Brevity"
	body := fmt.Sprintf(`
		<html>
		<head>
			<style>
				body { font-family: 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { text-align: center; margin-bottom: 30px; }
				.logo { color: #2563eb; font-size: 24px; font-weight: bold; margin-bottom: 10px; }
				.content { background-color: #f9fafb; padding: 25px; border-radius: 8px; }
				.button { display: inline-block; background-color: #2563eb; color: white !important; text-decoration: none; padding: 12px 24px; border-radius: 6px; font-weight: 500; margin: 20px 0; }
				.footer { margin-top: 30px; font-size: 12px; color: #6b7280; text-align: center; }
				hr { border: none; height: 1px; background-color: #e5e7eb; margin: 25px 0; }
			</style>
		</head>
		<body>
			<div class="header">
				<div class="logo">Brevity</div>
				<h2 style="margin: 0; font-weight: 500;">Verify Your Email Address</h2>
			</div>
			
			<div class="content">
				<p>Welcome to Brevity! We're excited to have you on board.</p>
				<p>To complete your registration, please verify your email address by clicking the button below:</p>
				
				<div style="text-align: center;">
					<a href="%s" class="button">Verify Email Address</a>
				</div>
				
				<p style="font-size: 14px; color: #6b7280;">If you didn't create an account with Brevity, you can safely ignore this email.</p>
			</div>
			
			<div class="footer">
				<hr>
				<p>This verification link will expire in 24 hours.</p>
				<p>&copy; %d Brevity. All rights reserved.</p>
			</div>
		</body>
		</html>
	`, verificationLink, time.Now().Year())

	return e.sendEmail(to, subject, body)
}

func (e *EmailService) SendPasswordResetEmail(to, resetLink string) error {
	const subject = "Reset Your Brevity Password"
	body := fmt.Sprintf(`
		<html>
		<head>
			<style>
				body { font-family: 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { text-align: center; margin-bottom: 30px; }
				.logo { color: #2563eb; font-size: 24px; font-weight: bold; margin-bottom: 10px; }
				.content { background-color: #f9fafb; padding: 25px; border-radius: 8px; }
				.button { display: inline-block; background-color: #dc2626; color: white !important; text-decoration: none; padding: 12px 24px; border-radius: 6px; font-weight: 500; margin: 20px 0; }
				.footer { margin-top: 30px; font-size: 12px; color: #6b7280; text-align: center; }
				hr { border: none; height: 1px; background-color: #e5e7eb; margin: 25px 0; }
				.warning { background-color: #fef2f2; padding: 12px; border-radius: 6px; border-left: 4px solid #dc2626; margin: 15px 0; }
			</style>
		</head>
		<body>
			<div class="header">
				<div class="logo">Brevity</div>
				<h2 style="margin: 0; font-weight: 500;">Password Reset Request</h2>
			</div>
			
			<div class="content">
				<p>We received a request to reset your Brevity account password.</p>
				
				<div style="text-align: center;">
					<a href="%s" class="button">Reset Password</a>
				</div>
				
				<p>This link will expire in 15 minutes for your security.</p>
				
				<div class="warning">
					<p style="margin: 0; color: #dc2626;">If you didn't request this password reset, please secure your account immediately as someone else may be trying to access it.</p>
				</div>
			</div>
			
			<div class="footer">
				<hr>
				<p>For security reasons, we don't store your password. This link gives you access to reset it.</p>
				<p>&copy; %d Brevity Security Team</p>
			</div>
		</body>
		</html>
	`, resetLink, time.Now().Year())

	return e.sendEmail(to, subject, body)
}

func (e *EmailService) sendEmail(to, subject, body string) error {
	from := e.cfg.SMTP.FromEmail
	if from == "" {
		return fmt.Errorf("from email address not configured")
	}

	// Validate recipient email
	if to == "" || !strings.Contains(to, "@") {
		return fmt.Errorf("invalid recipient email address")
	}

	// Construct MIME email
	headers := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n" + body)

	auth := smtp.PlainAuth("", e.cfg.SMTP.Username, e.cfg.SMTP.Password, e.cfg.SMTP.Host)
	addr := fmt.Sprintf("%s:%d", e.cfg.SMTP.Host, e.cfg.SMTP.Port)

	// Create a channel to handle SMTP send with timeout
	done := make(chan error, 1)
	go func() {
		err := smtp.SendMail(
			addr,
			auth,
			from,
			strings.Split(to, ";"),
			[]byte(msg.String()),
		)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			e.logger.Error("Failed to send email",
				logger.ErrorField(err),
				logger.String("to", to),
				logger.String("subject", subject))
			return fmt.Errorf("email delivery failed: %w", err)
		}
		e.logger.Info("Email sent successfully",
			logger.String("to", to),
			logger.String("subject", subject))
		return nil
	case <-time.After(emailTimeout):
		e.logger.Error("Email sending timed out",
			logger.String("to", to),
			logger.String("subject", subject))
		return fmt.Errorf("email sending timed out after %s", emailTimeout)
	}
}
