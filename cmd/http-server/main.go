package main

import (
	"fmt"
	"net/http"

	"git.qowevisa.me/Qowevisa/gonuts/db"
)

func main() {
	dbc := db.Connect()
	if dbc != nil {
		fmt.Printf("yay\n")
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.ListenAndServe(":8080", nil)
}
