package main

import (
	"fmt"
	"net/http"
)

func handler_(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Привет, мир!")
	w.Write([]byte("!!!"))
}

func main() {
	http.HandleFunc("/", handler_)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
