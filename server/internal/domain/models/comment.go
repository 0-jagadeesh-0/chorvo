package models

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// Common validation errors
var (
    ErrEmptyContent = errors.New("comment content cannot be empty")
    ErrMissingTask = errors.New("task ID is required")
    ErrMissingUser = errors.New("user ID is required")
)

// Comment represents a comment on a task
type Comment struct {
    gorm.Model
    Content  string `json:"content" gorm:"not null"`
    TaskID   uint   `json:"task_id" gorm:"not null"`
    Task     Task   `json:"-" gorm:"foreignKey:TaskID"`
    UserID   uint   `json:"user_id" gorm:"not null"`
    User     User   `json:"user" gorm:"foreignKey:UserID"`
    ParentID *uint  `json:"parent_id"`
    Parent   *Comment `json:"-" gorm:"foreignKey:ParentID"`
    Replies  []Comment `json:"replies" gorm:"foreignKey:ParentID"`
}

// Validate performs validation on the Comment model
func (c *Comment) Validate() error {
    if strings.TrimSpace(c.Content) == "" {
        return ErrEmptyContent
    }

    if c.TaskID == 0 {
        return ErrMissingTask
    }

    if c.UserID == 0 {
        return ErrMissingUser
    }

    return nil
}

// IsReply checks if this comment is a reply to another comment
func (c *Comment) IsReply() bool {
    return c.ParentID != nil
}

// BeforeCreate is a GORM hook that runs before creating a new comment
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
    return c.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a comment
func (c *Comment) BeforeUpdate(tx *gorm.DB) error {
    return c.Validate()
} 