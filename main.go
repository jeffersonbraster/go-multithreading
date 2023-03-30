package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	VIA_CEP_URL = "https://viacep.com.br/ws/%s/json/"
	API_CEP_URL = "https://cdn.apicep.com/file/apicep/%s.json"
	CEP         = "60821-490"
)

type ViaCepResponse struct {
	Cep          string `json:"cep"`
	Logradouro   string `json:"logradouro"`
	Complemento  string `json:"complemento"`
	Bairro       string `json:"bairro"`
	Localidade   string `json:"localidade"`
	Uf           string `json:"uf"`
	Ibge         string `json:"ibge"`
	Gia          string `json:"gia"`
	Ddd          string `json:"ddd"`
	Siafi        string `json:"siafi"`
}

type ApiCepResponse struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func main() {
	chViaCep := make(chan ViaCepResponse)
	chApiCep := make(chan ApiCepResponse)

	go func() {
		response := RequestApi(fmt.Sprintf(VIA_CEP_URL, CEP))

		if response != nil {
			var viaCepResponse ViaCepResponse
			json.Unmarshal(response, &viaCepResponse)
			chViaCep <- viaCepResponse
		}
	}()

	go func() {
		response := RequestApi(fmt.Sprintf(API_CEP_URL, CEP))

		if response != nil {
			var apiCepResponse ApiCepResponse
			json.Unmarshal(response, &apiCepResponse)
			chApiCep <- apiCepResponse
		}
	}()

	select {
	case response := <-chViaCep:
		fmt.Printf("ViaCepResponse: %+v\n", response)
	case response := <-chApiCep:
		fmt.Printf("ApiCepResponse: %+v\n", response)
	case <-time.After(time.Second * 1):
		fmt.Println("Timeout")
	default:
		fmt.Println("Nenhum case selecionado")
	}
}

func RequestApi(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return body
}