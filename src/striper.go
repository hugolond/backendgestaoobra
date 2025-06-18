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

		customer, err := pkg.GetStripeCustomer(subscription.Customer.ID)
		if err != nil {
			log.Fatal("Erro ao consultar cliente:", err)
		}

		account, err := StartAtivacao(customer.Name, customer.Email, subscription.Items.Data[0].Plan.Product.ID)
		if err != nil {
			log.Fatal("Erro na ativação:", err)
		}
		log.Println("Conta ativada:", account.ID)

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			log.Printf("Erro ao decodificar Subscription: %v", err)
			c.String(400, "Erro de parsing")
			return
		}
		log.Printf("✅ Cancelamento sucedido do cliente: %s", subscription.Customer.ID)
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
			log.Println("❌ Erro ao salvar cancelamento:", err)
			c.String(http.StatusInternalServerError, "Erro ao salvar")
			return
		}
		log.Println("✅ Cancelamento salvo com sucesso:", sub)

	default:
		log.Printf("🔔 Evento não tratado: %s", event.Type)
	}

	c.Status(200)
}
