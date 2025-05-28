package models

type Subscription struct {
	ID             int64
	UserID         int64
	StripeCustomer string
	StripeSession  string
	Plan           string
	Status         string
}
