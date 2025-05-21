package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Common validation errors
var (
	ErrEmptyProjectName = errors.New("project name cannot be empty")
	ErrInvalidDateRange = errors.New("end date must be after start date")
	ErrInvalidBudget = errors.New("budget must be non-negative")
	ErrMissingOrganization = errors.New("organization ID is required")
)

// ProjectStatus represents the current status of a project
type ProjectStatus string

const (
	ProjectStatusPlanning   ProjectStatus = "planning"
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusOnHold    ProjectStatus = "on_hold"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusCancelled ProjectStatus = "cancelled"
)

// Project represents a project within an organization
type Project struct {
	gorm.Model
	Name           string        `json:"name" gorm:"not null"`
	Description    string        `json:"description"`
	Status        ProjectStatus `json:"status" gorm:"type:varchar(20);default:'planning'"`
	StartDate     *time.Time    `json:"start_date"`
	EndDate       *time.Time    `json:"end_date"`
	Budget        float64       `json:"budget"`
	OrganizationID uint         `json:"organization_id" gorm:"not null"`
	Organization   Organization  `json:"-" gorm:"foreignKey:OrganizationID"`
	Teams          []Team       `json:"teams" gorm:"many2many:team_projects;"`
	Tasks          []Task       `json:"tasks" gorm:"foreignKey:ProjectID"`
	Members        []User       `json:"members" gorm:"many2many:project_members;"`
}

// ProjectMember represents the many-to-many relationship between projects and users
type ProjectMember struct {
	ProjectID uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"primaryKey"`
	Role      string `json:"role" gorm:"default:'member'"` // manager, member
}

// Validate performs validation on the Project model
 