package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/joshckidd/gm_tools/internal/database"
	"github.com/joshckidd/gm_tools/internal/responses"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Database error.")
		os.Exit(1)
	}

	dbQueries := database.New(db)

	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	var apiCfg responses.ApiConfig
	apiCfg.DB = dbQueries
	apiCfg.TokenSecret = os.Getenv("SECRET")
	serveMux.HandleFunc("POST /api/rolls", apiCfg.PostRoll)
	serveMux.HandleFunc("POST /api/users", apiCfg.PostUser)
	serveMux.HandleFunc("POST /api/login", apiCfg.UserLogin)

	server.ListenAndServe()
}
