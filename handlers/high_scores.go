package handlers

import (
	"log"
	"net/http"
	"text/template"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func add(x, y int) int {
	return x + y
}

// HighScoresHandler handles the stats page
func HighScoresHandler(w http.ResponseWriter, r *http.Request) {
	var highScores []struct {
		Repo    string
		Files   int
		Average float64
	}

	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Println("ERROR: could not get collection:", err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer session.Close()
	coll := session.DB(mongoDatabase).C(mongoCollection)

	err = coll.Find(bson.M{"files": bson.M{"$gt": 100}}).Sort("-average").All(&highScores)
	if err != nil {
		log.Println("ERROR: could not get high scores: ", err)
		http.Error(w, err.Error(), 500)
		return
	}

	count, err := coll.Count()
	if err != nil {
		log.Println("ERROR: could not get high scores: ", err)
		http.Error(w, err.Error(), 500)
		return
	}

	funcs := template.FuncMap{"add": add}
	t := template.Must(template.New("high_scores.html").Funcs(funcs).ParseFiles("templates/high_scores.html"))

	t.Execute(w, map[string]interface{}{"HighScores": highScores, "Count": count})
}
