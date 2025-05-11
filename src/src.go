package src

import (
	"backendgestaoobra/pkg"

	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	AccountID    string `db:"account_id"`
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
	accountID := c.GetString("account_id")
	if accountID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Conta não identificada"})
		return
	}
	dados, err := pkg.SelectObraPagamentoJoin(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dados)
}

func CadastraObra(c *gin.Context) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	accountID := c.GetString("account_id")
	if accountID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Conta não identificada"})
		return
	}
	userID := c.GetString("id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não identificada"})
		return
	}
	userName := c.GetString("username")
	if userName == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não identificada"})
		return
	}

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
	pkg.InsertObra(obra2, accountID, userID, userName)
	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "OBRA", obra.Nome, "nome", "backendgestaoobra", "Obra registrada com sucesso!", "")
	c.JSON(http.StatusOK, gin.H{"message": "Obra cadastrada com sucesso!"})
}

func GetObraByID(c *gin.Context) {
	idObra := c.Param("idobra")

	if idObra == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'idobra' é obrigatório"})
		return
	}
	accountID := c.GetString("account_id")
	if accountID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Conta não identificada"})
		return
	}

	obra, err := pkg.GetObraByID(idObra, accountID)
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
	accountID := c.GetString("account_id")
	if accountID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Conta não identificada"})
		return
	}
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | CA - Consulta lista de obra")

	dados, err := pkg.GetAllObra(accountID)
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
	accountID := c.GetString("account_id")
	if accountID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Conta não identificada"})
		return
	}
	userID := c.GetString("id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não identificada"})
		return
	}
	userName := c.GetString("username")
	if userName == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não identificada"})
		return
	}

	pagamento := pkg.Pagamento{}
	err := json.Unmarshal(buf.Bytes(), &pagamento)
	if err != nil {
		log.Println("Erro ao decodificar pagamento:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao processar pagamento"})
		return
	}

	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006-01-02 - 15:04:05") + " | PG - Insert pagamento da obra: " + pagamento.IDObra)

	err = pkg.InsertPagamentoStruct(pagamento, accountID, userID, userName)
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
	accountID := c.GetString("account_id")
	if accountID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Conta não identificada"})
		return
	}

	pagamentos, err := pkg.GetAllPagamentoByObra(idObra, accountID)
	if err != nil {
		log.Println("Erro ao buscar pagamentos:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar pagamentos"})
		return
	}

	c.JSON(http.StatusOK, pagamentos)
}

func AtualizaObra(c *gin.Context) {
	accountID := c.GetString("account_id")
	if accountID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Conta não identificada"})
		return
	}
	userID := c.GetString("id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não identificada"})
		return
	}
	userName := c.GetString("username")
	if userName == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não identificada"})
		return
	}
	var obra pkg.Obra
	if err := c.ShouldBindJSON(&obra); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	err := pkg.UpdateObra(obra, accountID, userID, userName)
	if err != nil {
		log.Println("Erro ao atualizar obra:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar obra"})
		return
	}

	pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "OBRA", obra.ID, "idObra", "backendgestaoobra", "Obra atualizada com sucesso!", "")
	c.JSON(http.StatusOK, gin.H{"message": "Obra atualizada com sucesso!"})
}

func AtualizaPagamento(c *gin.Context) {
	accountID := c.GetString("account_id")
	if accountID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Conta não identificada"})
		return
	}
	userID := c.GetString("id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não identificada"})
		return
	}
	userName := c.GetString("username")
	if userName == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não identificada"})
		return
	}
	var pagamento pkg.Pagamento
	if err := c.ShouldBindJSON(&pagamento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	err := pkg.UpdatePagamento(pagamento, accountID, userID, userName)
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
	accountID := c.GetString("account_id")
	if accountID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Conta não identificada"})
		return
	}
	err := pkg.DeletePagamento(id, accountID)
	if err != nil {
		log.Println("Erro ao excluir pagamento:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao excluir pagamento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pagamento excluído com sucesso"})
}
