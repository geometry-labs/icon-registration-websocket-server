package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port_env := os.Getenv("MOCK_REGISTRATION_API_PORT")

	http.HandleFunc("/broadcaster/register", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"broadcaster_id": "test-broadcaster-id"}`)
	})
	http.HandleFunc("/broadcaster/unregister", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"err": ""}`)
	})

	log.Fatal(http.ListenAndServe(":"+port_env, nil))
}
