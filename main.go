package main

import (
	"encoding/json"
	"fmt"
)

type cep struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
}

var endpoints = map[string]string{
	"correios":  "",
	"viacep":    "viacep.com.br/ws/%s/json/",
	"cepaberto": "",
}

func main() {
	tube := make(chan string, 3)

	for _, url := range endpoints {
		zip := fmt.Sprintf(url, "cep")
		// TODO: realizar chamadas usando goroutines para cada metodo que fornece
		// informações sobre o cep
	}

	response := <-tube
	var information cep
	content := json.Unmarshal(response, &information)

	json.NewEncoder(w).Encode(information)
}

func fromCorreios(zip string, tube chan []byte) {
	// TODO: implementar requisição para os correios
}

func fromViaCep(zip string, tube chan []byte) {
	// TODO: implementar requisiçãoo para o viacep
}

func fromCepAberto(zip string, tube chan []byte) {
	// TODO: implementar requisição para cep aberto
}
