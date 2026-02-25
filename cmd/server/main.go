package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/joshckidd/gm_tools/internal/database"
	"github.com/joshckidd/gm_tools/internal/responses"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	checkMinutes, _ := strconv.Atoi(os.Getenv("CHECK_MINUTES"))
	keepRollDays, _ := strconv.Atoi(os.Getenv("KEEP_ROLL_DAYS"))

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Database error.")
		os.Exit(1)
	}

	dbQueries := database.New(db)

	ticker := time.NewTicker(time.Minute * time.Duration(checkMinutes))

	go func() error {
		for ; ; <-ticker.C {
			err = dbQueries.DeleteOldRolls(context.Background(), time.Now().Add(time.Duration(-keepRollDays)*24*time.Hour))
			if err != nil {
				return err
			}
		}
	}()

	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	var apiCfg responses.ApiConfig
	apiCfg.DB = dbQueries
	apiCfg.TokenSecret = os.Getenv("SECRET")
	serveMux.HandleFunc("POST /api/rolls", apiCfg.ApiLogin(responses.PostRoll))
	serveMux.HandleFunc("GET /api/rolls", apiCfg.ApiLogin(responses.GetRolls))
	serveMux.HandleFunc("POST /api/users", apiCfg.PostUser)
	serveMux.HandleFunc("POST /api/login", apiCfg.UserLogin)
	serveMux.HandleFunc("POST /api/types", apiCfg.ApiLogin(responses.PostType))
	serveMux.HandleFunc("GET /api/types", apiCfg.ApiLogin(responses.GetTypes))
	serveMux.HandleFunc("POST /api/custom_fields", apiCfg.ApiLogin(responses.PostCustomField))
	serveMux.HandleFunc("GET /api/custom_fields", apiCfg.ApiLogin(responses.GetCustomFields))

	server.ListenAndServe()
}
