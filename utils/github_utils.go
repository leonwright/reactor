package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/leonwright/reactor/models"
)

const apiBase = "https://api.github.com"

// GetUser ...
func GetUser(token string) models.GithubUser {
	URL := apiBase + "/user"
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

	var githubResp models.GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubResp); err != nil {
		log.Fatalln(err)
	}

	return githubResp
}

// GetUserRepos runs th
func GetUserRepos(token string) models.UserRepositories {
	URL := apiBase + "/user/repos"
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

	var githubResp models.UserRepositories
	if err := json.NewDecoder(resp.Body).Decode(&githubResp); err != nil {
		log.Fatal(err)
	}

	return githubResp
}

// FilterRepositoriesByUsername is a method
func FilterRepositoriesByUsername(username string, repos models.UserRepositories) []models.UserRepositories {
	var filtered []models.UserRepositories = []models.UserRepositories{}

	for _, repo := range repos {
		if repo.Owner.Login == username {
			filtered = append(filtered, []models.UserRepository{repo})
		}

	}

	return filtered
}
