package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CdnCEP struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
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

type CEP struct {
	Codigo   string
	Endereco string
	Bairro   string
	Cidade   string
	Uf       string
	Api      string
}

func main() {
	c := make(chan CEP)
	go BuscaCdnCEP("71218-010", c)
	go BuscaViaCEP("71218-010", c)

	select {
	case cep := <-c:
		fmt.Printf("Received from %s:\n", cep.Api)
		fmt.Println(cep)
		close(c)

	case <-time.After(time.Second * 1):
		fmt.Println("timeout")
		close(c)
	}
}

func BuscaCdnCEP(codigo string, c chan CEP) {
	resp, err := http.Get("https://cdn.apicep.com/file/apicep/" + codigo + ".json")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var cdnCEP CdnCEP
	err = json.Unmarshal(body, &cdnCEP)
	if err != nil {
		return
	}

	if cdnCEP.Ok {
		cep := CEP{
			Codigo:   cdnCEP.Code,
			Endereco: cdnCEP.Address,
			Bairro:   cdnCEP.District,
			Cidade:   cdnCEP.City,
			Uf:       cdnCEP.State,
			Api:      "CDN",
		}

		c <- cep
	}
}

func BuscaViaCEP(codigo string, c chan CEP) {
	resp, err := http.Get("http://viacep.com.br/ws/" + codigo + "/json/")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var viaCEP ViaCEP
	err = json.Unmarshal(body, &viaCEP)
	if err != nil {
		return
	}

	cep := CEP{
		Codigo:   viaCEP.Cep,
		Endereco: viaCEP.Logradouro,
		Bairro:   viaCEP.Bairro,
		Cidade:   viaCEP.Localidade,
		Uf:       viaCEP.Uf,
		Api:      "VIACEP",
	}

	c <- cep
}
