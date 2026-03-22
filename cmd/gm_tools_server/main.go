// the main package for the rest api server

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
	// load variables from the .env file
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	checkMinutes, _ := strconv.Atoi(os.Getenv("CHECK_MINUTES"))
	keepRollDays, _ := strconv.Atoi(os.Getenv("KEEP_ROLL_DAYS"))

	// connect to the postgres database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Database error.")
		os.Exit(1)
	}

	dbQueries := database.New(db)

	// create a function that runs periodically and clears out old rolls and instances
	// rolls and instances are expected to be temporary
	// duration to keep old rolls and the interval to check are in the .env file
	ticker := time.NewTicker(time.Minute * time.Duration(checkMinutes))

	go func() error {
		for ; ; <-ticker.C {
			err = dbQueries.DeleteOldRolls(context.Background(), time.Now().Add(time.Duration(-keepRollDays)*24*time.Hour))
			if err != nil {
				return err
			}
			err = dbQueries.DeleteOldInstances(context.Background(), time.Now().Add(time.Duration(-keepRollDays)*24*time.Hour))
			if err != nil {
				return err
			}
		}
	}()

	// set up a new server
	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	// define api endpoints
	var apiCfg responses.ApiConfig
	apiCfg.DB = dbQueries
	apiCfg.TokenSecret = os.Getenv("SECRET")
	serveMux.HandleFunc("GET /api/gm_tools", apiCfg.GetStatus)
	serveMux.HandleFunc("POST /api/rolls", apiCfg.ApiLogin(responses.PostRoll))
	serveMux.HandleFunc("GET /api/rolls", apiCfg.ApiLogin(responses.GetRolls))
	serveMux.HandleFunc("POST /api/users", apiCfg.PostUser)
	serveMux.HandleFunc("POST /api/login", apiCfg.UserLogin)
	serveMux.HandleFunc("POST /api/types", apiCfg.ApiLogin(responses.PostType))
	serveMux.HandleFunc("GET /api/types", apiCfg.ApiLogin(responses.GetTypes))
	serveMux.HandleFunc("GET /api/types/{typeId}", apiCfg.ApiLogin(responses.GetType))
	serveMux.HandleFunc("DELETE /api/types/{typeId}", apiCfg.ApiLogin(responses.DeleteType))
	serveMux.HandleFunc("PUT /api/types/{typeId}", apiCfg.ApiLogin(responses.PutType))
	serveMux.HandleFunc("POST /api/custom_fields", apiCfg.ApiLogin(responses.PostCustomField))
	serveMux.HandleFunc("GET /api/custom_fields", apiCfg.ApiLogin(responses.GetCustomFields))
	serveMux.HandleFunc("GET /api/custom_fields/{fieldId}", apiCfg.ApiLogin(responses.GetCustomField))
	serveMux.HandleFunc("DELETE /api/custom_fields/{fieldId}", apiCfg.ApiLogin(responses.DeleteCustomField))
	serveMux.HandleFunc("PUT /api/custom_fields/{fieldId}", apiCfg.ApiLogin(responses.PutCustomField))
	serveMux.HandleFunc("POST /api/items", apiCfg.ApiLogin(responses.PostItem))
	serveMux.HandleFunc("GET /api/items", apiCfg.ApiLogin(responses.GetItems))
	serveMux.HandleFunc("GET /api/items/{itemId}", apiCfg.ApiLogin(responses.GetItem))
	serveMux.HandleFunc("DELETE /api/items/{itemId}", apiCfg.ApiLogin(responses.DeleteItem))
	serveMux.HandleFunc("PUT /api/items/{itemId}", apiCfg.ApiLogin(responses.PutItem))
	serveMux.HandleFunc("POST /api/instances", apiCfg.ApiLogin(responses.PostInstances))
	serveMux.HandleFunc("GET /api/instances", apiCfg.ApiLogin(responses.GetInstances))

	server.ListenAndServe()
}
