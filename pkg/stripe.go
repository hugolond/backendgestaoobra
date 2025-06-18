package pkg

import (
	models "backendgestaoobra/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetStripeCustomer(stripeCustomerID string) (*models.StripeCustomer, error) {
	stripeSecret := os.Getenv("STRIPE_SECRET_KEY")
	if stripeSecret == "" {
		return nil, fmt.Errorf("Stripe secret key não definido no .env")
	}

	url := fmt.Sprintf("https://api.stripe.com/v1/customers/%s", stripeCustomerID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+stripeSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na resposta do Stripe: %s", string(bodyBytes))
	}

	var customer models.StripeCustomer
	if err := json.Unmarshal(bodyBytes, &customer); err != nil {
		return nil, err
	}

	return &customer, nil
}
