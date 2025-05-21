package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// PaymentStatus represents the current status of a payment
type PaymentStatus string

const (
    PaymentStatusPending   PaymentStatus = "pending"
    PaymentStatusSucceeded PaymentStatus = "succeeded"
    PaymentStatusFailed    PaymentStatus = "failed"
    PaymentStatusRefunded  PaymentStatus = "refunded"
)

// PaymentMethod represents the method used for payment
type PaymentMethod string

const (
    PaymentMethodCard    PaymentMethod = "card"
    PaymentMethodPayPal  PaymentMethod = "paypal"
    PaymentMethodBank    PaymentMethod = "bank_transfer"
)

// Invoice represents a billing invoice
type Invoice struct {
    gorm.Model
    OrganizationID  uint          `json:"organization_id" gorm:"not null"`
    Organization    Organization  `json:"-" gorm:"foreignKey:OrganizationID"`
    SubscriptionID  uint          `json:"subscription_id" gorm:"not null"`
    Subscription    Subscription  `json:"-" gorm:"foreignKey:SubscriptionID"`
    
    InvoiceNumber   string        `json:"invoice_number" gorm:"unique;not null"`
    Amount          float64       `json:"amount" gorm:"not null"`
    Currency        string        `json:"currency" gorm:"default:USD"`
    DueDate         time.Time     `json:"due_date"`
    PaidAt         *time.Time     `json:"paid_at"`
    Status         PaymentStatus  `json:"status" gorm:"type:varchar(20);default:'pending'"`
    
    // Billing details
    BillingName    string        `json:"billing_name"`
    BillingEmail   string        `json:"billing_email"`
    BillingAddress string        `json:"billing_address"`
    
    // Payment details
    PaymentMethod  PaymentMethod `json:"payment_method" gorm:"type:varchar(20)"`
    PaymentID      string        `json:"payment_id"` // External payment reference
    
    // Line items
    Items         []InvoiceItem `json:"items" gorm:"foreignKey:InvoiceID"`
}

// InvoiceItem represents a line item in an invoice
type InvoiceItem struct {
    gorm.Model
    InvoiceID     uint    `json:"invoice_id" gorm:"not null"`
    Description   string  `json:"description" gorm:"not null"`
    Quantity      int     `json:"quantity" gorm:"not null"`
    UnitPrice     float64 `json:"unit_price" gorm:"not null"`
    Amount        float64 `json:"amount" gorm:"not null"`
}

// PaymentTransaction represents a payment transaction
type PaymentTransaction struct {
    gorm.Model
    InvoiceID     uint          `json:"invoice_id" gorm:"not null"`
    Invoice       Invoice       `json:"-" gorm:"foreignKey:InvoiceID"`
    Amount        float64       `json:"amount" gorm:"not null"`
    Currency      string        `json:"currency" gorm:"default:USD"`
    Status        PaymentStatus `json:"status" gorm:"type:varchar(20)"`
    PaymentMethod PaymentMethod `json:"payment_method" gorm:"type:varchar(20)"`
    
    // Payment provider details
    ProviderID    string        `json:"provider_id"` // Payment provider's transaction ID
    ProviderFee   float64       `json:"provider_fee"`
    
    // Error handling
    ErrorCode     string        `json:"error_code"`
    ErrorMessage  string        `json:"error_message"`
}

// Validate performs validation on the Invoice model
func (i *Invoice) Validate() error {
    if i.OrganizationID == 0 {
        return errors.New("organization ID is required")
    }

    if i.SubscriptionID == 0 {
        return errors.New("subscription ID is required")
    }

    if i.Amount < 0 {
        return errors.New("amount must be non-negative")
    }

    if i.InvoiceNumber == "" {
        return errors.New("invoice number is required")
    }

    return nil
}

// IsPaid checks if the invoice has been paid
func (i *Invoice) IsPaid() bool {
    return i.Status == PaymentStatusSucceeded && i.PaidAt != nil
}

// IsOverdue checks if the invoice is overdue
func (i *Invoice) IsOverdue() bool {
    return !i.IsPaid() && time.Now().After(i.DueDate)
}

// CalculateTotal calculates the total amount for the invoice
func (i *Invoice) CalculateTotal() float64 {
    var total float64
    for _, item := range i.Items {
        total += item.Amount
    }
    return total
}

// BeforeCreate is a GORM hook that runs before creating a new invoice
func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
    return i.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating an invoice
func (i *Invoice) BeforeUpdate(tx *gorm.DB) error {
    return i.Validate()
} 