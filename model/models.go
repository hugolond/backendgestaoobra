package models

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
