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

	sess.DB(set).C(collection).EnsureIndexKey("deviceid", "timestamp")

	return &Recorder{
		session:    sess,
		set:        set,
		collection: collection,
	}, nil
}

func (r *Recorder) Insert(args interface{}) error {
	dbHandler := r.session.DB(r.set).C(r.collection)

	err := dbHandler.Insert(args)
	if err != nil {
		return err
	}

	return nil
}

func (r *Recorder) FindLatest(deviceid uint64, record interface{}) error {
	dbHandler := r.session.DB(r.set).C(r.collection)
	err := dbHandler.Find(bson.M{
		"$query":   bson.M{"deviceid": deviceid},
		"$orderby": bson.M{"timestamp": -1},
	}).Limit(1).One(record)

	return err
}

func (r *Recorder) FindByTimestamp(deviceid uint64, start uint64, end uint64, records interface{}) error {
	dbHandler := r.session.DB(r.set).C(r.collection)
	err := dbHandler.Find(bson.M{
		"deviceid":  deviceid,
		"timestamp": bson.M{"$gte": start, "$lte": end},
	}).All(records)

	return err
}
