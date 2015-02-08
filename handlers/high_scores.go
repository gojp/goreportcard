package handlers

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gojp/goreportcard/db"
	"gopkg.in/mgo.v2/bson"
)

var highScores []struct {
	Repo    string
	Files   int
	Average float64
}

func HighScoresHandler(w http.ResponseWriter, r *http.Request) {
	db := db.Mongo{URL: mongoURL, Database: mongoDatabase, CollectionName: mongoCollection}
	coll, err := db.Collection()
	if err != nil {
		log.Println("ERROR: could not get collection:", err)
		http.Error(w, err.Error(), 500)
		return
	}

	err = coll.Find(bson.M{"files": bson.M{"$gt": 100}}).Sort("-average").All(&highScores)
	if err != nil {
		log.Println("ERROR: could not get high scores: ", err)
		http.Error(w, err.Error(), 500)
		return
	}

	t, err := template.New("high_scores.html").ParseFiles("templates/high_scores.html")
	if err != nil {
		log.Println("ERROR: ", err)
		http.Error(w, err.Error(), 500)
		return
	}

	t.Execute(w, map[string]interface{}{"HighScores": highScores})
}
