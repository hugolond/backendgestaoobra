package src

import (
	models "backendgestaoobra/model"
	"backendgestaoobra/pkg"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// StartAtivacao cria ou atualiza uma conta com base no email e plano Stripe
func StartAtivacao(nome string, email string, stripeProductID string) (*models.Account, error) {
	// 1. Verifica se já existe account com este email
	existingAccount, err := pkg.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	if existingAccount != nil {
		// 2. Atualiza o stripe_product_id
		err := pkg.UpdateAccountPlan(existingAccount.ID, stripeProductID)
		if err != nil {
			return nil, err
		}
		existingAccount.StripeProductID = stripeProductID
		// Verifica se possui um plano inferior ativo e desativa
		fmt.Println("Consulta plano anterior ")
		planoAnterior, err := pkg.BuscarAssinaturaAtivaAnterior(email, stripeProductID)
		if err != nil {
			log.Println("Erro:", err)
		}
		if planoAnterior != "" {
			fmt.Println("Plano anterior ativo encontrado:", planoAnterior)
			_, err := pkg.CancelSubscription(planoAnterior)
			if err != nil {
				return nil, err
			}
		}
		return existingAccount, nil
	}

	// 3. Cria nova conta
	newAccount := models.Account{
		ID:              uuid.NewString(),
		Nome:            nome,
		Email:           email,
		StripeProductID: stripeProductID,
		Status:          true,
		CreatedAt:       time.Now(),
	}

	// Inserir no banco
	err = pkg.CreateAccount(newAccount)
	if err != nil {
		return nil, err
	}

	return &newAccount, nil
}

func StartDesativacao(nome string, email string, stripeProductID string) (*models.Account, error) {
	// 1. Verifica se já existe account com este email
	existingAccount, err := pkg.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	// verficia o plano atual ativo
	planoAnterior, err := pkg.BuscarAssinaturaAtivaAnterior(email, stripeProductID)
	if err != nil {
		log.Println("Erro:", err)
	}
	if planoAnterior == "" {
		fmt.Println("Não há outro plano ativo anterior")
		err := pkg.UpdateAccountPlan(existingAccount.ID, "free")
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Plano anterior ativo:", planoAnterior)
		err := pkg.UpdateAccountPlan(existingAccount.ID, planoAnterior)
		if err != nil {
			return nil, err
		}
	}

	newAccount, err := pkg.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}
	return newAccount, nil
}
