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
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type Subscription struct {
	ID             int64
	UserID         int64
	StripeCustomer string
	StripeSession  string
	Plan           string
	Status         string
}

/*
	func StripeWebhookHandler(c *gin.Context) {
		const MaxBodyBytes = int64(65536)
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

		payload, err := ioutil.ReadAll(c.Request.Body)
		//payload, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println("❌ Erro ao ler corpo:", err)
			c.String(http.StatusRequestEntityTooLarge, "Request too large")
			return
		}

		event := stripe.Event{}

		if err := json.Unmarshal(payload, &event); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Webhook error while parsing basic request. %v\n", err.Error())
			c.Status(http.StatusBadRequest)
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
*/
func HandleWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("❌ Erro ao ler o corpo da requisição: %v\n", err)
		c.String(503, "Erro ao ler o corpo")
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	signatureHeader := c.GetHeader("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		log.Printf("⚠️  Falha na verificação da assinatura: %v\n", err)
		c.String(400, "Assinatura inválida")
		return
	}

	switch event.Type {

	case "customer.subscription.created":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			log.Printf("Erro ao decodificar Subscription: %v", err)
			c.String(400, "Erro de parsing")
			return
		}
		log.Printf("✅ Pagamento bem-sucedido do cliente: %s", subscription.Customer.ID)
		sub := models.Subscription{
			UserID:             0, // você pode atualizar se tiver user logado
			StripeCustomer:     subscription.Customer.ID,
			StripeSubscription: subscription.ID,
			StripePriceID:      subscription.Items.Data[0].Plan.ID,
			StripeProductID:    subscription.Items.Data[0].Plan.Product.ID,
			StripePlanAmount:   subscription.Items.Data[0].Plan.Amount,
			Currency:           string(subscription.Currency),
			Interval:           string(subscription.Items.Data[0].Plan.Interval),
			IntervalCount:      int64(subscription.Items.Data[0].Plan.IntervalCount),
			Status:             string(subscription.Status),
		}

		if err := pkg.SaveSubscription(sub); err != nil {
			log.Println("❌ Erro ao salvar assinatura:", err)
			c.String(http.StatusInternalServerError, "Erro ao salvar")
			return
		}
		log.Println("✅ Assinatura salva com sucesso:", sub)

	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			log.Printf("Erro ao decodificar PaymentIntent: %v", err)
			c.String(400, "Erro de parsing")
			return
		}
		log.Printf("✅ Pagamento bem-sucedido no valor de %d", paymentIntent.Amount)

	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		if err := json.Unmarshal(event.Data.Raw, &paymentMethod); err != nil {
			log.Printf("Erro ao decodificar PaymentMethod: %v", err)
			c.String(400, "Erro de parsing")
			return
		}
		log.Printf("✅ Método de pagamento anexado: %s", paymentMethod.ID)

	default:
		log.Printf("🔔 Evento não tratado: %s", event.Type)
	}

	c.Status(200)
}
