package models

import "time"

type Subscription struct {
	ID                 int64  `db:"id"` // ID interno no banco
	UserID             int64  `db:"user_id"`
	StripeCustomer     string `db:"stripe_customer"`
	StripeSubscription string `db:"stripe_subscription"`
	StripePriceID      string `db:"stripe_price_id"`
	StripeProductID    string `db:"stripe_product_id"`
	StripePlanAmount   int64  `db:"stripe_plan_amount"`
	Currency           string `db:"currency"`
	Interval           string `db:"interval"`
	IntervalCount      int64  `db:"interval_count"`
	Status             string `db:"status"`
}

type Account struct {
	ID              string    `json:"id"`
	Nome            string    `json:"nome"`
	Email           string    `json:"email"`
	StripeProductID string    `json:"stripe_product_id"`
	CreatedAt       time.Time `json:"created_at"`
	Status          bool      `json:"status"`
}

type StripeCustomer struct {
	ID               string   `json:"id"`
	Object           string   `json:"object"`
	Email            string   `json:"email"`
	Name             string   `json:"name"`
	Currency         string   `json:"currency"`
	InvoicePrefix    string   `json:"invoice_prefix"`
	NextInvoiceSeq   int      `json:"next_invoice_sequence"`
	TaxExempt        string   `json:"tax_exempt"`
	Phone            string   `json:"phone"`
	Created          int64    `json:"created"`
	Delinquent       bool     `json:"delinquent"`
	Livemode         bool     `json:"livemode"`
	PreferredLocales []string `json:"preferred_locales"`

	Address struct {
		City       string `json:"city"`
		Country    string `json:"country"`
		Line1      string `json:"line1"`
		Line2      string `json:"line2"`
		PostalCode string `json:"postal_code"`
		State      string `json:"state"`
	} `json:"address"`

	InvoiceSettings struct {
		CustomFields         interface{} `json:"custom_fields"`
		DefaultPaymentMethod interface{} `json:"default_payment_method"`
		Footer               interface{} `json:"footer"`
		RenderingOptions     interface{} `json:"rendering_options"`
	} `json:"invoice_settings"`

	Metadata map[string]string `json:"metadata"`
}
