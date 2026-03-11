package email

import (
	"fmt"
	"net/smtp"
)

// Config holds email configuration
type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// Service handles email sending
type Service struct {
	config Config
}

// NewService creates a new email service
func NewService(config Config) *Service {
	return &Service{config: config}
}

// SendConfirmationEmail sends an email confirmation link to the user
func (s *Service) SendConfirmationEmail(to, username, token, baseURL string) error {
	confirmURL := fmt.Sprintf("%s/confirm-email?token=%s", baseURL, token)
	
	subject := "Confirm your email address"
	body := fmt.Sprintf(`
Hello %s,

Thank you for registering! Please confirm your email address by clicking the link below:

%s

This link will expire in 24 hours.

If you didn't create an account, you can safely ignore this email.

Best regards,
Divergence 2%% Writer Portal Team
`, username, confirmURL)

	return s.SendEmail(to, subject, body)
}

// SendEmail sends a plain text email
func (s *Service) SendEmail(to, subject, body string) error {
	from := s.config.FromEmail
	
	// Set up authentication
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
	
	// Compose message
	msg := []byte(fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.config.FromName, from, to, subject, body))
	
	// Send email
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	
	return nil
}

