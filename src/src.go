package src

import (
	"backendgestaoobra/config"
	"backendgestaoobra/pkg"

	"backendgestaoobra/queue"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type jwtModel struct {
	Username  string
	CreatedAt time.Time
}

type BodyDelivered struct {
	OrderId string
	Status  string
}

type ConsultaSaldoPickup struct {
	Pickup string
	Sku    string
}

type AuthReceived struct {
	Username string
	Password string
}
type TokenBody struct {
	Token string `"token"`
	User  string `"user"`
}
type Frequencia struct {
	IdLoja        string
	DepartureDate string
	ArrivalDate   string
	OperationTime string
}

type Obra struct {
	ID             string // Usado apenas se quiser armazenar o retorno
	Nome           string `json:"nome"`
	Endereco       string `json:"endereco"`
	Bairro         string `json:"bairro"`
	Area           string `json:"area"`
	Tipo           int    `json:"tipo"`
	Casagerminada  bool   `json:"casagerminada"`
	Status         bool   `json:"status"`
	DataInicioObra string `json:"datainicioobra"` // ou time.Time, dependendo da necessidade
	DataFinalObra  string `json:"datafinalobra"`  // idem
}

type RcCarrinho struct {
	Document     string `json:"document"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	DocumentType string `json:"documentType"`
	Checkouttag  struct {
		DisplayValue string `json:"DisplayValue"`
		Scores       struct {
			DadosPessoais []struct {
				Point float64   `json:"Point"`
				Date  time.Time `json:"Date"`
				Until time.Time `json:"Until"`
			} `json:"DadosPessoais"`
			Endereco []struct {
				Point float64   `json:"Point"`
				Date  time.Time `json:"Date"`
				Until time.Time `json:"Until"`
			} `json:"Endereco"`
			FormaPagamento []struct {
				Point float64   `json:"Point"`
				Date  time.Time `json:"Date"`
				Until time.Time `json:"Until"`
			} `json:"FormaPagamento"`
		} `json:"Scores"`
	} `json:"checkouttag"`
	CorporateDocument string    `json:"corporateDocument"`
	Rclastcart        string    `json:"rclastcart"`
	Rclastcartvalue   string    `json:"rclastcartvalue"`
	Rclastsession     string    `json:"rclastsession"`
	Rclastsessiondate time.Time `json:"rclastsessiondate"`
}

func Healthz(c *gin.Context) {
	post := gin.H{
		"message": "Ok",
	}
	c.JSON(http.StatusOK, post)

}

const maxLoja = 5     // quantidade de divisões da lista de lojas ativas
const timeReload = 60 // em minutos

func RetryPdv(c *gin.Context) {
	cfg := &config.Config{}
	config.New(cfg)
	var pedidoPdv pkg.FranchiseOrder
	orderid := c.Query("orderid")
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | C1 - Inquiry received - OrderId: " + orderid)

	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Franchise-order - OrderId", orderid, "OrderId", "web-server-pnb", "Inquiry received", "")
	pedidoPdv, err := pkg.ConsultaOrderId(orderid)
	if err != nil {
		log.Fatalf("Erro ao consultar pedido")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar pedido"})
		return
	}
	order, err := pkg.GetOrder(cfg.Account, cfg.KeyVtex, cfg.TokenVtex, orderid)
	if err != nil {
		log.Fatalf("Erro ao consultar pedido")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar pedido"})
		return
	}
	if len(order.Status) > 0 {
		var requestPdv pkg.ResponsePdv

		reg, err := regexp.Compile("[^0-9]")
		if err != nil {
			log.Fatalf("Erro ao consultar pedido")
			c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar pedido"})
			return
		}
		loja := reg.ReplaceAllString(order.Sellers[0].Name, "")
		numerber, err := strconv.Atoi(loja)

		if (pedidoPdv.InvoiceKey != nil) && (pedidoPdv.InvoiceKey != "") {
			requestPdv.OrderDate = time.Now()
			requestPdv.Customer.Name = order.ClientProfileData.FirstName + " " + order.ClientProfileData.LastName
			requestPdv.StoreNumber = loja
			requestPdv.Customer.Document = order.ClientProfileData.Document
			requestPdv.TotalValue = order.Value
			requestPdv.OrderPayments = nil
			requestPdv.Message = "Pedido já faturado!"
			c.JSON(http.StatusOK, requestPdv)
			return
		}

		if pedidoPdv.Message == "" {
			// prepara dados para envio do pedido para o PDV
			var preparaPdv pkg.Pdv
			preparaPdv.StoreID = numerber
			preparaPdv.Origem = "CONTA-FRANQUIA"
			preparaPdv.Document = pedidoPdv.CustomerCpf
			preparaPdv.ChapaColaborador = "0800"
			for _, item := range pedidoPdv.Items {
				preparaPdv.Items = append(preparaPdv.Items, struct {
					CodigoBarcode string "json:\"codigoBarcode\""
					Quantity      int    "json:\"quantity\""
				}{CodigoBarcode: item.AutomationCodeBeep, Quantity: item.QuantityBeep})

			}
			preparaPdv.Payment.VtexOrder = pedidoPdv.OrderID
			pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Franchise-order - OrderId", orderid, "OrderId", "web-server-pnb", "Order send", "")

			requestPdv, err := pkg.EnviaPdv(preparaPdv)
			if err != nil {
				log.Fatalf("Erro ao consultar pedido")
				c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar pedido"})
				return
			}
			pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Franchise-order - OrderId", strconv.Itoa(requestPdv.OrderNumber), "OrderNumber", "web-server-pnb", "Order Pdv received", "")
			currentTime = time.Now()
			fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | C1 - Franchise Order Complete  - OrderNumber: " + strconv.Itoa(requestPdv.OrderNumber))

			c.JSON(http.StatusOK, requestPdv)
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": pedidoPdv.Message})
}

func InvoicedService(c *gin.Context) {
	cfg := &config.Config{}
	config.New(cfg)
	orderid := c.Query("orderid")
	tipo := c.Query("type")
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | G1 - Resquest Invoiced - OrderId: " + orderid)
	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "GERF - OrderId", orderid, "OrderId", "web-server-pnb", "Resquest Invoiced", "")

	order, err := pkg.GetOrder(cfg.Account, cfg.KeyVtex, cfg.TokenVtex, orderid)
	if err != nil {
		log.Fatalf("Erro ao consultar pedido")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar pedido"})
		return
	}
	var nomeConta string
	if len(order.Status) > 0 {
		if order.Sellers[0].ID == "1" {
			nomeConta = cfg.Account
		} else {
			nomeConta = order.Sellers[0].ID
			order, err = pkg.GetOrder(nomeConta, cfg.KeyVtex, cfg.TokenVtex, order.SellerOrderID)
			if err != nil {
				log.Fatalf("Erro ao consultar pedido")
				c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar pedido"})
				return
			}
		}

		var ListaItems []pkg.InvoicedItem
		if order.Status != "invoiced" {
			if order.Status == "ready-for-handling" {
				pkg.RegistraManuseio(nomeConta, cfg.KeyVtex, cfg.TokenVtex, order.OrderID)
			}
			for posicao, itens := range order.Items {
				var itemEncontrado = false
				for _, Pacotes := range order.PackageAttachment.Packages {
					for _, itensPacotes := range Pacotes.Items {
						if posicao == itensPacotes.ItemIndex {
							itemEncontrado = true
						}
					}
				}
				if itemEncontrado == false {
					sku, err := strconv.Atoi(itens.ID)
					if err != nil {
					}
					if tipo == "service" {
						if strings.Contains(itens.RefID, "#") {
							ListaItems = append(ListaItems, pkg.InvoicedItem{Sku: sku, Quantity: itens.Quantity, Price: itens.Price, Description: nil, UnitMultiplier: 0})
						}
					}
					if tipo == "all" {
						ListaItems = append(ListaItems, pkg.InvoicedItem{Sku: sku, Quantity: itens.Quantity, Price: itens.Price, Description: nil, UnitMultiplier: 0})
					}
				}
			}
			var valueShipping int
			for _, totals := range order.Totals {
				if totals.ID == "Shipping" {
					valueShipping = totals.Value
				}
			}

			if len(ListaItems) > 0 {
				recibo, err := pkg.FaturaProdutosPedido(nomeConta, cfg.KeyVtex, cfg.TokenVtex, order.OrderID, order.OrderID+"-1", valueShipping, ListaItems)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao faturar pedido"})
					return
				}

				currentTime = time.Now()
				pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Invoiced - OrderId", orderid, "OrderId", "web-server-pnb", "Order invoiced Receipt: "+recibo.Receipt, "")
				fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | G1 - Invoiced Service  - OrderNumber: " + orderid)
				c.JSON(http.StatusOK, gin.H{"message": "Pedido faturado com sucesso!", "Receipt": recibo.Receipt})
			}
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Pedido já faturado!"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Pedido não encontrado!"})
	}
}

func Delivered(c *gin.Context) {
	cfg := &config.Config{}
	config.New(cfg)
	orderid := c.Query("orderid")
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | G2 - Resquest Delivered - OrderId: " + orderid)
	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Delivered - OrderId", orderid, "OrderId", "web-server-pnb", "Resquest Delivered", "")

	order, err := pkg.GetOrder(cfg.Account, cfg.KeyVtex, cfg.TokenVtex, orderid)
	if err != nil {
		log.Fatalf("Erro ao consultar pedido")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar pedido"})
		return
	}
	var nomeConta string
	if len(order.Status) > 0 {
		if order.Sellers[0].ID == "1" {
			nomeConta = cfg.Account
		} else {
			nomeConta = order.Sellers[0].ID
			order, err = pkg.GetOrder(nomeConta, cfg.KeyVtex, cfg.TokenVtex, order.SellerOrderID)
			if err != nil {
				log.Fatalf("Erro ao consultar pedido")
				c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar pedido"})
				return
			}
		}
		if order.Status == "invoiced" {
			var recibos []string
			for _, Pacotes := range order.PackageAttachment.Packages {
				recibo, err := pkg.MarcarPedidoEntregue(nomeConta, cfg.KeyVtex, cfg.TokenVtex, order.OrderID, Pacotes.InvoiceNumber)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao marcar pedido entregue"})
					return
				}
				currentTime = time.Now()
				pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Delivered - OrderId", orderid, "OrderId", "web-server-pnb", "Order Delivered Receipt: "+recibo.Receipt, "")
				fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | G2 - Delivered Service  - OrderNumber: " + orderid)
				recibos = append(recibos, recibo.Receipt)
			}
			c.JSON(http.StatusOK, gin.H{"message": "Pedido marcado como entregue com sucesso!!", "Receipt": recibos})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Pedido não está faturado!"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Pedido não encontrado!"})
	}
}

func CliqueRetireSearch(c *gin.Context) {
	cfg := &config.Config{}
	config.New(cfg)
	var pedidoCliqueRetire pkg.CliqueRetire
	termo := c.Query("termo")
	campo := c.Query("campo")
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | CR1 - Search Termo: " + termo + " campo: " + campo)

	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Clique-Retire - Termo", campo, "OrderId", "web-server-pnb", "Search Received", "")
	pedidoCliqueRetire, err := pkg.GetCliqueRetire(cfg.Account, cfg.KeyVtex, cfg.TokenVtex, campo, termo)
	if err != nil {
		log.Fatalf("Erro ao consultar pedido")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar pedido!"})
		return
	}

	if len(pedidoCliqueRetire) > 0 {
		if pedidoCliqueRetire[0].IsShippingCompany == true {
			order, err := pkg.GetOrder(cfg.Account, cfg.KeyVtex, cfg.TokenVtex, campo)
			if err != nil {
				log.Fatalf("Erro ao consultar pedido")
				c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar pedido!"})
				return
			}
			if len(order.PackageAttachment.Packages) > 0 {
				pedidoCliqueRetire[0].UrlRastreio = order.PackageAttachment.Packages[0].TrackingURL
				pedidoCliqueRetire[0].CourierName = order.PackageAttachment.Packages[0].Courier
				if len(order.PackageAttachment.Packages[0].CourierStatus.Status) > 0 {
					pedidoCliqueRetire[0].FinishedCourier = order.PackageAttachment.Packages[0].CourierStatus.Finished
					pedidoCliqueRetire[0].DeliveredDate = order.PackageAttachment.Packages[0].CourierStatus.DeliveredDate
				}
			}
		}
		c.JSON(http.StatusOK, pedidoCliqueRetire)
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar pedido!"})
		return
	}
}
func CliqueRetireUpdateBase(c *gin.Context) {

	body := csv.NewReader(c.Request.Body)
	body.Comma = ';'
	body.LazyQuotes = false
	body.FieldsPerRecord = 4
	var frequenciaslojas []Frequencia

	for {
		linha, err := body.Read()

		if err == io.EOF {
			break
		}
		var frequencia Frequencia
		for i, date := range linha {

			if i == 0 {
				frequencia.IdLoja = date
			}
			if i == 1 {
				frequencia.DepartureDate = date
			}
			if i == 2 {
				frequencia.ArrivalDate = date
			}
			if i == 3 {
				frequencia.OperationTime = date
			}
		}
		if _, err := strconv.Atoi(frequencia.IdLoja); err == nil {
			frequenciaslojas = append(frequenciaslojas, frequencia)
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Total de linhas: " + strconv.Itoa(len(frequenciaslojas))})
}
func CliqueRetireValidatePickup(c *gin.Context) {

}
func RequestDelivered(c *gin.Context) {
	body := new(BodyDelivered)
	if err := c.Bind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao receber pedido!"})
		return
	}

	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | C2 - POST Delivered - Order Id: " + body.OrderId + " Status: " + body.Status)

	//connection := queue.Connect()
	data, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err.Error())
	}
	queue.Notify(data, "delivered_ex", "")
	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Service", body.OrderId, "OrderId", "web-server-pnb", "Payload Recebido com sucesso!", "")
	if err != nil {
		log.Fatalf("Erro ao gravar log")
	}
	c.JSON(200, body)
}

func RcVarejo(c *gin.Context) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	carrinho := RcCarrinho{}
	err := json.Unmarshal(buf.Bytes(), &carrinho)
	if err != nil {
		log.Println(err.Error())
		return
	}
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | CA - Insert dados carrinho document: " + carrinho.Document)
	pkg.InsertCarrinho(time.Now().Format("2006-01-02 15:04:05"), carrinho.Document, carrinho.Phone, carrinho.DocumentType, carrinho.CorporateDocument, carrinho.Checkouttag.DisplayValue, carrinho.Rclastcart, carrinho.Rclastcartvalue, carrinho.Rclastsession, carrinho.Rclastsessiondate, carrinho.Email)
	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "RCVTEX", carrinho.Document, "document", "web-server-pnb", "Carrinho Registrado com sucesso!", "")
	c.JSON(http.StatusOK, gin.H{"message": "Carrinho gravado com sucesso!"})
}

func CadastraObra(c *gin.Context) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	obra := Obra{}
	err := json.Unmarshal(buf.Bytes(), &obra)
	if err != nil {
		log.Println(err.Error())
		return
	}
	obra2 := pkg.Obra{}
	obra2.Nome = obra.Nome
	obra2.Endereco = obra.Endereco
	obra2.Bairro = obra.Bairro
	obra2.Area = obra.Area
	obra2.Tipo = obra.Tipo
	obra2.Casagerminada = obra.Casagerminada
	obra2.Status = obra.Status
	obra2.DataInicioObra = obra.DataInicioObra
	obra2.DataFinalObra = obra.DataFinalObra

	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | CA - Insert dados obra: " + obra.Nome)
	pkg.InsertObra(obra2)
	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "OBRA", obra.Nome, "nome", "backendgestaoobra", "Obra registrada com sucesso!", "")
	c.JSON(http.StatusOK, gin.H{"message": "Obra cadastrada com sucesso!"})
}

func ListObra(c *gin.Context) {
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | CA - Consulta lista de obra")

	dados, err := pkg.GetAllObra()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar dados da obra"})
		fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | CA - Consulta lista de obra - Error:" + err.Error())
		return
	}

	pkg.InsertLog(
		time.Now().Format("2006-01-02 15:04:05"),
		"OBRA",
		"All",
		"Nome",
		"backendgestaoobra",
		"Consulta realizada com sucesso!",
		"",
	)

	c.JSON(http.StatusOK, dados)
}

func GetSaldo(c *gin.Context) {
	/*body := new(ConsultaSaldoPickup)
	if err := c.Bind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao receber pedido!"})
		return
	}*/

	cfg := &config.Config{}
	config.New(cfg)
	sku := c.Query("sku")
	pickup := c.Query("pickup")

	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | C3 - GET Consulta Saldo Pickup: " + pickup + " Sku: " + sku)

	tokenVtex, err := pkg.ConsultaTokenVtex()
	dados, err := pkg.GetSaldoSkuPickup(pickup, tokenVtex.Token, sku)

	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Conta-franquia", sku, "Sku", "web-server-pnb", "Consulta Recebido com sucesso!", "")
	if err != nil {
		log.Fatalf("Erro ao gravar log")
	}
	if len(dados.Data.ProductHistory.SkuName) > 0 {
		c.JSON(200, dados)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar SKU"})
	}

}

func RequestAuth(c *gin.Context) {
	body := new(AuthReceived)
	if err := c.Bind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao autenticar"})
		return
	}
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | A1 - Auth Received - User: " + body.Username + " Status: " + body.Password)
	//log.Println("Auth Received - User: " + body.Username + " Status: " + body.Password)
	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Auth", body.Username, "UserName", "web-server-pnb", "Username tentative login!", "")
	var user = "hugodias"
	var pass = "2403Alice@"

	if body.Username == user && body.Password == pass {
		token, err := pkg.GenerateToken(body.Username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao autenticar"})
			return
		}
		dados := new(TokenBody)
		dados.Token = token
		dados.User = user
		c.JSON(200, dados)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "Usuário ou senha inválidos!"})
}

func splitList(input []int, parts int) [][]int {
	if parts <= 0 {
		return nil
	}

	size := (len(input) + parts - 1) / parts
	result := make([][]int, parts)

	for i := 0; i < len(input); i += size {
		end := i + size
		if end > len(input) {
			end = len(input)
		}
		result[i/size] = input[i:end]
	}

	return result
}

func SolicitaAtualizacaoSaldo(lista []int) {

	//url := "a"
	url := "https://conta-franquia-backoffice-api.pernambucanas.com.br/api/v4/estoque/sync-vtex-lojas"
	method := "POST"

	payload := strings.NewReader(`{"cod_estabelecimento": [` + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(lista)), ","), "[]") + `],"sc": 4}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Access-Control-Allow-Origin", "no-cors")

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
	fmt.Println(string(body))
}

func SacolaDesconto(c *gin.Context) {
	cfg := &config.Config{}
	config.New(cfg)
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | SC1 - Sacola Desconto 1 - Envia Push Sacola Desconto")
	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Sacola Desconto", "Sacola Desconto", "Sacola Desconto", "web-server-pnb", "Rotina envio de Push Sacola de Desconto", "")

	listaSacola, err := pkg.GetSacolaAbandonada()
	if err != nil {
		//log.Fatalf("Erro ao converter texto")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar sacolas de desconto!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rotina envio Push para sacola enviada com sucesso!"})

	for _, sacola := range listaSacola {
		response, err := pkg.SendPush(sacola)
		if err != nil {
			//log.Fatalf("Erro ao converter texto")
			c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao consultar sacolas de desconto!"})
			return
		} else {
			currentTime := time.Now()
			if len(response.Enviados) > 0 {
				fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | SC1 - Sacola Desconto 1 - Envia Push Sacola Desconto id:" + strconv.Itoa(sacola) + " - Status Enviado: " + strconv.Itoa(response.Enviados[0]))
			} else {
				fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | SC1 - Sacola Desconto 1 - Envia Push Sacola Desconto id:" + strconv.Itoa(sacola) + " - Status Não Enviado: " + strconv.Itoa(response.NaoEnviados[0]))
			}

		}
	}
}

func AtualizaSaldo(c *gin.Context) {
	nowStart := time.Now()
	cfg := &config.Config{}
	config.New(cfg)
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | CRON1 - Cron 1 - Atualiza Saldo Conta Franquia")
	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Cron", "Conta Franquia", "Conta Franquia", "web-server-pnb", "Cron Atualiza Saldo Conta Franquia iniciada", "")

	var min = 0
	var max = 15
	var ListaLojasConsultadas []int
	var count = 1
	for i := 0; i < count; i++ {
		lojasContaFranquia, _ := pkg.GetLojasContaFranquia(cfg.Account, cfg.KeyVtex, cfg.TokenVtex, min, max)
		min = max + 1
		max = max + (max * (i + 1))
		if len(lojasContaFranquia) > 0 {
			for _, loja := range lojasContaFranquia {
				ListaLojasConsultadas = append(ListaLojasConsultadas, loja.IDLoja)
			}
			count = count + 1
		} else {
			count = 0
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rotina iniciada com sucesso! Quant. Lojas: " + strconv.Itoa(len(ListaLojasConsultadas))})

	sort.Slice(ListaLojasConsultadas, func(i, j int) bool {
		return ListaLojasConsultadas[i] < ListaLojasConsultadas[j] // Organiza do maior para o menor
	})
	var novaListaLojasUnique = removeDuplicates(ListaLojasConsultadas)

	result := splitList(novaListaLojasUnique, maxLoja) // lista de x lojas

	for _, part := range result {
		SolicitaAtualizacaoSaldo(part)
		pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Cron", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(part)), ","), "[]"), "Conta Franquia", "web-server-pnb", "Cron Atualiza Saldo Conta Franquia finalizada", "")
		time.Sleep(timeReload * 60 * time.Second) // x minutos
	}
	nowEnd := time.Now()
	nowStart.Sub(nowEnd)
	time := nowStart.Sub(nowEnd)
	c.JSON(http.StatusOK, gin.H{"message": "Rotina finalizada com sucesso! Quant. Lojas: " + strconv.Itoa(len(ListaLojasConsultadas)) + "Time: " + shortDur(time)})
}
func RequestVerifyToken(c *gin.Context) {
	token := c.Query("token")
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | A2 - Verify Received - token: " + token)
	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "VerifyToken", token, "Token", "web-server-pnb", "Token Verify!", "")
	valida := pkg.AuthMiddleware(token)
	if valida == "ok" {
		c.JSON(200, gin.H{"message": "Token verificado com sucesso"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao verificar token"})
}
func AtualizaSaldoLoja(c *gin.Context) {

	cfg := &config.Config{}
	account := c.Query("account")
	config.New(cfg)
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | C3 - Atualiza Saldo Loja - Atualiza Saldo Conta Franquia: " + account)
	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Cron", "Conta Franquia", "Conta Franquia", "web-server-pnb", "Atualiza Saldo Loja Conta Franquia iniciada", "")
	var lista []int

	reg, _ := regexp.Compile("[^0-9]")

	loja := reg.ReplaceAllString(account, "")
	numerber, err := strconv.Atoi(loja)
	if err != nil {
		//log.Fatalf("Erro ao converter texto")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao enviar atualização!"})
		return
	}
	lista = append(lista, numerber)
	SolicitaAtualizacaoSaldo(lista)

	c.JSON(http.StatusOK, gin.H{"message": "Atualização enviada com sucesso! Loja: " + account})
}

func removeDuplicates(nums []int) []int {
	// Mapa para manter controle dos elementos únicos
	seen := make(map[int]bool)
	unique := []int{}

	// Iterar pela lista original
	for _, num := range nums {
		// Verificar se o elemento já foi visto
		if !seen[num] {
			// Se não foi visto, adicionar ao mapa e à lista de elementos únicos
			seen[num] = true
			unique = append(unique, num)
		}
	}

	return unique
}

func ClienteSearch(c *gin.Context) {
	cfg := &config.Config{}
	config.New(cfg)
	termo := c.Query("termo")
	campo := c.Query("campo")
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | CL1 - Search Client: " + termo + " campo: " + campo)

	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Client - Termo", campo, "Client", "web-server-pnb", "Search Client", "")
	clienteGer, err := pkg.GetCliente(cfg.TokenGer, termo, campo)
	if err != nil {
		log.Fatalf("Erro ao pesquisar cliente")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao pesquisar cliente!"})
		return
	}
	if clienteGer.NumberOfElements > 0 {
		c.JSON(http.StatusOK, clienteGer)
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao pesquisar cliente!"})
		return
	}
}

func shortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}
