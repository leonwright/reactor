package routes

import (
	"fmt"
	"net/http"

	"github.com/leonwright/reactor/utils"
)

// RegisterWithBuilder gets the user's Github Access Token and store's it in memory
// for future requests.
func RegisterWithBuilder(w http.ResponseWriter, r *http.Request) {
	deb.Info("RegisterWithBuilder called...")
	token := utils.ExtractToken(r)

	hasScope, username := utils.CheckScope("read:self", token)
	if !hasScope {
		message := "Insufficient scope."
		utils.SendJSONResponse(message, w)
		return
	}
	err := utils.ValidateRequestIsJSON(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	auth0Token, err := utils.GetManagementAPIToken()
	if err != nil {
		utils.SendJSONResponse(err, w)
	}

	var user = utils.GetUserByID(username, auth0Token)

	var githubAccessToken string

	for _, s := range user.Identities {
		if s.Provider == "github" {
			githubAccessToken = *s.AccessToken
		}
	}

	utils.UpdateGithubToken(username, githubAccessToken)

	utils.SendResponseWithoutData("Successfully Registered!", http.StatusCreated, w)

}
