package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/joshckidd/gm_tools/internal/config"

	"net/http"
	"net/url"
)

func CallApi[T any](cfg *config.CliConfig, endpoint, method string, payload any) (T, error) {
	var records T

	body, err := sendRequest(*cfg, endpoint, method, payload)
	if err != nil {
		return records, err
	}

	err = json.Unmarshal(body, &records)
	return records, nil
}

func LoginUser(cfg *config.CliConfig, username, password string) error {
	body, err := sendRequest(*cfg, "login", "POST", map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return err
	}

	var loginResult struct {
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Token     string    `json:"token"`
	}
	err = json.Unmarshal(body, &loginResult)
	if err != nil {
		return err
	}
	cfg.CurrentUserToken = loginResult.Token
	return nil
}

func sendRequest(cfg config.CliConfig, endpoint, method string, payload any) ([]byte, error) {
	val, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, err
	}

	apiURL, err := url.JoinPath(cfg.APIUrl, endpoint)
	if err != nil {
		return []byte{}, err
	}

	client := &http.Client{}

	req, err := http.NewRequestWithContext(context.Background(), method, apiURL, bytes.NewBuffer([]byte(val)))
	if err != nil {
		return []byte{}, err
	}

	tok := fmt.Sprintf("Bearer %s", cfg.CurrentUserToken)

	req.Header.Set("User-Agent", "gm-tools")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tok)

	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return body, err
	}

	if res.StatusCode > 299 {
		return body, fmt.Errorf("Response code: %v", res.StatusCode)
	}

	return body, nil
}
