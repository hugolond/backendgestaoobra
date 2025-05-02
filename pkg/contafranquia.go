package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type FranchiseOrder struct {
	Message     string    `json:"message"`
	Stack       string    `json:"stack"`
	ID          int       `json:"id"`
	CustomerCpf string    `json:"customerCpf"`
	OrderID     string    `json:"orderId"`
	StatusID    int       `json:"statusId"`
	StoreID     any       `json:"storeId"`
	ImageURL    any       `json:"imageUrl"`
	InvoiceKey  any       `json:"invoiceKey"`
	OrderNumber int       `json:"orderNumber"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Items       []struct {
		ID                 int       `json:"id"`
		UniqueID           string    `json:"uniqueId"`
		FranchiseOrderID   int       `json:"franchiseOrderId"`
		StatusID           int       `json:"statusId"`
		AutomationCode     string    `json:"automationCode"`
		AutomationCodeBeep string    `json:"automationCodeBeep"`
		Quantity           int       `json:"quantity"`
		QuantityBeep       int       `json:"quantityBeep"`
		BeepUpdatedAt      time.Time `json:"beepUpdatedAt"`
		CreatedAt          time.Time `json:"createdAt"`
		UpdatedAt          time.Time `json:"updatedAt"`
	} `json:"items"`
}

type Pdv struct {
	StoreID          int    `json:"storeId"`
	Origem           string `json:"origem"`
	Document         string `json:"document"`
	ChapaColaborador string `json:"chapaColaborador"`
	Items            []struct {
		CodigoBarcode string `json:"codigoBarcode"`
		Quantity      int    `json:"quantity"`
	} `json:"items"`
	Payment struct {
		VtexOrder string `json:"vtexOrder"`
	} `json:"payment"`
}

type ResponsePdv struct {
	Message       string `json:"message"`
	CompanyNumber string `json:"companyNumber"`
	StoreNumber   string `json:"storeNumber"`
	OperatorCode  string `json:"operatorCode"`
	OrderNumber   int    `json:"orderNumber"`
	Status        string `json:"status"`
	OriginSystem  string `json:"originSystem"`
	Customer      struct {
		CustomerID   string `json:"customerId"`
		Name         string `json:"name"`
		Document     string `json:"document"`
		DocumentType string `json:"documentType"`
		Address      struct {
			AddressName       string `json:"addressName"`
			AddressNumber     string `json:"addressNumber"`
			AddressComplement string `json:"addressComplement"`
			City              string `json:"city"`
			State             string `json:"state"`
			Country           string `json:"country"`
			ZipCode           string `json:"zipCode"`
			Neighborhood      string `json:"neighborhood"`
			Phone             string `json:"phone"`
		} `json:"address"`
		CustomerType string `json:"customerType"`
		Email        string `json:"email"`
	} `json:"customer"`
	FiscalItems []struct {
		SequenceID                int     `json:"sequenceId"`
		ProductCode               string  `json:"productCode"`
		AutomationCode            string  `json:"automationCode"`
		ProductDescription        string  `json:"productDescription"`
		ProductCompactDescription string  `json:"productCompactDescription"`
		DeliveryOption            string  `json:"deliveryOption"`
		UnitValue                 float64 `json:"unitValue"`
		QuantitySold              int     `json:"quantitySold"`
		TotalValue                float64 `json:"totalValue"`
		UnitQuantity              int     `json:"unitQuantity"`
		SaleUnit                  string  `json:"saleUnit"`
		MaterialOrigin            string  `json:"materialOrigin"`
		CbModality                int     `json:"cbModality"`
		Discounts                 []struct {
			SequenceID int     `json:"sequenceId"`
			Value      float64 `json:"value"`
			Type       string  `json:"type"`
			Source     string  `json:"source"`
		} `json:"discounts"`
		ExternalSequenceID string `json:"externalSequenceId"`
	} `json:"fiscalItems"`
	OrderPayments []struct {
		PaymentSequence        string  `json:"paymentSequence"`
		PaymentCodeMethod      string  `json:"paymentCodeMethod"`
		PaymentCodePlan        string  `json:"paymentCodePlan"`
		PaymentValue           float64 `json:"paymentValue"`
		NumberOfInstallments   int     `json:"numberOfInstallments"`
		MandatoryMethod        bool    `json:"mandatoryMethod"`
		PaymentMethodIndicator int     `json:"paymentMethodIndicator"`
		IntegrationType        int     `json:"integrationType"`
	} `json:"orderPayments"`
	OrderDate  time.Time `json:"orderDate"`
	TotalValue int       `json:"totalValue"`
	Discounts  []struct {
		SequenceID int     `json:"sequenceId"`
		Value      float64 `json:"value"`
		Type       string  `json:"type"`
		Source     string  `json:"source"`
	} `json:"discounts"`
}

func EnviaPdv(pedido Pdv) (resp ResponsePdv, err error) {
	// mock
	//url := "https://f24ccbeb-9471-4c10-9a5b-a81f3f63872c.mock.pstmn.io/api/conta-franquia/v1/pdv"
	url := "https://cms-conta-franquia-api-hml.pernambucanas.com.br/api/conta-franquia/v1/pdv"
	method := "POST"
	client := &http.Client{}
	jsonBytes, err := json.Marshal(pedido)
	pedidoString := string(jsonBytes)
	payload := strings.NewReader(pedidoString)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Erro ao converter Json" + err.Error())
		return
	}
	return
}

func ConsultaOrderId(OrderId string) (resp FranchiseOrder, err error) {
	url := "https://cms-conta-franquia-api.pernambucanas.com.br/api/conta-franquia/v1/franchise-order/" + OrderId
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Erro ao converter Json" + err.Error())
		return
	}
	return
}
