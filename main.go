package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	//http.HandleFunc("/projectinfo/v1/", serviceHandler)
	http.HandleFunc("/hello", helloHandler)
	log.Printf("Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}
