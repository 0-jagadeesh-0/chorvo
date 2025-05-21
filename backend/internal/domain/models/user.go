package models

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Common validation errors
var (
	ErrEmptyFirstName = errors.New("first name cannot be empty")
	ErrEmptyLastName  = errors.New("last name cannot be empty")
	ErrInvalidEmail   = errors.New("invalid email format")
	ErrEmptyPassword  = errors.New("password cannot be empty")
	ErrInvalidPhone   = errors.New("invalid phone number format")
)

// UserStatus represents the current status of a user account
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// User represents a user in the system who can be part of organizations,
// teams, and projects. Users can create and be assigned to tasks.
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"` // "-" means don't show in JSON
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	IsActive  bool          `json:"is_active" gorm:"default:false"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Basic Information
	PhoneNumber  string     `json:"phone_number"`
	Avatar       string     `json:"avatar"`
	
	// Account Status
	Status       UserStatus `json:"status" gorm:"type:varchar(20);default:'inactive'"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	
	// Email Verification
	VerificationCode string     `json:"-" gorm:"size:6"`
	CodeExpiresAt   time.Time  `json:"-"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	
	// Password Reset
	ResetToken       string     `json:"-"`
	ResetTokenExpiry time.Time  `json:"-"`
	
	// Preferences
	TimeZone     string     `json:"time_zone" gorm:"default:'UTC'"`
	Language     string     `json:"language" gorm:"default:'en'"`
	
	// Relationships
	Organizations []Organization `json:"organizations" gorm:"many2many:organization_users;"`
	Teams         []Team        `json:"teams" gorm:"many2many:team_members;"`
 
	// Task Related
	AssignedTasks []Task     `json:"assigned_tasks" gorm:"foreignKey:AssigneeID"`
	CreatedTasks  []Task     `json:"created_tasks" gorm:"foreignKey:CreatedByID"`
	Comments      []Comment  `json:"comments" gorm:"foreignKey:UserID"`
	
	// Notifications
	EmailNotifications bool `json:"email_notifications" gorm:"default:true"`
	PushNotifications  bool `json:"push_notifications" gorm:"default:true"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// Validate performs validation on the User model
func (u *User) Validate() error {
	if strings.TrimSpace(u.FirstName) == "" {
		return ErrEmptyFirstName
	}

	if strings.TrimSpace(u.LastName) == "" {
		return ErrEmptyLastName
	}

	if err := u.validateEmail(); err != nil {
		return err
	}

	if strings.TrimSpace(u.Password) == "" {
		return ErrEmptyPassword
	}

	if u.PhoneNumber != "" {
		if err := u.validatePhoneNumber(); err != nil {
			return err
		}
	}

	return nil
}

// validateEmail checks if the email is in a valid format
func (u *User) validateEmail() error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return ErrInvalidEmail
	}
	return nil
}

// validatePhoneNumber checks if the phone number is in a valid format
func (u *User) validatePhoneNumber() error {
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(strings.ReplaceAll(u.PhoneNumber, " ", "")) {
		return ErrInvalidPhone
	}
	return nil
}

// BeforeCreate is a GORM hook that runs before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a user
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return u.Validate()
}