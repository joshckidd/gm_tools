package main

import (
	"net/http"

	"github.com/joshckidd/gm_tools/internal/requests"

	_ "github.com/lib/pq"
)

func main() {

	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	serveMux.HandleFunc("GET /api", requests.GetRoll)

	server.ListenAndServe()
}
