package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"os"
)

type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

var emailConfig *EmailConfig

func InitEmailConfig() {
	emailConfig = &EmailConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}
}

func GenerateVerificationCode() (string, error) {
	bytes := make([]byte, 3) // 6 characters in hex
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func SendVerificationEmail(to, code string) error {
	subject := "Verify Your Email - Chorvo"
	body := fmt.Sprintf(`
		<h2>Welcome to Chorvo!</h2>
		<p>Your verification code is: <strong>%s</strong></p>
		<p>This code will expire in 10 minutes.</p>
	`, code)

	return sendEmail(to, subject, body)
}

func SendPasswordResetEmail(to, token string) error {
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), token)
	subject := "Reset Your Password - Chorvo"
	body := fmt.Sprintf(`
		<h2>Password Reset Request</h2>
		<p>Click the link below to reset your password:</p>
		<p><a href="%s">Reset Password</a></p>
		<p>This link will expire in 1 hour.</p>
		<p>If you didn't request this, please ignore this email.</p>
	`, resetLink)

	return sendEmail(to, subject, body)
}

func sendEmail(to, subject, body string) error {
	if emailConfig == nil {
		InitEmailConfig()
	}

	auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.Host)
	
	msg := fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", to, emailConfig.From, subject, body)

	addr := fmt.Sprintf("%s:%s", emailConfig.Host, emailConfig.Port)
	return smtp.SendMail(addr, auth, emailConfig.From, []string{to}, []byte(msg))
} 