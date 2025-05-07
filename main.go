package main

import (
	"backendgestaoobra/config"
	"backendgestaoobra/pkg"
	"backendgestaoobra/queue"
	"backendgestaoobra/src"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// ConfigRuntime sets the number of operating system threads.
func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

type BodyDelivered struct {
	OrderId string
	Status  string
}

// StartWorkers start starsWorker by goroutine.
/*func StartWorkers() {
	go statsWorker()
}
*/

// StartGin starts gin web server with setting router.

func main() {
	// Rabbit
	//go consumer()

	gin.SetMode(gin.ReleaseMode)
	ConfigRuntime()
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.GET("/healthz", src.Healthz)

	//Obra
	router.POST("/api/obra/v1/sendnewobra", src.CadastraObra) //GE 1
	router.GET("/api/obra/v1/listObra", src.ListObra)         //GE 1

	// Carrinho Abandonado
	router.POST("/api/carrinho/v1/rcvarejo", src.RcVarejo) // CA1

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}

func consumer() {
	// Env
	cfg := &config.Config{}
	err := config.New(cfg)
	if err != nil {
		fmt.Println("Arquivo '.env' não encontrado")
	}

	in := make(chan []byte)
	conn := queue.Connect()
	ch, err := conn.Channel()
	if err != nil {
		panic(err.Error())
	}
	//fmt.Println("aa")
	queue.StartCosuming(ch, in)
	//fmt.Println("a")

	for payload := range in {
		var resp BodyDelivered
		err := json.Unmarshal(payload, &resp)
		if err != nil {
			fmt.Println("Erro ao converter Json" + err.Error())
		}
		currentTime := time.Now()
		fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | C2 - Consumindo Mensagem Order: " + resp.OrderId)
		pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Consumer", resp.OrderId, "OrderId", "web-server-pnb", "Mensagem consumida com sucesso!", "")

		order, err := pkg.GetOrder(cfg.Account, cfg.KeyVtex, cfg.TokenVtex, resp.OrderId)
		if err != nil {
			fmt.Println(err)
			return
		}
		currentTime = time.Now()
		fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | C2 - Order Status: " + string(order.Status) + " Order: " + resp.OrderId)

		for _, pacote := range order.PackageAttachment.Packages {
			if pacote.TrackingNumber == "" && order.Status == "invoiced" {
				pkg.RegistraEntregue(cfg.Account, cfg.KeyVtex, cfg.TokenVtex, resp.OrderId, string(pacote.InvoiceNumber))
				pkg.InsertLog(time.Now().Format("2006-01-02 15:04:05"), "Consumer", resp.OrderId, "OrderId", "web-server-pnb", "Notificacao entrega enviado OrderId!", "")
			}
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
