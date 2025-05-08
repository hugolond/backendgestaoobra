package main

import (
	"backendgestaoobra/src"
	"backendgestaoobra/src/middleware"
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

	router.POST("/login", src.LoginHandler)

	router.GET("/healthz", src.Healthz)

	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		//Obra
		protected.POST("/api/obra/v1/sendnewobra", src.CadastraObra) //GE 1
		protected.GET("/api/obra/v1/:idobra", src.GetObraByID)       //GE 1
		protected.GET("/api/obra/v1/listallobra", src.ListObra)      //GE 1
		protected.PUT("/api/obra/v1/update", src.AtualizaObra)       //GE 1

		//Pagamento
		protected.POST("/api/payment/v1/sendnewpayment", src.CadastraPagamento) //GE 1
		protected.GET("/api/payment/v1/listpayment", src.ListPagamentoPorObra)  //GE 1
		protected.PUT("/api/payment/v1/update", src.AtualizaPagamento)          // PG 3
		protected.DELETE("/api/payment/v1/delete", src.DeletePagamento)

	}
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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
