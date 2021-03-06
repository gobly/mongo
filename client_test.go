package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

type HelloWorld struct {
	Id    bson.ObjectId `bson:"_id,omitempty"`
	Name  string        `bson:"name,omitempty" gobly:"slug"`
	Value string        `bson:"value,omitempty"`
}

type Hello struct {
	HelloWorld `bson:",inline"`
}

const collectionName = "testCollection"

func connect(t *testing.T) *Client {
	m := Client{}
	err := m.Connect("localhost", "local", collectionName)

	if err != nil {
		t.Errorf("Could not connect to database: %s", err.Error())
	}

	return &m
}

func TestAnonymousConnection(t *testing.T) {
	connect(t)
}

func TestCreateCollection(t *testing.T) {
	m := connect(t)
	if m == nil {
		return
	}

	err := m.CreateCollection()
	if err != nil {
		t.Errorf("Could not cteate test collection: %s", err.Error())
	}
}

func TestInsertCollection(t *testing.T) {
	m := connect(t)
	if m == nil {
		return
	}

	v := HelloWorld{Name: "world", Value: "Hello"}
	err := m.Insert(&v)
	if err != nil {
		t.Errorf("Could not insert test data to collection: %s", err.Error())
	}
}

func TestInsertInlineCollection(t *testing.T) {
	m := connect(t)
	if m == nil {
		return
	}

	v := Hello{}
	v.Name = "world"
	v.Value = "Hello"
	err := m.Insert(&v)
	if err != nil {
		t.Errorf("Could not insert test data to collection: %s", err.Error())
	}
}

func TestFindDocuments(t *testing.T) {
	m := connect(t)
	if m == nil {
		return
	}

	q := HelloWorld{Name: "world"}
	v := []HelloWorld{}
	err := m.FindByValue(q, &v)
	if err != nil {
		t.Errorf("Could not lookup collection: %s", err.Error())
		return
	}

	if v[0].Name != "world" || v[0].Value != "Hello" {
		t.Errorf("Got back invalid data: name=%s, value=%s", v[0].Name, v[0].Value)
	}
}

func TestReadDocument(t *testing.T) {
	m := connect(t)
	if m == nil {
		return
	}

	q := HelloWorld{Name: "world"}
	err := m.ReadByValue(&q)
	if err != nil {
		t.Errorf("Could not read from collection %s", err.Error())
		return
	}

	if q.Name != "world" || q.Value != "Hello" {
		t.Errorf("Got back invalid data: name=%s, value=%s", q.Name, q.Value)
	}
}

func TestReturnDocumentFromSlug(t *testing.T) {
	m := connect(t)
	if m == nil {
		return
	}

	q := HelloWorld{}
	err := m.ReadBySlug("world", &q)
	if err != nil {
		t.Errorf("Could not read from collection %s", err.Error())
		return
	}

	if q.Name != "world" || q.Value != "Hello" {
		t.Errorf("Got back invalid data: name=%s, value=%s", q.Name, q.Value)
	}
}

func TestDeleteCollection(t *testing.T) {
	m := connect(t)
	if m == nil {
		return
	}

	err := m.DropCollection()
	if err != nil {
		t.Errorf("Could not delete test collection: %s", err.Error())
	}
}
