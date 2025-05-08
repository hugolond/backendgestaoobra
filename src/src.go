package src

import (
	"backendgestaoobra/config"
	"backendgestaoobra/pkg"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID           string `db:"id"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	Password     string `db:"password"`
	Active       bool   `db:"active"`
	Roles        string `db:"roles"`
	Departament  string `db:"departament"`
	EmailManager string `db:"emailmanager"`
}

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
	ID             string `json:"id,omitempty"`
	Nome           string `json:"nome"`
	Endereco       string `json:"endereco"`
	Bairro         string `json:"bairro"`
	Area           string `json:"area"`
	Tipo           int    `json:"tipo"`
	Casagerminada  bool   `json:"casagerminada"`
	Status         bool   `json:"status"`
	DataInicioObra string `json:"datainicioobra"`
	DataFinalObra  string `json:"datafinalobra"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

type Pagamento struct {
	ID            int     `json:"id"`
	IDObra        string  `json:"idobra"`
	DataPagamento string  `json:"datapagamento"`
	Detalhe       string  `json:"detalhe"`
	Categoria     string  `json:"categoria"`
	Valor         float64 `json:"valor"`
	Observacao    string  `json:"observacao"`
	CreatedAt     string  `json:"created_at,omitempty"`
	UpdatedAt     string  `json:"updated_at,omitempty"`
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

type ObraPagamento struct {
	IDObra        string  `json:"idobra"`
	Nome          string  `json:"nome"`
	DataPagamento string  `json:"datapagamento"`
	Valor         float64 `json:"valor"`
	Categoria     string  `json:"categoria"`
}

func Healthz(c *gin.Context) {
	post := gin.H{
		"message": "Ok",
	}
	c.JSON(http.StatusOK, post)

}

// GET /api/dashboard/obra-pagamento
func GetObraPagamentoUnificado(c *gin.Context) {
	dados, err := pkg.SelectObraPagamentoJoin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dados)
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

func GetObraByID(c *gin.Context) {
	idObra := c.Param("idobra")
	if idObra == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'idobra' é obrigatório"})
		return
	}

	obra, err := pkg.GetObraByID(idObra)
	if err != nil {
		log.Println("Erro ao consultar obra:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar obra"})
		return
	}

	if obra.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "Obra não encontrada"})
		return
	}

	c.JSON(http.StatusOK, obra)
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

func CadastraPagamento(c *gin.Context) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	pagamento := pkg.Pagamento{}
	err := json.Unmarshal(buf.Bytes(), &pagamento)
	if err != nil {
		log.Println("Erro ao decodificar pagamento:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao processar pagamento"})
		return
	}

	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006-01-02 - 15:04:05") + " | PG - Insert pagamento da obra: " + pagamento.IDObra)

	err = pkg.InsertPagamentoStruct(pagamento)
	if err != nil {
		log.Println("Erro ao inserir pagamento:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao registrar pagamento"})
		return
	}

	pkg.InsertLog(currentTime.Format("2006-01-02 15:04:05"), "PAGAMENTO", pagamento.IDObra, "idObra", "backendgestaoobra", "Pagamento registrado com sucesso!", "")
	c.JSON(http.StatusOK, gin.H{"message": "Pagamento cadastrado com sucesso!"})
}

func ListPagamentoPorObra(c *gin.Context) {
	idObra := c.Query("idobra")
	if idObra == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'idobra' é obrigatório"})
		return
	}

	pagamentos, err := pkg.GetAllPagamentoByObra(idObra)
	if err != nil {
		log.Println("Erro ao buscar pagamentos:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar pagamentos"})
		return
	}

	c.JSON(http.StatusOK, pagamentos)
}

func AtualizaObra(c *gin.Context) {
	var obra pkg.Obra
	if err := c.ShouldBindJSON(&obra); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	err := pkg.UpdateObra(obra)
	if err != nil {
		log.Println("Erro ao atualizar obra:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar obra"})
		return
	}

	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "OBRA", obra.ID, "idObra", "backendgestaoobra", "Obra atualizada com sucesso!", "")
	c.JSON(http.StatusOK, gin.H{"message": "Obra atualizada com sucesso!"})
}

func AtualizaPagamento(c *gin.Context) {
	var pagamento pkg.Pagamento
	if err := c.ShouldBindJSON(&pagamento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	err := pkg.UpdatePagamento(pagamento)
	if err != nil {
		log.Println("Erro ao atualizar pagamento:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar pagamento"})
		return
	}

	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "PAGAMENTO", strconv.Itoa(pagamento.ID), "idPagamento", "backendgestaoobra", "Pagamento atualizado com sucesso!", "")
	c.JSON(http.StatusOK, gin.H{"message": "Pagamento atualizado com sucesso!"})
}

func DeletePagamento(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID do pagamento não informado"})
		return
	}
	err := pkg.DeletePagamento(id)
	if err != nil {
		log.Println("Erro ao excluir pagamento:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao excluir pagamento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pagamento excluído com sucesso"})
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
