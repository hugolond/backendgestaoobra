package pkg

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ObraPagamento struct {
	IDObra        string  `json:"idobra"`
	Nome          string  `json:"nome"`
	DataPagamento string  `json:"datapagamento"`
	Valor         float64 `json:"valor"`
	Categoria     string  `json:"categoria"`
}

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
	CreatedAt      string
	UpdatedAt      string
}

type Pagamento struct {
	ID            int     `json:"id,omitempty"`
	IDObra        string  `json:"idobra"`
	DataPagamento *string `json:"datapagamento"`
	Detalhe       string  `json:"detalhe"`
	Categoria     string  `json:"categoria"`
	Valor         float64 `json:"valor"`
	Observacao    *string `json:"observacao"`
	CreatedAt     string  `json:"created_at,omitempty"`
	UpdatedAt     string  `json:"updated_at,omitempty"`
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

func SelectObraPagamentoJoin() ([]ObraPagamento, error) {
	conn, err := OpenConn()
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	query := `
		SELECT
			o.idObra,
			o.nome,
			COALESCE(p.data_do_pagamento, '') AS data_do_pagamento,
			COALESCE(p.valor, 0),
			COALESCE(p.categoria, '')
		FROM obra.cadastroobra o
		LEFT JOIN obra.pagamento p ON p.idObra = o.idObra
		ORDER BY o.nome, p.data_do_pagamento DESC
	`

	rows, err := conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	var dados []ObraPagamento
	for rows.Next() {
		var linha ObraPagamento
		err := rows.Scan(&linha.IDObra, &linha.Nome, &linha.DataPagamento, &linha.Valor, &linha.Categoria)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler linha: %w", err)
		}
		dados = append(dados, linha)
	}

	return dados, nil
}

func InsertObra(obra Obra) (string, error) {
	conn, err := OpenConn()
	if err != nil {
		return "", fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `
		INSERT INTO obra.cadastroobra (
			nome, endereco, bairro, area, tipo, casagerminada, status, data_inicio_obra, data_final_obra, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now()
		) RETURNING idObra`

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
		SELECT idObra, nome, endereco, bairro, area, tipo, casagerminada, status, data_inicio_obra, data_final_obra, created_at, updated_at
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
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler linha: %w", err)
		}
		obras = append(obras, u)
	}

	return obras, nil
}

func InsertPagamentoStruct(p Pagamento) error {
	conn, err := OpenConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	sqlStatement := `
		INSERT INTO obra.pagamento (
			idObra, data_do_pagamento, detalhe, categoria, valor, observacao, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, now(), now()
		)`

	_, err = conn.Exec(sqlStatement,
		p.IDObra,
		p.DataPagamento,
		p.Detalhe,
		p.Categoria,
		p.Valor,
		p.Observacao,
	)

	return err
}

func GetAllPagamentoByObra(idObra string) ([]Pagamento, error) {
	conn, err := OpenConn()
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `
		SELECT id, idObra, data_do_pagamento, detalhe, categoria, valor, observacao, created_at, updated_at
		FROM obra.pagamento
		WHERE idObra = $1
		ORDER BY data_do_pagamento DESC`

	rows, err := conn.Query(sqlStatement, idObra)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	var pagamentos []Pagamento
	for rows.Next() {
		var p Pagamento
		err := rows.Scan(
			&p.ID,
			&p.IDObra,
			&p.DataPagamento,
			&p.Detalhe,
			&p.Categoria,
			&p.Valor,
			&p.Observacao,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler linha: %w", err)
		}
		pagamentos = append(pagamentos, p)
	}

	return pagamentos, nil
}

func DeletePagamento(id string) error {
	conn, err := OpenConn()
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `DELETE FROM obra.pagamento WHERE id = $1`
	_, err = conn.Exec(sqlStatement, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir pagamento: %w", err)
	}

	return nil
}

func UpdateObra(obra Obra) error {
	conn, err := OpenConn()
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `
		UPDATE obra.cadastroobra SET
			nome = $1,
			endereco = $2,
			bairro = $3,
			area = $4,
			tipo = $5,
			casagerminada = $6,
			status = $7,
			data_inicio_obra = $8,
			data_final_obra = $9,
			updated_at = now()
		WHERE idObra = $10`

	_, err = conn.Exec(sqlStatement,
		obra.Nome,
		obra.Endereco,
		obra.Bairro,
		obra.Area,
		obra.Tipo,
		obra.Casagerminada,
		obra.Status,
		obra.DataInicioObra,
		obra.DataFinalObra,
		obra.ID,
	)

	return err
}

func UpdatePagamento(p Pagamento) error {
	conn, err := OpenConn()
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `
		UPDATE obra.pagamento SET
			data_do_pagamento = $1,
			detalhe = $2,
			categoria = $3,
			valor = $4,
			observacao = $5,
			updated_at = now()
		WHERE id = $6`

	_, err = conn.Exec(sqlStatement,
		p.DataPagamento,
		p.Detalhe,
		p.Categoria,
		p.Valor,
		p.Observacao,
		p.ID,
	)

	return err
}

func GetObraByID(idObra string) (Obra, error) {
	conn, err := OpenConn()
	if err != nil {
		return Obra{}, fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `
		SELECT idObra, nome, endereco, bairro, area, tipo, casagerminada, status, data_inicio_obra, data_final_obra, created_at, updated_at
		FROM obra.cadastroobra
		WHERE idObra = $1`

	var u Obra
	err = conn.QueryRow(sqlStatement, idObra).Scan(
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
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return Obra{}, nil
		}
		return Obra{}, fmt.Errorf("erro ao buscar obra: %w", err)
	}

	return u, nil
}
