package services

import (
	"context"
	"errors"
	"time"

	"github.com/0-jagadeesh-0/chorvo/internal/api/v1/utils"
	"github.com/0-jagadeesh-0/chorvo/internal/domain/models"
	"github.com/0-jagadeesh-0/chorvo/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

var (
	ErrUserExists          = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCode        = errors.New("invalid verification code")
	ErrCodeExpired        = errors.New("verification code expired")
)

func (s *AuthService) Register(email, password, firstName, lastName string) (*models.User, error) {
	ctx := context.Background()
	
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Generate verification code
	verificationCode, err := utils.GenerateVerificationCode()
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &models.User{
		Email:            email,
		Password:         string(hashedPassword),
		FirstName:        firstName,
		LastName:         lastName,
		Status:           models.UserStatusInactive,
		IsActive:         false,
		VerificationCode: verificationCode,
		CodeExpiresAt:    time.Now().Add(10 * time.Minute),
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Send verification email
	if err := utils.SendVerificationEmail(user.Email, verificationCode); err != nil {
		// Log the error but don't fail the registration
		// In a production environment, you might want to implement retry logic
		return user, nil
	}

	return user, nil
}

func (s *AuthService) VerifyEmail(email, code string) error {
	ctx := context.Background()
	
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}

	if user.VerificationCode != code {
		return ErrInvalidCode
	}

	if time.Now().After(user.CodeExpiresAt) {
		return ErrCodeExpired
	}

	now := time.Now()
	user.Status = models.UserStatusActive
	user.IsActive = true
	user.EmailVerifiedAt = &now
	user.VerificationCode = ""
	user.CodeExpiresAt = time.Time{}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResendVerificationCode(email string) error {
	ctx := context.Background()
	
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}

	// Check if email is already verified
	if user.EmailVerifiedAt != nil {
		return errors.New("email is already verified")
	}

	// Generate new verification code
	verificationCode, err := utils.GenerateVerificationCode()
	if err != nil {
		return err
	}

	// Update user with new verification code
	user.VerificationCode = verificationCode
	user.CodeExpiresAt = time.Now().Add(10 * time.Minute)

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Send new verification email
	return utils.SendVerificationEmail(user.Email, verificationCode)
}

func (s *AuthService) Login(email, password string) (string, error) {
	ctx := context.Background()
	
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Check if email is verified
	if user.EmailVerifiedAt == nil {
		return "", errors.New("email not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(ctx, user)

	return token, nil
}

func (s *AuthService) RequestPasswordReset(email string) error {
	ctx := context.Background()
	
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}

	// Generate reset token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return err
	}

	// Store reset token with expiry
	user.ResetToken = token
	user.ResetTokenExpiry = time.Now().Add(1 * time.Hour)
	
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Send reset email
	return utils.SendPasswordResetEmail(user.Email, token)
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Validate token
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return err
	}

	ctx := context.Background()
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil || user.ResetToken != token {
		return ErrInvalidCredentials
	}

	if time.Now().After(user.ResetTokenExpiry) {
		return ErrCodeExpired
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password and clear reset token
	user.Password = string(hashedPassword)
	user.ResetToken = ""
	user.ResetTokenExpiry = time.Time{}
	
	return s.userRepo.Update(ctx, user)
}