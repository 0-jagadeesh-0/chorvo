package models

import (
	"errors"
	"net/url"
	"strings"

	"gorm.io/gorm"
)

// Common validation errors
var (
	ErrEmptyOrgName = errors.New("organization name cannot be empty")
	ErrInvalidWebsite = errors.New("invalid website URL")
)

// Organization represents a company or group that can have multiple users,
// projects, and teams. It is the top-level entity in the multi-tenant architecture.
type Organization struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null"`                           // Name of the organization
	Description string `json:"description"`                                    // Optional description
	Website     string `json:"website"`                                       // Optional website URL
	Logo        string `json:"logo"`                                          // Optional logo URL/path
	
	// Relationships
	Users       []User `json:"users" gorm:"many2many:organization_users;"`    // Users belonging to this organization
	Projects    []Project `json:"projects" gorm:"foreignKey:OrganizationID"`  // Projects owned by this organization
	Teams       []Team    `json:"teams" gorm:"foreignKey:OrganizationID"`    // Teams within this organization
	
	// Subscription related
	CurrentSubscription *Subscription `json:"current_subscription" gorm:"foreignKey:OrganizationID"`
	Subscriptions      []Subscription `json:"-" gorm:"foreignKey:OrganizationID"` // Subscription history
	Invoices           []Invoice      `json:"-" gorm:"foreignKey:OrganizationID"` // Billing history
	
	// Billing information
	BillingEmail   string `json:"billing_email"`
	BillingName    string `json:"billing_name"`
	BillingAddress string `json:"billing_address"`
	TaxID          string `json:"tax_id"`
	
	// Feature flags based on subscription
	CustomDomainEnabled bool `json:"custom_domain_enabled" gorm:"default:false"`
	APIAccessEnabled    bool `json:"api_access_enabled" gorm:"default:false"`
	StorageLimit       int  `json:"storage_limit" gorm:"default:5"` // in GB
}

// OrganizationUser represents the many-to-many relationship between
// organizations and users, including the user's role within the organization.
type OrganizationUser struct {
	OrganizationID uint   `gorm:"primaryKey"`                    // Foreign key to Organization
	UserID         uint   `gorm:"primaryKey"`                    // Foreign key to User
	Role          string `json:"role" gorm:"default:'member'"` // User's role: admin or member
}

// Validate performs validation on the Organization model
func (o *Organization) Validate() error {
	if strings.TrimSpace(o.Name) == "" {
		return ErrEmptyOrgName
	}

	if o.Website != "" {
		if err := o.validateWebsite(); err != nil {
			return err
		}
	}

	return nil
}

// validateWebsite checks if the website URL is valid
func (o *Organization) validateWebsite() error {
	_, err := url.ParseRequestURI(o.Website)
	if err != nil {
		return ErrInvalidWebsite
	}
	return nil
}

// HasActiveSubscription checks if the organization has an active subscription
func (o *Organization) HasActiveSubscription() bool {
	return o.CurrentSubscription != nil && o.CurrentSubscription.IsActive()
}

// IsTrialing checks if the organization is in trial period
func (o *Organization) IsTrialing() bool {
	return o.CurrentSubscription != nil && o.CurrentSubscription.IsTrialing()
}

// HasFeature checks if the organization has access to a specific feature
func (o *Organization) HasFeature(featureName string) bool {
	if o.CurrentSubscription == nil {
		return false
	}
	return o.CurrentSubscription.HasFeature(featureName)
}

// WithinLimits checks if the organization is within its subscription limits
func (o *Organization) WithinLimits() bool {
	if o.CurrentSubscription == nil {
		return false
	}
	return o.CurrentSubscription.WithinLimits()
}

// BeforeCreate is a GORM hook that runs before creating a new organization
func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	return o.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating an organization
func (o *Organization) BeforeUpdate(tx *gorm.DB) error {
	return o.Validate()
} 