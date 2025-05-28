package src

import (
	models "backendgestaoobra/model"
	"backendgestaoobra/pkg"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

type Subscription struct {
	ID             int64
	UserID         int64
	StripeCustomer string
	StripeSession  string
	Plan           string
	Status         string
}

func StripeWebhookHandler(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusRequestEntityTooLarge, "Request too large")
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	sigHeader := c.GetHeader("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		log.Println("⚠️  Webhook signature verification failed.")
		c.String(http.StatusBadRequest, "Webhook signature verification failed")
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			c.String(http.StatusBadRequest, "Error parsing session data")
			return
		}

		sub := models.Subscription{
			UserID:         0,
			StripeCustomer: session.CustomerEmail,
			StripeSession:  session.ID,
			Plan:           string(session.Currency), // se estiver usando metadata
			Status:         "active",
		}

		if err := pkg.SaveSubscription(sub); err != nil {
			log.Println("❌ Erro ao salvar assinatura:", err)
			c.String(http.StatusInternalServerError, "Erro ao salvar")
			return
		}
		log.Println("✅ Sessão de assinatura salva com sucesso")

	default:
		log.Println("🔔 Evento recebido:", event.Type)
	}

	c.Status(http.StatusOK)
}
