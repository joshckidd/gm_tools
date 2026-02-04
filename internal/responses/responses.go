package responses

import (
	"encoding/json"
	"fmt"

	"github.com/joshckidd/gm_tools/internal/auth"
	"github.com/joshckidd/gm_tools/internal/database"
	"github.com/joshckidd/gm_tools/internal/rolls"

	"net/http"
)

type ApiConfig struct {
	DB *database.Queries
}

func GetRoll(w http.ResponseWriter, r *http.Request) {
	rollString := r.URL.Query().Get("roll")

	tot := rolls.RollAll(rolls.ParseRoll(rollString))

	respondWithJSON(w, 200, tot)
}

func (cfg *ApiConfig) PostUser(w http.ResponseWriter, r *http.Request) {
	type userParam struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	inParams := userParam{}

	err := decoder.Decode(&inParams)
	if err != nil {
		respondWithError(w, 500, "Invalid request")
		return
	}

	hashedPassword, err := auth.HashPassword(inParams.Password)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	params := database.CreateUserParams{
		Username:       inParams.Username,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.DB.CreateUser(r.Context(), params)
	if err.Error() == "pq: duplicate key value violates unique constraint \"users_pkey\"" {
		respondWithError(w, 409, fmt.Sprintf("%s is already in use as a username. Please select another.", inParams.Username))
		return
	} else if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 201, database.User{
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Username:  user.Username,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	val, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(val)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnError struct {
		Error string `json:"error"`
	}

	respError := returnError{
		Error: msg,
	}

	dat, err := json.Marshal(respError)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
