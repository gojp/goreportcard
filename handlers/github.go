package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type githubRepository struct {
	FullName string `json:"full_name"`
}

type githubPayload struct {
	Repository githubRepository `json:"repository"`
}

// GithubHookHandler handles "application/json" POST request from github hooks.
func GithubHookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Fail to read request body: %s", err), http.StatusBadRequest)
		return
	}

	var payload githubPayload
	err = json.Unmarshal(b, &payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("fail to unmarshal request body: %s", err), http.StatusBadRequest)
		return
	}

	repo := "github.com/" + payload.Repository.FullName
	handleHTTPCheckRepo(repo, true, w)
}
