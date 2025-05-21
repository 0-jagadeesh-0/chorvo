package models

import "gorm.io/gorm"

// Team represents a group of users working together on projects
type Team struct {
    gorm.Model
    Name           string `json:"name" gorm:"not null"`
    Description    string `json:"description"`
    OrganizationID uint   `json:"organization_id" gorm:"not null"`
    Organization   Organization `json:"-" gorm:"foreignKey:OrganizationID"`
    Members        []User  `json:"members" gorm:"many2many:team_members;"`
    Projects       []Project `json:"projects" gorm:"many2many:team_projects;"`
}

// TeamMember represents the many-to-many relationship between teams and users
type TeamMember struct {
    TeamID  uint   `gorm:"primaryKey"`
    UserID  uint   `gorm:"primaryKey"`
    Role    string `json:"role" gorm:"default:'member'"` // lead, member
}

// TeamProject represents the many-to-many relationship between teams and projects
type TeamProject struct {
    TeamID    uint `gorm:"primaryKey"`
    ProjectID uint `gorm:"primaryKey"`
} 