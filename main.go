package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	db "DB"
	pages "Pages"
)

const connection = "mongodb+srv://ImBlue:abubakat2ahmed@charactersheetdb.zre7u.mongodb.net/CharacterSheetDB?retryWrites=true&w=majority"

func main() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	//http.HandleFunc("/projectinfo/v1/", serviceHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/index/", indexHandler)
	log.Printf("Listening on %s...\n", addr)

	client, ctx, err := db.GetDBConnection()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := pages.Index{}
	data.Title = "Hello World"
	data.Sheets = []string{"One", "Two", "Three"}
	data.LoggedIn = false
	t, _ := template.ParseFiles("./templates/index.html")
	t.Execute(w, data)
}
