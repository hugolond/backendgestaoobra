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
		log.Println("❌ Erro ao ler corpo:", err)
		c.String(http.StatusRequestEntityTooLarge, "Request too large")
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if endpointSecret == "" {
		log.Println("❌ STRIPE_WEBHOOK_SECRET não definido")
		c.String(http.StatusInternalServerError, "Webhook secret ausente")
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		log.Println("⚠️  Webhook signature verification failed:", err)
		c.String(http.StatusBadRequest, "Assinatura inválida")
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Println("❌ Erro ao decodificar session:", err)
			c.String(http.StatusBadRequest, "Erro parsing")
			return
		}

		sub := models.Subscription{
			UserID:         0,
			StripeCustomer: session.Customer.ID, // <- Melhor que CustomerEmail
			StripeSession:  session.ID,
			Plan:           session.Metadata["plan"], // Requer que você envie metadata no checkout
			Status:         "active",
		}

		if err := pkg.SaveSubscription(sub); err != nil {
			log.Println("❌ Erro ao salvar assinatura:", err)
			c.String(http.StatusInternalServerError, "Erro ao salvar")
			return
		}
		log.Println("✅ Assinatura salva com sucesso:", sub)

	default:
		log.Println("📦 Evento ignorado:", event.Type)
	}

	c.Status(http.StatusOK)
}
