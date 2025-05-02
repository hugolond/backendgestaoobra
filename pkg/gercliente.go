package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Cliente struct {
	Elements []struct {
		ID              int    `json:"id"`
		DataNascimento  string `json:"dataNascimento"`
		Email           string `json:"email"`
		NomeCompleto    string `json:"nomeCompleto"`
		Cpf             string `json:"cpf"`
		Celular         string `json:"celular"`
		IsPrime         bool   `json:"isPrime"`
		IsFuncionario   bool   `json:"isFuncionario"`
		IsVendaTablet   bool   `json:"isVendaTablet"`
		Genero          int    `json:"genero"`
		DataCriacao     string `json:"dataCriacao"`
		DataAtualizacao string `json:"dataAtualizacao"`
		DataUltimoLogin string `json:"dataUltimoLogin"`
		Ativo           bool   `json:"ativo"`
		WhatsAppOption  bool   `json:"whatsAppOption"`
		Enderecos       []struct {
			ID              int         `json:"id"`
			Tipo            string      `json:"tipo"`
			Logradouro      string      `json:"logradouro"`
			Bairro          string      `json:"bairro"`
			Cidade          string      `json:"cidade"`
			Estado          string      `json:"estado"`
			Cep             string      `json:"cep"`
			Numero          string      `json:"numero"`
			Complemento     string      `json:"complemento"`
			Referencia      interface{} `json:"referencia"`
			NomeRecebedor   interface{} `json:"nomeRecebedor"`
			DefaultBilling  bool        `json:"defaultBilling"`
			DefaultShipping bool        `json:"defaultShipping"`
			Ativo           bool        `json:"ativo"`
		} `json:"enderecos"`
	} `json:"elements"`
	TotalElements    int `json:"totalElements"`
	TotalPages       int `json:"totalPages"`
	PageNumber       int `json:"pageNumber"`
	PageSize         int `json:"pageSize"`
	NumberOfElements int `json:"numberOfElements"`
}

func GetCliente(key string, campo string, termo string) (resp Cliente, err error) {
	url := "http://ger-clientes.pernambucanas.com.br/api/v1/clientes/parametros?" + campo + "=" + termo
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+key)

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

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Erro ao converter Json" + err.Error())
		return
	}

	return
}
