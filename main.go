package main

import (
	"backendgestaoobra/src"
	"fmt"
	"log"
	"os"
	"runtime"

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
	router.PUT("/api/obra/v1/update", src.AtualizaObra)

	//Pagamento
	router.POST("/api/payment/v1/sendnewpayment", src.CadastraPagamento) //GE 1
	router.GET("/api/payment/v1/listpayment", src.ListPagamentoPorObra)  //GE 1
	router.PUT("/api/payment/v1/update", src.AtualizaPagamento)          // PG 3

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
