package requests

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/joshckidd/gm_tools/internal/rolls"

	"net/http"
)

func GenerateRoll(rollString string) (rolls.RollTotalResult, error) {
	apiURL := fmt.Sprintf("http://localhost:8080/api?roll=%s", rollString)

	res, err := http.Get(apiURL)
	if err != nil {
		return rolls.RollTotalResult{}, err
	}

	body, err := io.ReadAll(res.Body)

	res.Body.Close()
	if res.StatusCode > 299 {
		return rolls.RollTotalResult{}, fmt.Errorf("Response code: %v", res.StatusCode)
	}
	if err != nil {
		return rolls.RollTotalResult{}, err
	}

	var rollResult rolls.RollTotalResult
	err = json.Unmarshal(body, &rollResult)
	if err != nil {
		return rolls.RollTotalResult{}, err
	}
	return rollResult, nil
}
