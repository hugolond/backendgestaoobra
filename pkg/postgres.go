package pkg

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Obra struct {
	ID             string // Usado apenas se quiser armazenar o retorno
	Nome           string
	Endereco       string
	Bairro         string
	Area           string
	Tipo           int
	Casagerminada  bool
	Status         bool
	DataInicioObra string // ou time.Time, dependendo da necessidade
	DataFinalObra  string // idem
}

func OpenConn() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}
	var localhost = os.Getenv("LOCALHOST_POSTGRES")
	var port = os.Getenv("PORT_POSTGRES")
	var user = os.Getenv("USER_POSTGRES")
	var password = os.Getenv("PASSWORD_POSTGRES")
	var dbname = os.Getenv("DBNAME_POSTGRES")

	db, err := sql.Open("postgres", "host="+localhost+" port="+port+" user="+user+" password="+password+" dbname="+dbname+" sslmode=require")
	//fmt.Println("host=" + localhost + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = db.Ping()
	return db, err
}

func InsertLog(data string, app string, keyValue string, keyName string, user string, mensagem string, erro string) {
	conn, err := OpenConn()
	if err != nil {
		return
	}
	defer conn.Close()

	sqlStatement :=
		`INSERT INTO sistema.logs ("createdDate" ,app,"keyValue","keyName","user",mensage,error)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`

	conn.QueryRow(sqlStatement, data, app, keyValue, keyName, user, mensagem, erro)
	if err != nil {
		panic(err)
	}
}

func InsertCarrinho(data string, document string, phone string, documentType string, checkouttag string, corporateDocument string, rclastcart string, rclastcartvalue string, rclastsession string, rclastsessiondate time.Time, email string) {
	conn, err := OpenConn()
	if err != nil {
		return
	}
	defer conn.Close()

	sqlStatement :=
		`INSERT INTO varejo.carrinhoabandonado ("createdDate","document",phone,"documentType",checkouttag,"corporateDocument",rclastcart,rclastcartvalue,rclastsession,rclastsessiondate,email)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

	conn.QueryRow(sqlStatement, data, document, phone, documentType, checkouttag, corporateDocument, rclastcart, rclastcartvalue, rclastsession, rclastsessiondate, email)
	if err != nil {
		panic(err)
	}
}

func InsertObra(obra Obra) (string, error) {
	conn, err := OpenConn()
	if err != nil {
		return "", fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `
		INSERT INTO obra.cadastroobra (
			"nome", "endereco", "bairro", "area", "tipo", "casagerminada", "status", "data_inicio_obra", "data_final_obra"
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING idObra`

	var idObra string
	err = conn.QueryRow(sqlStatement,
		obra.Nome,
		obra.Endereco,
		obra.Bairro,
		obra.Area,
		obra.Tipo,
		obra.Casagerminada,
		obra.Status,
		obra.DataInicioObra,
		obra.DataFinalObra,
	).Scan(&idObra)

	if err != nil {
		return "", fmt.Errorf("erro ao inserir obra: %w", err)
	}

	return idObra, nil
}

func GetAllObra() ([]Obra, error) {
	conn, err := OpenConn()
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `
		SELECT idObra, nome, endereco, bairro, area, tipo, casagerminada, status, data_inicio_obra, data_final_obra
		FROM obra.cadastroobra
		ORDER BY data_inicio_obra DESC`

	rows, err := conn.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	var obras []Obra
	for rows.Next() {
		var u Obra
		err := rows.Scan(
			&u.ID,
			&u.Nome,
			&u.Endereco,
			&u.Bairro,
			&u.Area,
			&u.Tipo,
			&u.Casagerminada,
			&u.Status,
			&u.DataInicioObra,
			&u.DataFinalObra,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler linha: %w", err)
		}
		obras = append(obras, u)
	}

	return obras, nil
}

func InsertPagamento(idObra string, dataPagamento string, detalhe string, categoria string, valorStr string, observacao string) {
	conn, err := OpenConn()
	if err != nil {
		return
	}
	defer conn.Close()

	// Converte valor string "123,22" para float64
	valorStr = strings.Replace(valorStr, ",", ".", 1)
	valor, err := strconv.ParseFloat(valorStr, 64)
	if err != nil {
		panic(fmt.Sprintf("Erro ao converter valor: %v", err))
	}

	sqlStatement := `
		INSERT INTO obra.pagamento ("idObra", "data_do_pagamento", "detalhe", "categoria", "valor", "observacao")
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = conn.Exec(sqlStatement, idObra, dataPagamento, detalhe, categoria, valor, observacao)
	if err != nil {
		panic(fmt.Sprintf("Erro ao inserir pagamento: %v", err))
	}
}
