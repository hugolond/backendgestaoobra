package pkg

import (
	"log"
	"os"
	"strings"
	"time"
)

func GravaLog(linha string) error {
	dia := strings.Split(time.Now().String(), " ")
	// Cria o arquivo de texto

	logFile, err := os.OpenFile(dia[0]+".txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Erro ao abrir o arquivo de log:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println(linha)
	return nil
}

/*func AplicaMascara(text string) (retorno string) {
	antes, depois, found := strings.Cut(text,"@")
	if (found == true){
		// regra email
	}else{

	}

	return "0"
}
*/
