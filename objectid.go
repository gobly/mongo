package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
)

type ObjectId struct {
	objectID bson.ObjectId
	idField  reflect.Value
}

func (oid *ObjectId) GetMirror(v interface{}) (objectId bson.ObjectId, ok bool) {
	mirror := reflect.Indirect(reflect.ValueOf(v))
	mirrorType := mirror.Type()
	oidType := reflect.TypeOf(oid.objectID)

	for i := 0; i < mirrorType.NumField(); i++ {
		field := mirrorType.Field(i)
		if field.Type != oidType {
			continue
		}

		val, ok := field.Tag.Lookup("bson")
		if !ok || !strings.Contains(val, "_id") {
			continue
		}

		oid.idField = mirror.Field(i)
		oid.objectID = oid.idField.Interface().(bson.ObjectId)

		if oid.objectID.Valid() {
			return oid.objectID, true
		}

		return bson.NewObjectId(), false
	}

	panic(`No ID field found! Use tag bson:"_id,omitempty" to define one!`)
}

func (oid *ObjectId) SetOid(id bson.ObjectId) bson.ObjectId {
	oid.objectID = id
	oid.idField.Set(reflect.ValueOf(oid.objectID))
	return oid.objectID
}
