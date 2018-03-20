package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hey dude %s\n", r.URL.Query().Get("name"))
	})

	http.ListenAndServe(":8000", nil)
}
