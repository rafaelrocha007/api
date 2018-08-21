package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var endpoints = map[string]string{
	"viacep":           "https://viacep.com.br/ws/%s/json/",
	"postmon":          "https://api.postmon.com.br/v1/cep/%s",
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
	fmt.Fprintf(w, "%v", string(<-tube))
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
		fmt.Printf("Endpoint %s won - time %s \n", source, time.Since(start))
		tube <- response
	}
}
