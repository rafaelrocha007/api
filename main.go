package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
		port = "8080"
	}

	fmt.Println("Listening on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "not allowed", http.StatusMethodNotAllowed)
		return
	}

	tube := make(chan []byte, 3)
	for source, url := range endpoints {
		endpoint := fmt.Sprintf(url, r.URL.Path[1:])
		go request(endpoint, source, tube)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", string(<-tube))
}

func request(endpoint, source string, tube chan []byte) {
	request, err := http.Get(endpoint)
	if err != nil {
		fmt.Printf("Could not get from %s - %s \n", source, err.Error())
	}
	defer request.Body.Close()

	response, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Printf("Could not get payload from %s - %s", source, err.Error())
	}

	fmt.Printf("endpoint %s won \n", source)
	tube <- response
}
