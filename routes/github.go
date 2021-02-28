package routes

import (
	"fmt"
	"net/http"

	"github.com/leonwright/reactor/utils"
)

// GetUserRepositories ...
func GetUserRepositories(w http.ResponseWriter, r *http.Request) {
	// If the Content-Type header is present, check that it has the value
	// application/json. Note that we are using the gddo/httputil/header
	// package to parse and extract the value here, so the check works
	// even if the client includes additional charset or boundary
	// information in the header.
	deb.Info("GetUserRepositories called...")
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
	githubToken, err := utils.GetGithubToken(username)
	if err != nil {
		http.Error(w, "Couldn't get github credentials.", http.StatusInternalServerError)
		return
	}

	resp := utils.FilterRepositoriesByUsername("leonwright", utils.GetUserRepos(githubToken))

	utils.SendJSONResponse(resp, w)
}
