package mongo

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Recorder struct {
	session    *mgo.Session
	set        string
	collection string
}

func NewRecorder(host string, set string, collection string) (*Recorder, error) {
	sess, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}

	return &Recorder{
		session:    sess,
		set:        set,
		collection: collection,
	}
}

func (r *Recorder) Insert(args interface{}) error {
	dbHandler := r.session.DB(r.set).C(r.collection)
	err := dbHandler.Insert(args)
	if err != nil {
		return err
	}

	return nil
}

func (r *Recorder) FindLatest(deviceid uint64, records interface{}) error {
	dbHandler := r.session.DB(r.set).C(r.collection)
	err := dbHandler.Find(bson.M{
		"$query":   bson.M{"deviceid": deviceid},
		"$orderby": bson.M{"timestamp": -1},
	}).Limit(1).One(records)

	return nil
}

func (r *Recorder) FindByTimestamp(deviceid uint64, start uint64, end uint64, records interface{}) error {
	dbHandler := r.session.DB(r.set).C(r.collection)
	err := dbHandler.Find(bson.M{
		"$query":    bson.M{"deviceid": deviceid},
		"timestamp": bson.M{"$gte": start, "$lte": end},
	}).All(records)

	return nil
}
