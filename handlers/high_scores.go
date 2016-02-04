package handlers

import (
	"fmt"
	"net/http"
)

func add(x, y int) int {
	return x + y
}

func formatScore(x float64) string {
	return fmt.Sprintf("%.2f", x*100)
}

// HighScoresHandler handles the stats page
func HighScoresHandler(w http.ResponseWriter, r *http.Request) {
	// session, err := mgo.Dial(mongoURL)
	// if err != nil {
	// 	log.Println("ERROR: could not get collection:", err)
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }
	// defer session.Close()
	// coll := session.DB(mongoDatabase).C(mongoCollection)
	//
	// var highScores []struct {
	// 	Repo    string
	// 	Files   int
	// 	Average float64
	// }
	// err = coll.Find(bson.M{"files": bson.M{"$gt": 100}}).Sort("-average").Limit(50).All(&highScores)
	// if err != nil {
	// 	log.Println("ERROR: could not get high scores: ", err)
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }
	//
	// count, err := coll.Count()
	// if err != nil {
	// 	log.Println("ERROR: could not get count of high scores: ", err)
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }
	//
	// funcs := template.FuncMap{"add": add, "formatScore": formatScore}
	// t := template.Must(template.New("high_scores.html").Funcs(funcs).ParseFiles("templates/high_scores.html"))
	//
	// t.Execute(w, map[string]interface{}{"HighScores": highScores, "Count": humanize.Comma(int64(count))})
}
