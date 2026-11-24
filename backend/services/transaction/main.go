package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from transaction service")
	})

	log.Printf("transaction service listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
