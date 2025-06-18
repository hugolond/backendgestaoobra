package pkg

import (
	models "backendgestaoobra/model"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type CategoriaProps struct {
	ID       string `json:"id"`
	Tipo     string `json:"tipo"`
	Campo    string `json:"campo"`
	Subcampo string `json:"subcampo"`
	Titulo   string `json:"titulo"`
	Status   bool   `json:"status"`
}

type ObraPagamento struct {
	IDObra        string  `json:"idobra"`
	Nome          string  `json:"nome"`
	Previsto      float64 `json:"previsto"`
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
	Previsto       float64
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

func GetAllProps() ([]CategoriaProps, error) {
	conn, err := OpenConn()
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `
		SELECT id, tipo, campo, subcampo, titulo, status
		FROM obra.props
		ORDER BY campo, subcampo;
	`

	rows, err := conn.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	var categorias []CategoriaProps
	for rows.Next() {
		var cat CategoriaProps
		err := rows.Scan(&cat.ID, &cat.Tipo, &cat.Campo, &cat.Subcampo, &cat.Titulo, &cat.Status)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear linha: %w", err)
		}
		categorias = append(categorias, cat)
	}

	return categorias, nil
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

func SelectObraPagamentoJoin(accountID string) ([]ObraPagamento, error) {
	conn, err := OpenConn()
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	query := `
		SELECT
			o.idObra,
			o.nome,
			COALESCE(o.previsto, 0) as previsto,
			COALESCE(p.data_do_pagamento, '2024-01-01') AS data_do_pagamento,
			COALESCE(p.valor, 0),
			COALESCE(p.categoria, '')
			FROM obra.cadastroobra o
			INNER JOIN obra.pagamento p ON p.idObra = o.idObra AND p.account_id = $1
			ORDER BY o.nome, p.data_do_pagamento DESC;
	`

	rows, err := conn.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	var dados []ObraPagamento
	for rows.Next() {
		var linha ObraPagamento
		err := rows.Scan(&linha.IDObra, &linha.Nome, &linha.Previsto, &linha.DataPagamento, &linha.Valor, &linha.Categoria)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler linha: %w", err)
		}
		dados = append(dados, linha)
	}

	return dados, nil
}

func InsertObra(obra Obra, accountID string, userID string, userName string) (string, error) {
	conn, err := OpenConn()
	if err != nil {
		return "", fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `
		INSERT INTO obra.cadastroobra (
			nome, endereco, bairro, area, tipo, previsto, casagerminada, status, data_inicio_obra, data_final_obra, created_at, updated_at,account_id,userid_at,username_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now(), now() ,$11 ,$12 ,$13
		) RETURNING idObra`

	var idObra string
	err = conn.QueryRow(sqlStatement,
		obra.Nome,
		obra.Endereco,
		obra.Bairro,
		obra.Area,
		obra.Tipo,
		obra.Previsto,
		obra.Casagerminada,
		obra.Status,
		obra.DataInicioObra,
		obra.DataFinalObra,
		accountID,
		userID,
		userName,
	).Scan(&idObra)

	if err != nil {
		return "", fmt.Errorf("erro ao inserir obra: %w", err)
	}

	return idObra, nil
}

func GetAllObra(accountID string) ([]Obra, error) {
	conn, err := OpenConn()
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `
		SELECT idObra, nome, endereco, bairro, area, tipo, COALESCE(previsto, 0), casagerminada, status, data_inicio_obra, data_final_obra, created_at, updated_at
		FROM obra.cadastroobra
		WHERE account_id = $1
		ORDER BY data_inicio_obra DESC`

	rows, err := conn.Query(sqlStatement, accountID)
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
			&u.Previsto,
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

func InsertPagamentoStruct(p Pagamento, accountID string, userID string, userName string) error {
	conn, err := OpenConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	sqlStatement := `
		INSERT INTO obra.pagamento (
			idObra, data_do_pagamento, detalhe, categoria, valor, observacao, created_at, updated_at, account_id,userid_at,username_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, now(), now(), $7, $8, $9
		)`

	_, err = conn.Exec(sqlStatement,
		p.IDObra,
		p.DataPagamento,
		p.Detalhe,
		p.Categoria,
		p.Valor,
		p.Observacao,
		accountID,
		userID,
		userName,
	)

	return err
}

func GetAllPagamentoByObra(idObra string, accountID string) ([]Pagamento, error) {
	conn, err := OpenConn()
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `
		SELECT id, idObra, data_do_pagamento, detalhe, categoria, valor, observacao, created_at, updated_at
		FROM obra.pagamento
		WHERE idObra = $1 AND account_id = $2
		ORDER BY data_do_pagamento DESC`

	rows, err := conn.Query(sqlStatement, idObra, accountID)
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

func DeletePagamento(id string, accountId string) error {
	conn, err := OpenConn()
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `DELETE FROM obra.pagamento WHERE id = $1 AND account_id = $2`
	_, err = conn.Exec(sqlStatement, id, accountId)
	if err != nil {
		return fmt.Errorf("erro ao excluir pagamento: %w", err)
	}

	return nil
}

func UpdateObra(obra Obra, accountID string, userID string, userName string) error {
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
			previsto = &6,
			casagerminada = $7,
			status = $8,
			data_inicio_obra = $9,
			data_final_obra = $10,
			updated_at = now(), 
			userid_at = $11,
			username_at = $12
		WHERE idObra = $13 AND account_id = $14`

	_, err = conn.Exec(sqlStatement,
		obra.Nome,
		obra.Endereco,
		obra.Bairro,
		obra.Area,
		obra.Tipo,
		obra.Previsto,
		obra.Casagerminada,
		obra.Status,
		obra.DataInicioObra,
		obra.DataFinalObra,
		userID,    // $10
		userName,  // $11
		obra.ID,   // $12
		accountID, // $13
	)

	return err
}

func UpdatePagamento(p Pagamento, accountID string, userID string, userName string) error {
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
			updated_at = now(),
			userid_at = $6,
			username_at = $7
		WHERE id = $8 AND account_id = $9`

	_, err = conn.Exec(sqlStatement,
		p.DataPagamento,
		p.Detalhe,
		p.Categoria,
		p.Valor,
		p.Observacao,
		userID,
		userName,
		p.ID,
		accountID,
	)

	return err
}

func GetObraByID(idObra string, accountID string) (Obra, error) {
	conn, err := OpenConn()
	if err != nil {
		return Obra{}, fmt.Errorf("erro ao abrir conexão: %w", err)
	}
	defer conn.Close()

	sqlStatement := `
		SELECT idObra, nome, endereco, bairro, area, tipo, previsto, casagerminada, status, data_inicio_obra, data_final_obra, created_at, updated_at
		FROM obra.cadastroobra
		WHERE idObra = $1 AND account_id = $2`

	var u Obra
	err = conn.QueryRow(sqlStatement, idObra, accountID).Scan(
		&u.ID,
		&u.Nome,
		&u.Endereco,
		&u.Bairro,
		&u.Area,
		&u.Tipo,
		&u.Previsto,
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

func SaveSubscription(sub models.Subscription) error {
	conn, err := OpenConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	sqlStatement := `
		INSERT INTO obra.subscriptions (
			user_id, stripe_customer, stripe_subscription,
			stripe_price_id, stripe_product_id, stripe_plan_amount,
			currency, interval, interval_count, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = conn.Exec(sqlStatement, sub.UserID,
		sub.StripeCustomer,
		sub.StripeSubscription,
		sub.StripePriceID,
		sub.StripeProductID,
		sub.StripePlanAmount,
		sub.Currency,
		sub.Interval,
		sub.IntervalCount,
		sub.Status)
	if err != nil {
		log.Println("Erro ao executar INSERT:", err)
	}
	return err
}

func CreateAccount(account models.Account) error {
	conn, err := OpenConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	query := `
		INSERT INTO obra.account (id, nome, email, stripe_product_id, status,created_at)
		VALUES ($1, $2, $3, $4, $5 ,$6)
	`
	_, err = conn.Exec(query, account.ID, account.Nome, account.Email, account.StripeProductID, account.Status, account.CreatedAt)
	if err != nil {
		log.Println("Erro ao criar account:", err)
		return err
	}
	return nil
}

func GetAccountByEmail(email string) (*models.Account, error) {
	conn, err := OpenConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := `
		SELECT id, nome, email, stripe_product_id, created_at
		FROM obra.account
		WHERE email = $1
		LIMIT 1
	`
	row := conn.QueryRow(query, email)

	var account models.Account
	err = row.Scan(&account.ID, &account.Nome, &account.Email, &account.StripeProductID, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // não encontrado
		}
		log.Println("Erro ao buscar account:", err)
		return nil, err
	}
	return &account, nil
}

func UpdateAccountPlan(accountID string, newPlan string) error {
	conn, err := OpenConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	query := `
		UPDATE obra.account
		SET stripe_product_id = $1
		WHERE id = $2
	`
	_, err = conn.Exec(query, newPlan, accountID)
	if err != nil {
		log.Printf("Erro ao atualizar plano do account (%s): %v", accountID, err)
		return err
	}
	return nil
}
