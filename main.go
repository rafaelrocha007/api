package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type cep struct {
	Cidade     string `json:"cidade"`
	Bairro     string `json:"bairro"`
	Logradouro string `json:"logradouro"`
	UF         string `json:"uf"`
}

func (c cep) exist() bool {
	return len(c.UF) != 0
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

	requestedCep := strings.Split(r.URL.Path[1:], "/")[1]
	if err := isValidCEP(requestedCep); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	timeout := time.Second * 300
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tube := make(chan []byte, 1)
	for source, url := range endpoints {
		endpoint := fmt.Sprintf(url, requestedCep)
		go request(ctx, endpoint, source, tube)
	}

	w.Header().Set("Content-Type", "application/json")

	for index := 0; index < 3; index++ {
		cepInformation, err := parseResponse(<-tube)
		if err != nil {
			continue
		}

		if cepInformation.exist() {
			json.NewEncoder(w).Encode(cepInformation)
			return
		}
	}

	http.Error(w, "", http.StatusNoContent)
}

func request(ctx context.Context, endpoint, source string, tube chan []byte) {
	start := time.Now()

	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Printf("Could not get from %s - %s \n", source, err.Error())
		return
	}

	request = request.WithContext(ctx)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("Fail to request data from %s \n", source)
		return
	}
	defer response.Body.Close()

	requestContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Could not get payload from %s - %s \n", source, err.Error())
		return
	}

	if len(requestContent) != 0 && response.StatusCode == http.StatusOK {
		fmt.Printf("Endpoint %s took: %s \n", source, time.Since(start))
		tube <- requestContent
	}
}

// parseResponse formata a resposta para uma saida padrão baseada na struct de cep,
// alguns serviços de cep tem respostas diferentes que usa a palavra
// "localidade" para definir a cidade e estado para definir a UF.
func parseResponse(content []byte) (payload cep, err error) {
	response := make(map[string]interface{})
	_ = json.Unmarshal(content, &response)

	if err := isValidResponse(response); !err {
		return payload, errors.New("invalid response")
	}

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

func isValidCEP(cep string) error {
	re := regexp.MustCompile(`[^0-9]`)
	formatedCEP := re.ReplaceAllString(cep, `$1`)

	if len(formatedCEP) < 8 {
		return errors.New("Cep deve conter apenas números e no minimo 8 digitos")
	}

	return nil
}

func isValidResponse(requestContent map[string]interface{}) bool {
	if _, ok := requestContent["erro"]; ok {
		return false
	}

	if _, ok := requestContent["fail"]; ok {
		return false
	}

	return true
}
