package mongo

import (
	"context"
	"errors"
	"github.com/gobly/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"runtime"
)

type Client struct {
	db         *mongo.Database
	c          *mongo.Collection
	client     *mongo.Client
	collection string
}

func (m *Client) Connect(url string, db string, collection string) error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		return err
	}

	m.client = client
	m.db = m.client.Database(db)
	m.c = m.db.Collection(collection)
	m.collection = collection

	runtime.SetFinalizer(m, func(m *Client) { m.Close() })
	return nil
}

func (m *Client) Close() {
	_ = m.client.Disconnect(context.TODO())
}

func (m *Client) Insert(v interface{}) error {
	objectId := NewObjectId(v)
	oid, ok := objectId.Value()
	if !ok {
		objectId.SetValue(oid)
	}
	_, err := m.c.InsertOne(context.TODO(), v)
	return err
}

func (m *Client) Update(v interface{}) error {
	objectId := NewObjectId(v)
	oid, ok := objectId.Value()
	if !ok {
		return errors.New("missing objectID")
	}

	_, err := m.c.ReplaceOne(context.TODO(), bson.M{tag_value_id: oid}, v)
	return err
}

func (m *Client) UpdateRawFiltered(f bson.M, q bson.M) error {
	res := m.c.FindOneAndReplace(context.TODO(), f, q)
	return res.Err()
}

func (m *Client) UpdateRaw(q bson.M) error {
	_, err := m.c.ReplaceOne(context.TODO(), nil, q)
	return err
}

func (m *Client) ReadByValueFiltered(f interface{}, v interface{}) error {
	return m.c.FindOne(context.TODO(), f).Decode(v)
}

func (m *Client) ReadByValue(v interface{}) error {
	return m.c.FindOne(context.TODO(), v).Decode(v)
}

func (m *Client) ReadRaw(q bson.M, v interface{}) error {
	return m.c.FindOne(context.TODO(), q).Decode(v)
}

func (m *Client) ReadByID(objectId string, v interface{}) error {
	oid, err := primitive.ObjectIDFromHex(objectId)
	if err != nil {
		return err
	}

	return m.c.FindOne(context.TODO(), bson.M{tag_value_id: oid}).Decode(v)
}

func (m *Client) ReadBySlug(slug string, v interface{}) error {
	oid, err := primitive.ObjectIDFromHex(slug)
	if err == nil {
		return m.c.FindOne(context.TODO(), bson.M{tag_value_id: oid}).Decode(v)
	}

	s := core.NewSlug(v)
	s.SetValue(slug)
	return m.c.FindOne(context.TODO(), v).Decode(v)
}

func (m *Client) FindAll(v interface{}) error {
	res, err := m.c.Find(context.TODO(), bson.D{})
	if err != nil {
		return err
	}
	return res.All(context.TODO(), v)
}

func (m *Client) FindByValue(q interface{}, v interface{}) error {
	res, err := m.c.Find(context.TODO(), q)
	if err != nil {
		return err
	}
	return res.All(context.TODO(), v)
}

func (m *Client) FindByValueSorted(q interface{}, v interface{}, fields ...string) error {
	res, err := m.c.Find(context.TODO(), q, options.Find().SetSort(fields))
	if err != nil {
		return err
	}

	return res.All(context.TODO(), v)
}

func (m *Client) FindById(objectId string, v interface{}) error {
	oid, err := primitive.ObjectIDFromHex(objectId)
	if err != nil {
		return err
	}

	res, err := m.c.Find(context.TODO(), bson.M{tag_value_id: oid})
	if err != nil {
		return err
	}

	return res.All(context.TODO(), v)
}

func (m *Client) FindGroup(q interface{}, groupPipe bson.M, sortPipe bson.M, v interface{}) error {
	pipe := []bson.M{
		{"$match": q},
		{"$group": groupPipe},
	}

	if len(sortPipe) > 0 {
		pipe = append(pipe, bson.M{"$sort": sortPipe})
	}

	res, err := m.c.Aggregate(context.TODO(), pipe)
	if err != nil {
		return err
	}

	return res.All(context.TODO(), v)
}

func (m *Client) FindRedact(q interface{}, redactPipe bson.M, sortPipe bson.M, v interface{}) error {
	pipe := []bson.M{
		{"$match": q},
		{"$redact": bson.M{
			"$cond": []interface{}{redactPipe, "$$KEEP", "$$PRUNE"},
		}},
	}

	if len(sortPipe) > 0 {
		pipe = append(pipe, bson.M{"$sort": sortPipe})
	}

	res, err := m.c.Aggregate(context.TODO(), pipe)
	if err != nil {
		return err
	}

	return res.All(context.TODO(), v)
}

func (m *Client) DeleteById(objectId string) error {
	oid, err := primitive.ObjectIDFromHex(objectId)
	if err != nil {
		return err
	}

	_, err = m.c.DeleteOne(context.TODO(), bson.M{tag_value_id: oid})
	return err
}

func (m *Client) DeleteBySlug(slug string, v interface{}) error {
	oid, err := primitive.ObjectIDFromHex(slug)
	if err == nil {
		_, err := m.c.DeleteOne(context.TODO(), bson.M{tag_value_id: oid})
		return err
	}

	s := core.NewSlug(v)
	s.SetValue(slug)
	_, err = m.c.DeleteOne(context.TODO(), s)
	return err
}

func (m *Client) CreateCollection() error {
	cols, err := m.db.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		return err
	}

	for _, collection := range cols {
		if collection == m.collection {
			return errors.New("collection already exists")
		}
	}

	return m.db.CreateCollection(context.TODO(), m.collection)
}

func (m *Client) DropCollection() error {
	return m.c.Drop(context.TODO())
}
