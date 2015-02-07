package db

import "gopkg.in/mgo.v2"

type Mongo struct {
	URL            string
	Database       string
	CollectionName string
}

func (db *Mongo) Collection() (*mgo.Collection, error) {
	session, err := mgo.Dial(db.URL)
	if err != nil {
		return nil, err
	}
	c := session.DB(db.Database).C(db.CollectionName)

	return c, nil
}
