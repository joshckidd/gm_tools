package requests

import (
	"encoding/json"

	"github.com/joshckidd/gm_tools/internal/rolls"

	"net/http"
)

func GetRoll(w http.ResponseWriter, r *http.Request) {
	rollString := r.URL.Query().Get("roll")

	tot := rolls.RollAll(rolls.ParseRoll(rollString))

	respondWithJSON(w, 200, tot)
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
