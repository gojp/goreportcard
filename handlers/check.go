package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gojp/goreportcard/db"
	"gopkg.in/mgo.v2/bson"
)

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	repo := r.FormValue("repo")
	forceRefresh := r.Method != "GET" // if this is a GET request, try fetch from cached version in mongo first
	resp, err := newChecksResp(repo, forceRefresh)
	if err != nil {
		log.Println("ERROR: ", err)
		b, _ := json.Marshal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(b)
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println("ERROR: could not marshal json:", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(b)

	// write to mongo
	db := db.Mongo{URL: mongoURL, Database: mongoDatabase, CollectionName: mongoCollection}
	coll, err := db.Collection()
	if err != nil {
		log.Println("Failed to get mongo collection: ", err)
	} else {
		log.Println("Writing to mongo...")
		_, err := coll.Upsert(bson.M{"repo": repo}, resp)
		if err != nil {
			log.Println("Mongo writing error:", err)
		}
	}
}
