package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gojp/goreportcard/db"
	"gopkg.in/mgo.v2/bson"
)

var highScores []struct {
	Repo    string
	Files   int
	Average float64
}

func HighScoresHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := db.Mongo{URL: mongoURL, Database: mongoDatabase, CollectionName: mongoCollection}
	coll, err := db.Collection()
	if err != nil {
		log.Println("ERROR: could not get collection:", err)
		http.Error(w, err.Error(), 500)
		return
	}

	err = coll.Find(bson.M{"files": bson.M{"$gt": 100}}).Sort("average").All(&highScores)
	if err != nil {
		log.Println("ERROR: could not get high scores: ", err)
		http.Error(w, err.Error(), 500)
		return
	}

	b, err := json.Marshal(highScores)
	if err != nil {
		log.Println("ERROR: could not marshal json:", err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(b)
}
