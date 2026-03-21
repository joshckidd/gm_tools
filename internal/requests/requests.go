// internal package that the cli uses to send requests to the rest api

package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/joshckidd/gm_tools/internal/config"

	"net/http"
	"net/url"
)

// call the api and unmarshal the result into the data type provided
func CallApi[T any](cfg *config.CliConfig, endpoint, method string, payload any) (T, error) {
	var records T

	body, err := sendRequest(*cfg, endpoint, method, payload)
	if err != nil {
		return records, err
	}

	err = json.Unmarshal(body, &records)
	return records, err
}

// make a login call to the api
func LoginUser(cfg *config.CliConfig, username, password string) error {
	loginResult, err := CallApi[map[string]string](
		cfg,
		"login",
		"POST",
		map[string]string{
			"username": username,
			"password": password,
		})
	if err != nil {
		return err
	}

	cfg.CurrentUserToken = loginResult["token"]
	return nil
}

// send a request to the api
func sendRequest(cfg config.CliConfig, endpoint, method string, payload any) ([]byte, error) {
	val, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, err
	}

	urlParts := strings.Split(endpoint, "?")

	apiURL, err := url.JoinPath(cfg.APIUrl, urlParts[0])
	if err != nil {
		return []byte{}, err
	}
	if len(urlParts) > 1 {
		apiURL = fmt.Sprintf("%s?%s", apiURL, urlParts[1])
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
		return body, fmt.Errorf("Response code: %v\nBody: %s", res.StatusCode, string(body))
	}

	return body, nil
}
