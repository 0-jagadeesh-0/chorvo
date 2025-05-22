package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Common validation errors
var (
	ErrEmptyTaskTitle = errors.New("task title cannot be empty")
	ErrMissingProject = errors.New("project ID is required")
	ErrMissingCreator = errors.New("creator ID is required")
	ErrInvalidHours = errors.New("hours must be non-negative")
	ErrInvalidTaskDates = errors.New("completion date must be after start date")
)

// TaskPriority represents the priority level of a task
type TaskPriority string

// TaskStatus represents the current status of a task
type TaskStatus string

const (
	TaskPriorityLow     TaskPriority = "low"
	TaskPriorityMedium  TaskPriority = "medium"
	TaskPriorityHigh    TaskPriority = "high"
	TaskPriorityCritical TaskPriority = "critical"
)

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusInReview   TaskStatus = "in_review"
	TaskStatusDone       TaskStatus = "done"
)

// Task represents a unit of work within a project
type Task struct {
	gorm.Model
	Title       string       `json:"title" gorm:"not null"`
	Description string       `json:"description"`
	Priority    TaskPriority `json:"priority" gorm:"type:varchar(20);default:'medium'"`
	Status      TaskStatus   `json:"status" gorm:"type:varchar(20);default:'todo'"`
	DueDate     *time.Time   `json:"due_date"`
	
	// Project relationship
	ProjectID   uint         `json:"project_id" gorm:"not null"`
	Project     Project      `json:"-" gorm:"foreignKey:ProjectID"`
	
	// User relationships
	CreatedByID uint         `json:"created_by_id" gorm:"not null"`
	CreatedBy   User         `json:"created_by" gorm:"foreignKey:CreatedByID"`
	AssigneeID  *uint        `json:"assignee_id"`
	Assignee    *User        `json:"assignee" gorm:"foreignKey:AssigneeID"`
	
	// Task hierarchy
	ParentID    *uint        `json:"parent_id"`
	Parent      *Task        `json:"-" gorm:"foreignKey:ParentID"`
	Subtasks    []Task       `json:"subtasks" gorm:"foreignKey:ParentID"`
	
	// Related entities
	Comments    []Comment    `json:"comments" gorm:"foreignKey:TaskID"`
	
	// Time tracking
	EstimatedHours float32   `json:"estimated_hours"`
	ActualHours    float32   `json:"actual_hours"`
	StartedAt      *time.Time `json:"started_at"`
	CompletedAt    *time.Time `json:"completed_at"`
}

// Validate performs validation on the Task model
func (t *Task) Validate() error {
	if strings.TrimSpace(t.Title) == "" {
		return ErrEmptyTaskTitle
	}

	if t.ProjectID == 0 {
		return ErrMissingProject
	}

	if t.CreatedByID == 0 {
		return ErrMissingCreator
	}

	if t.EstimatedHours < 0 || t.ActualHours < 0 {
		return ErrInvalidHours
	}

	if err := t.validateDates(); err != nil {
		return err
	}

	return nil
}

// validateDates checks if the task dates are valid
func (t *Task) validateDates() error {
	if t.StartedAt != nil && t.CompletedAt != nil {
		if t.CompletedAt.Before(*t.StartedAt) {
			return ErrInvalidTaskDates
		}
	}
	return nil
}

// IsComplete checks if the task is marked as done
func (t *Task) IsComplete() bool {
	return t.Status == TaskStatusDone
}

// IsInProgress checks if the task is currently being worked on
func (t *Task) IsInProgress() bool {
	return t.Status == TaskStatusInProgress
}

// BeforeCreate is a GORM hook that runs before creating a new task
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	return t.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a task
func (t *Task) BeforeUpdate(tx *gorm.DB) error {
	return t.Validate()
} 