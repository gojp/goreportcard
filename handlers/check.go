package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CheckHandler handles the request for checking a repo
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	repo := r.FormValue("repo")
	log.Printf("Checking repo %s...", repo)
	if strings.ToLower(repo) == "golang/go" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("We've decided to omit results for the Go repository because it has lots of test files that (purposely) don't pass our checks. Go gets an A+ in our books though!"))
		return
	}
	forceRefresh := r.Method != "GET" // if this is a GET request, try fetch from cached version in mongo first
	resp, err := newChecksResp(repo, forceRefresh)
	if err != nil {
		log.Println("ERROR: ", err)
		b, _ := json.Marshal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(b)
		return
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println("ERROR: could not marshal json:", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(b)

	// write to mongo
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Println("Failed to get mongo collection: ", err)
		return
	}
	defer session.Close()
	coll := session.DB(mongoDatabase).C(mongoCollection)
	log.Printf("Upserting repo %s...", repo)
	_, err = coll.Upsert(bson.M{"repo": repo}, resp)
	if err != nil {
		log.Println("Mongo writing error:", err)
		return
	}
}
