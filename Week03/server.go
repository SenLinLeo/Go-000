package main

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "func1......")
}

func main() {
	http.HandleFunc("/index", indexHandler)
	http.ListenAndServe(":8080", nil)
}
