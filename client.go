package mongo

import (
	"gopkg.in/mgo.v2"
	"runtime"
)

type Client struct {
	db      *mgo.Database
	c       *mgo.Collection
	session *mgo.Session
}

func (m *Client) Connect(url string, db string, collection string) error {
	session, err := mgo.Dial(url)

	if err != nil {
		return err
	}

	session.SetMode(mgo.Monotonic, true)

	m.db = session.DB(db)
	m.c = m.db.C(collection)
	m.session = session

	runtime.SetFinalizer(m, func(m *Client) { m.session.Close() })
	return nil
}

func (m *Client) Insert(v interface{}) error {
	return m.c.Insert(v)
}

func (m *Client) ReadByValue(v interface{}) error {
	return m.c.Find(v).One(v)
}

func (m *Client) ReadByID(v interface{}) error {
	return m.c.FindId(v).One(v)
}

func (m *Client) FindByValue(q interface{}, v interface{}) error {
	return m.c.Find(q).All(v)
}

func (m *Client) FindById(q interface{}, v interface{}) error {
	return m.c.FindId(q).All(v)
}

func (m *Client) CreateCollection() error {
	return m.c.Create(&mgo.CollectionInfo{})
}

func (m *Client) DropCollection() error {
	return m.c.DropCollection()
}
