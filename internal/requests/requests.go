package requests

import (
	"fmt"

	"github.com/joshckidd/gm_tools/internal/rolls"

	"net/http"
)

func GetRoll(w http.ResponseWriter, r *http.Request) {
	rollString := r.URL.Query().Get("roll")

	tot, _ := rolls.RollAll(rolls.ParseRoll(rollString))

	w.Header().Set("Content-Type", "application/html")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%d", tot)
}
