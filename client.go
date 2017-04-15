package mongo

import (
	"github.com/gobly/core"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"runtime"
)

type Client struct {
	db         *mgo.Database
	c          *mgo.Collection
	session    *mgo.Session
	collection string
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
	m.collection = collection

	runtime.SetFinalizer(m, func(m *Client) { m.session.Close() })
	return nil
}

func (m *Client) Insert(v interface{}) error {
	objectId := ObjectId{}
	oid, ok := objectId.GetMirror(v)
	if !ok {
		objectId.SetOid(oid)
	}
	err := m.c.Insert(v)
	return err
}

func (m *Client) ReadByValue(v interface{}) error {
	return m.c.Find(v).One(v)
}

func (m *Client) ReadByID(v interface{}) error {
	return m.c.FindId(v).One(v)
}

func (m *Client) ReadBySlug(slug string, v interface{}) error {
	if bson.IsObjectIdHex(slug) {
		return m.c.FindId(bson.M{"_id": bson.ObjectIdHex(slug)}).One(v)
	}

	s := core.NewSlug(v)
	s.SetValue(slug)
	return m.c.Find(v).One(v)
}

func (m *Client) FindByValue(q interface{}, v interface{}) error {
	return m.c.Find(q).All(v)
}

func (m *Client) FindById(q interface{}, v interface{}) error {
	return m.c.FindId(q).All(v)
}

func (m *Client) CreateCollection() error {
	cols, err := m.db.CollectionNames()
	if err != nil {
		return err
	}

	for _, collection := range cols {
		if collection == m.collection {
			return nil
		}
	}

	return m.c.Create(&mgo.CollectionInfo{})
}

func (m *Client) DropCollection() error {
	return m.c.DropCollection()
}
