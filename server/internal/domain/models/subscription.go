package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Common validation errors
var (
    ErrInvalidPlanID = errors.New("invalid plan ID")
    ErrInvalidPrice = errors.New("price must be non-negative")
    ErrInvalidDuration = errors.New("invalid subscription duration")
)

// SubscriptionStatus represents the current status of a subscription
type SubscriptionStatus string

const (
    SubscriptionStatusActive    SubscriptionStatus = "active"
    SubscriptionStatusInactive  SubscriptionStatus = "inactive"
    SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
    SubscriptionStatusPending   SubscriptionStatus = "pending"
)

// BillingInterval represents the billing frequency
type BillingInterval string

const (
    BillingIntervalMonthly BillingInterval = "monthly"
    BillingIntervalYearly  BillingInterval = "yearly"
)

// Plan represents a subscription plan with its features
type Plan struct {
    gorm.Model
    Name           string  `json:"name" gorm:"not null"`
    Description    string  `json:"description"`
    Price         float64 `json:"price" gorm:"not null"`
    BillingInterval BillingInterval `json:"billing_interval" gorm:"type:varchar(20);default:'monthly'"`
    Features      []PlanFeature `json:"features" gorm:"foreignKey:PlanID"`
    
    // Plan limits
    MaxUsers      int `json:"max_users" gorm:"not null"`
    MaxProjects   int `json:"max_projects" gorm:"not null"`
    MaxStorage    int `json:"max_storage" gorm:"not null"` // in GB
    
    // Additional features
    CustomDomain  bool `json:"custom_domain" gorm:"default:false"`
    APIAccess     bool `json:"api_access" gorm:"default:false"`
    Priority      bool `json:"priority_support" gorm:"default:false"`
}

// PlanFeature represents a feature available in a plan
type PlanFeature struct {
    gorm.Model
    PlanID      uint   `json:"plan_id" gorm:"not null"`
    Name        string `json:"name" gorm:"not null"`
    Description string `json:"description"`
    Included    bool   `json:"included" gorm:"default:true"`
}

// Subscription represents an organization's subscription to a plan
type Subscription struct {
    gorm.Model
    OrganizationID uint              `json:"organization_id" gorm:"not null"`
    Organization   Organization      `json:"-" gorm:"foreignKey:OrganizationID"`
    PlanID         uint              `json:"plan_id" gorm:"not null"`
    Plan           Plan              `json:"plan" gorm:"foreignKey:PlanID"`
    Status         SubscriptionStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
    
    // Billing dates
    StartDate      time.Time         `json:"start_date"`
    EndDate        time.Time         `json:"end_date"`
    TrialEndsAt    *time.Time        `json:"trial_ends_at"`
    
    // Payment info
    PaymentMethod  string            `json:"payment_method"`
    LastBilledAt   *time.Time        `json:"last_billed_at"`
    NextBillingAt  *time.Time        `json:"next_billing_at"`
    
    // Usage tracking
    CurrentUsers   int               `json:"current_users"`
    CurrentProjects int             `json:"current_projects"`
    CurrentStorage  float64         `json:"current_storage"` // in GB
}

// Validate performs validation on the Subscription model
func (s *Subscription) Validate() error {
    if s.OrganizationID == 0 {
        return errors.New("organization ID is required")
    }

    if s.PlanID == 0 {
        return ErrInvalidPlanID
    }

    if s.EndDate.Before(s.StartDate) {
        return ErrInvalidDuration
    }

    return nil
}

// IsActive checks if the subscription is currently active
func (s *Subscription) IsActive() bool {
    return s.Status == SubscriptionStatusActive && time.Now().Before(s.EndDate)
}

// IsTrialing checks if the subscription is in trial period
func (s *Subscription) IsTrialing() bool {
    return s.TrialEndsAt != nil && time.Now().Before(*s.TrialEndsAt)
}

// HasFeature checks if the subscription's plan includes a specific feature
func (s *Subscription) HasFeature(featureName string) bool {
    for _, feature := range s.Plan.Features {
        if feature.Name == featureName && feature.Included {
            return true
        }
    }
    return false
}

// WithinLimits checks if the organization is within the plan's limits
func (s *Subscription) WithinLimits() bool {
    return s.CurrentUsers <= s.Plan.MaxUsers &&
           s.CurrentProjects <= s.Plan.MaxProjects &&
           s.CurrentStorage <= float64(s.Plan.MaxStorage)
}

// BeforeCreate is a GORM hook that runs before creating a new subscription
func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
    return s.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a subscription
func (s *Subscription) BeforeUpdate(tx *gorm.DB) error {
    return s.Validate()
} 