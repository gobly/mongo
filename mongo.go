package mongo

import (
	"gopkg.in/mgo.v2"
	"runtime"
)

type Mongo struct {
	db      *mgo.Database
	c       *mgo.Collection
	session *mgo.Session
}

func (m *Mongo) connect(url string, db string, collection string) error {
	session, err := mgo.Dial(url)

	if err != nil {
		return err
	}

	session.SetMode(mgo.Monotonic, true)

	m.db = session.DB(db)
	m.c = m.db.C(collection)
	m.session = session

	runtime.SetFinalizer(m, func(m *Mongo) { m.session.Close() })
	return nil
}

func (m *Mongo) insert(v interface{}) error {
	return m.c.Insert(v)
}

func (m *Mongo) readByValue(v interface{}) error {
	return m.c.Find(v).One(v)
}

func (m *Mongo) readByID(v interface{}) error {
	return m.c.FindId(v).One(v)
}

func (m *Mongo) findByValue(q interface{}, v interface{}) error {
	return m.c.Find(q).All(v)
}

func (m *Mongo) findById(q interface{}, v interface{}) error {
	return m.c.FindId(q).All(v)
}

func (m *Mongo) create() error {
	return m.c.Create(&mgo.CollectionInfo{})
}

func (m *Mongo) drop() error {
	return m.c.DropCollection()
}
