package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ApiCEP struct {
	Code       string `json:"code"`
	Address    string `json:"address"`
	District   string `json:"district"`
	City       string `json:"city"`
	State      string `json:"state"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	chViaCEP := make(chan ViaCEP)
	chApiCEP := make(chan ApiCEP)

	for _, cep := range os.Args[1:] {
		if len(cep) == 8 {
			go getApiCep(cep, chApiCEP)
			go getViaCep(cep, chViaCEP)
		}
	}

	select {
	case resp := <-chViaCEP:
		fmt.Printf("ViaCEP\n\nCEP: %v\nLogradouro: %v\nBairro: %v\nLocalidade: %v\nUF: %v\n", resp.Cep, resp.Logradouro, resp.Bairro, resp.Localidade, resp.Uf)

	case resp := <-chApiCEP:
		fmt.Printf("ApiCEP\n\nCEP: %v\nLogradouro: %v\nBairro: %v\nLocalidade: %v\nUF: %v\n", resp.Code, resp.Address, resp.District, resp.City, resp.State)

	case <-time.After(time.Second):
		println("timeout")
	}
}

func getApiCep(cep string, ch chan<- ApiCEP) {
	req, err := http.Get("https://cdn.apicep.com/file/apicep/" + fmt.Sprintf("%s-%s", cep[:5], cep[len(cep)-3:]) + ".json")
	if err != nil {
		fmt.Printf("ApiCEP -> Erro ao fazer requisição: %v\n", err)
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("ApiCEP -> Erro ao ler resposta: %v\n", err)
	}

	var data ApiCEP
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Printf("ApiCEP -> Eerro ao fazer parse da resposta: %v\n", err)
	}

	ch <- data
}

func getViaCep(cep string, ch chan<- ViaCEP) {
	req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		fmt.Printf("ViaCEP -> Erro ao fazer requisição: %v\n", err)
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("ViaCEP -> Erro ao ler resposta: %v\n", err)
	}

	var data ViaCEP
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Printf("ViaCEP -> Erro ao fazer parse da resposta: %v\n", err)
	}

	ch <- data
}
