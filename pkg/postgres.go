package pkg

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Obra struct {
	sequence      int
	nome          string
	endereco      string
	bairro        string
	area          string
	tipo          int
	casagerminada bool
	status        bool
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

func InsertObra(data string, nome string, endereco string, bairro string, area string, tipo int, casagerminada bool) {
	conn, err := OpenConn()
	if err != nil {
		return
	}
	defer conn.Close()

	sqlStatement :=
		`INSERT INTO obra.cadastroobra ("createdDate","nome","endereco","bairro","area","tipo","casagerminada","status")
		VALUES ($1,$2,$3,$4,$5,$6,$7,true)`

	conn.QueryRow(sqlStatement, data, nome, endereco, bairro, area, tipo, casagerminada)
	if err != nil {
		panic(err)
	}
}

func GetAllObra() ([]Obra, error) {
	conn, err := OpenConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	sqlStatement :=
		`SELECT c.sequence , c.nome , c.endereco , c.bairro , c.area , c.tipo , c.casagerminada , c.status FROM obra.cadastroobra c ORDER BY sequence ASC`
	rows, err := conn.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	currentTime := time.Now()
	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | CA - Consulta realizada")
	var obras []Obra
	for rows.Next() {
		var u Obra
		if err := rows.Scan(&u.sequence, &u.nome, &u.endereco, &u.bairro, &u.area, &u.tipo, &u.casagerminada, &u.status); err != nil {
			return nil, err
		}
		fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | CA - Consulta realizada : " + u.nome)
		obras = append(obras, u)
	}
	return obras, nil
}

/*
func GetObraId(id string) ([]Obra, error) {
	conn, err := OpenConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	sqlStatement :=
		`SELECT c.sequence , c.nome , c.endereco , c.bairro , c.area , c.tipo , c.casagerminada , c.status FROM obra.cadastroobra c where c.sequence = '$1' ORDER BY sequence ASC`

	rows, err := conn.QueryRow(sqlStatement , id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var obras []Obra
	for rows.Next() {
		var u Obra
		if err := rows.Scan(&u.sequence, &u.nome, &u.endereco, &u.bairro, &u.area, &u.tipo, &u.casagerminada, &u.status); err != nil {
			return nil, err
		}
		obras = append(obras, u)
	}

	return obras, nil
}
*/
