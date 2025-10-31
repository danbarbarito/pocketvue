package types

import "time"

// WebhookEvent is the base structure for all Polar webhook events
type WebhookEvent struct {
	Type      string          `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Data      WebhookEventMap `json:"data"`
}

// WebhookEventMap allows flexible parsing of different event data types
type WebhookEventMap map[string]interface{}

// SubscriptionWebhookData represents subscription event data
type SubscriptionWebhookData struct {
	ID                 string                 `json:"id"`
	CreatedAt          time.Time              `json:"created_at"`
	ModifiedAt         time.Time              `json:"modified_at"`
	Amount             int                    `json:"amount"`
	Currency           string                 `json:"currency"`
	RecurringInterval  string                 `json:"recurring_interval"`
	Status             string                 `json:"status"`
	CurrentPeriodStart time.Time              `json:"current_period_start"`
	CurrentPeriodEnd   time.Time              `json:"current_period_end"`
	CancelAtPeriodEnd  bool                   `json:"cancel_at_period_end"`
	CanceledAt         *time.Time             `json:"canceled_at"`
	StartedAt          *time.Time             `json:"started_at"`
	EndsAt             *time.Time             `json:"ends_at"`
	EndedAt            *time.Time             `json:"ended_at"`
	CustomerID         string                 `json:"customer_id"`
	ProductID          string                 `json:"product_id"`
	DiscountID         *string                `json:"discount_id"`
	CheckoutID         *string                `json:"checkout_id"`
	Metadata           map[string]interface{} `json:"metadata"`
	Customer           CustomerData           `json:"customer"`
	Product            ProductData            `json:"product"`
}

// OrderWebhookData represents order event data
type OrderWebhookData struct {
	ID             string                 `json:"id"`
	CreatedAt      time.Time              `json:"created_at"`
	ModifiedAt     time.Time              `json:"modified_at"`
	Status         string                 `json:"status"`
	Paid           bool                   `json:"paid"`
	SubtotalAmount int                    `json:"subtotal_amount"`
	DiscountAmount int                    `json:"discount_amount"`
	NetAmount      int                    `json:"net_amount"`
	TaxAmount      int                    `json:"tax_amount"`
	TotalAmount    int                    `json:"total_amount"`
	Currency       string                 `json:"currency"`
	BillingReason  string                 `json:"billing_reason"`
	CustomerID     string                 `json:"customer_id"`
	ProductID      string                 `json:"product_id"`
	SubscriptionID *string                `json:"subscription_id"`
	CheckoutID     *string                `json:"checkout_id"`
	Metadata       map[string]interface{} `json:"metadata"`
	Customer       CustomerData           `json:"customer"`
	Product        ProductData            `json:"product"`
}

// CustomerData represents customer information in webhook events
type CustomerData struct {
	ID         string     `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	ModifiedAt time.Time  `json:"modified_at"`
	Email      string     `json:"email"`
	Name       *string    `json:"name"`
	ExternalID *string    `json:"external_id"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

// ProductData represents product information in webhook events
type ProductData struct {
	ID          string                 `json:"id"`
	CreatedAt   time.Time              `json:"created_at"`
	ModifiedAt  time.Time              `json:"modified_at"`
	Name        string                 `json:"name"`
	Description *string                `json:"description"`
	IsRecurring bool                   `json:"is_recurring"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ProductWebhookData represents product webhook event data (product.created, product.updated)
type ProductWebhookData struct {
	ID                     string                 `json:"id"`
	CreatedAt              time.Time              `json:"created_at"`
	ModifiedAt             time.Time              `json:"modified_at"`
	TrialInterval          *string                `json:"trial_interval"`
	TrialIntervalCount     *int                   `json:"trial_interval_count"`
	Name                   string                 `json:"name"`
	Description            *string                `json:"description"`
	RecurringInterval      string                 `json:"recurring_interval"`
	RecurringIntervalCount int                    `json:"recurring_interval_count"`
	IsRecurring            bool                   `json:"is_recurring"`
	IsArchived             bool                   `json:"is_archived"`
	OrganizationID         string                 `json:"organization_id"`
	Metadata               map[string]interface{} `json:"metadata"`
	Prices                 []ProductPrice         `json:"prices"`
}

// ProductPrice represents a price within a product
type ProductPrice struct {
	ID                string    `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	ModifiedAt        time.Time `json:"modified_at"`
	AmountType        string    `json:"amount_type"`
	IsArchived        bool      `json:"is_archived"`
	ProductID         string    `json:"product_id"`
	Type              string    `json:"type"`
	RecurringInterval *string   `json:"recurring_interval"`
	PriceCurrency     string    `json:"price_currency"`
	PriceAmount       int       `json:"price_amount"`
	Legacy            bool      `json:"legacy"`
}
