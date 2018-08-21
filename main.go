package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type cep struct {
	Cidade     string `json:"cidade"`
	Bairro     string `json:"bairro"`
	Logradouro string `json:"logradouro"`
	UF         string `json:"uf"`
}

var endpoints = map[string]string{
	"viacep":           "http://viacep.com.br/ws/%s/json/",
	"postmon":          "http://api.postmon.com.br/v1/cep/%s",
	"republicavirtual": "http://republicavirtual.com.br/web_cep.php?cep=%s&formato=json",
}

func main() {
	http.HandleFunc("/v1/", handler)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	fmt.Println("Listening on port: " + port)
	http.ListenAndServe(":"+port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "not allowed", http.StatusMethodNotAllowed)
		return
	}

	tube := make(chan []byte, 1)
	cep := strings.Split(r.URL.Path[1:], "/")[1]
	for source, url := range endpoints {
		endpoint := fmt.Sprintf(url, cep)
		go request(endpoint, source, tube)
	}

	w.Header().Set("Content-Type", "application/json")
	response := parseResponse(<-tube)
	json.NewEncoder(w).Encode(response)
}

func request(endpoint, source string, tube chan []byte) {
	start := time.Now()

	request, err := http.Get(endpoint)
	if err != nil {
		fmt.Printf("Could not get from %s - %s \n", source, err.Error())
	}
	defer request.Body.Close()

	response, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Printf("Could not get payload from %s - %s", source, err.Error())
	}

	if len(response) != 0 && request.StatusCode == http.StatusOK {
		fmt.Printf("Endpoint %s took: %s \n", source, time.Since(start))
		tube <- response
	}
}

func parseResponse(content []byte) (payload cep) {
	response := make(map[string]interface{})
	_ = json.Unmarshal(content, &response)

	if _, ok := response["localidade"]; ok {
		payload.Cidade = response["localidade"].(string)
	} else {
		payload.Cidade = response["cidade"].(string)
	}

	if _, ok := response["estado"]; ok {
		payload.UF = response["estado"].(string)
	} else {
		payload.UF = response["uf"].(string)
	}

	payload.Bairro = response["bairro"].(string)
	payload.Logradouro = response["logradouro"].(string)

	return
}
