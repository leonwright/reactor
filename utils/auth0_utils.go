package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"strings"

	"github.com/leonwright/reactor/models"
)

const auth0ApiBase = "https://reactorapp.us.auth0.com"

// RequestNewManagementAPIToken gets an access token for the management API and
// stores it to Redis Memory.
func RequestNewManagementAPIToken(cfg Config) {
	url := auth0ApiBase + "/oauth/token"

	payload := strings.NewReader("{\"client_id\":\"" + cfg.Auth0.APIClientID + "\",\"client_secret\":\"" + cfg.Auth0.APIClientSecret + "\",\"audience\":\"https://reactorapp.us.auth0.com/api/v2/\",\"grant_type\":\"client_credentials\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	var tokenResp models.Auth0TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResp); err != nil {
		log.Fatalln(err)
	}
	UpdateManagementAPIToken(tokenResp.AccessToken)
}

// GetUserByID gets a user from the Auth0 Management API
func GetUserByID(userID string, token string) models.Auth0User {
	userID = url.PathEscape(userID)
	URL := auth0ApiBase + "/api/v2/users/" + userID

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var auth0Resp models.Auth0User
	if err := json.NewDecoder(resp.Body).Decode(&auth0Resp); err != nil {

		log.Fatalln(err)
	}

	return auth0Resp
}
