package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/joshckidd/gm_tools/internal/config"
	"github.com/joshckidd/gm_tools/internal/rolls"

	"net/http"
)

func GenerateRoll(cfg config.CliConfig, rollString string) (rolls.RollTotalResult, error) {
	apiURL := fmt.Sprintf("%s/rolls", cfg.APIUrl)

	client := &http.Client{}

	roll := fmt.Sprintf("{\"roll\": \"%s\"}", rollString)

	req, err := http.NewRequestWithContext(context.Background(), "POST", apiURL, bytes.NewBuffer([]byte(roll)))
	if err != nil {
		return rolls.RollTotalResult{}, err
	}

	tok := fmt.Sprintf("Bearer %s", cfg.CurrentUserToken)

	req.Header.Set("User-Agent", "gm-tools")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tok)

	res, err := client.Do(req)
	if err != nil {
		return rolls.RollTotalResult{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return rolls.RollTotalResult{}, err
	}

	if res.StatusCode > 299 {
		return rolls.RollTotalResult{}, fmt.Errorf("Response code: %v", res.StatusCode)
	}

	var rollResult rolls.RollTotalResult
	err = json.Unmarshal(body, &rollResult)
	if err != nil {
		return rolls.RollTotalResult{}, err
	}
	return rollResult, nil
}
