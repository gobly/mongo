package mongo

import (
	"errors"
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

func (m *Client) Close() {
	m.db.Logout()
	m.session.Close()
}

func (m *Client) Insert(v interface{}) error {
	objectId := NewObjectId(v)
	oid, ok := objectId.Value()
	if !ok {
		objectId.SetValue(oid)
	}
	err := m.c.Insert(v)
	return err
}

func (m *Client) Update(v interface{}) error {
	objectId := NewObjectId(v)
	oid, ok := objectId.Value()
	if !ok {
		return errors.New("Missing objectID")
	}

	return m.c.UpdateId(oid, v)
}

func (m *Client) UpdateRawFiltered(f bson.M, q bson.M) error {
	return m.c.Update(f, q)
}

func (m *Client) UpdateRaw(q bson.M) error {
	return m.c.Update(nil, q)
}

func (m *Client) ReadByValueFiltered(f interface{}, v interface{}) error {
	return m.c.Find(v).Select(f).One(v)
}

func (m *Client) ReadByValue(v interface{}) error {
	return m.c.Find(v).One(v)
}

func (m *Client) ReadRaw(q bson.M, v interface{}) error {
	return m.c.Find(q).One(v)
}

func (m *Client) ReadRawFiltered(q bson.M, f bson.M, v interface{}) error {
	return m.c.Find(q).Select(f).One(v)
}

func (m *Client) ReadByID(objectId string, v interface{}) error {
	if !bson.IsObjectIdHex(objectId) {
		return errors.New("Invalid ObjectID format")
	}

	return m.c.FindId(bson.ObjectIdHex(objectId)).One(v)
}

func (m *Client) ReadBySlug(slug string, v interface{}) error {
	if bson.IsObjectIdHex(slug) {
		return m.c.FindId(bson.ObjectIdHex(slug)).One(v)
	}

	s := core.NewSlug(v)
	s.SetValue(slug)
	return m.c.Find(v).One(v)
}

func (m *Client) FindAll(v interface{}) error {
	return m.c.Find(nil).All(v)
}

func (m *Client) FindByValue(q interface{}, v interface{}) error {
	return m.c.Find(q).All(v)
}

func (m *Client) FindByValueSorted(q interface{}, v interface{}, fields ...string) error {
	return m.c.Find(q).Sort(fields...).All(v)
}

func (m *Client) FindById(objectId string, v interface{}) error {
	if !bson.IsObjectIdHex(objectId) {
		return errors.New("Invalid ObjectID format")
	}

	return m.c.FindId(bson.ObjectIdHex(objectId)).All(v)
}

func (m *Client) FindGroup(q interface{}, groupPipe bson.M, sortPipe bson.M, v interface{}) error {
	pipe := []bson.M{
		{"$match": q},
		{"$group": groupPipe},
	}

	if len(sortPipe) > 0 {
		return m.c.Pipe(append(pipe, bson.M{"$sort": sortPipe})).All(v)
	}

	return m.c.Pipe(pipe).All(v)
}

func (m *Client) FindRedact(q interface{}, redactPipe bson.M, sortPipe bson.M, v interface{}) error {
	pipe := []bson.M{
		{"$match": q},
		{"$redact": bson.M{
			"$cond": []interface{}{redactPipe, "$$KEEP", "$$PRUNE"},
		}},
	}

	if len(sortPipe) > 0 {
		return m.c.Pipe(append(pipe, bson.M{"$sort": sortPipe})).All(v)
	}

	return m.c.Pipe(pipe).All(v)
}

func (m *Client) DeleteById(objectId string) error {
	if !bson.IsObjectIdHex(objectId) {
		return errors.New("Invalid ObjectID format")
	}

	return m.c.RemoveId(bson.ObjectIdHex(objectId))
}

func (m *Client) DeleteBySlug(slug string, v interface{}) error {
	if bson.IsObjectIdHex(slug) {
		return m.c.RemoveId(bson.ObjectIdHex(slug))
	}

	s := core.NewSlug(v)
	s.SetValue(slug)
	return m.c.Remove(v)
}

func (m *Client) CreateCollection() error {
	cols, err := m.db.CollectionNames()
	if err != nil {
		return err
	}

	for _, collection := range cols {
		if collection == m.collection {
			return errors.New("Collection already exists")
		}
	}

	return m.c.Create(&mgo.CollectionInfo{})
}

func (m *Client) DropCollection() error {
	return m.c.DropCollection()
}
