package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
)

const tag_name = "bson"
const tag_value = "_id"

type objectId struct {
	value bson.ObjectId
	field reflect.Value
}

func NewObjectId(v interface{}) *objectId {
	oid := &objectId{}
	mirror := reflect.Indirect(reflect.ValueOf(v))
	mirrorType := mirror.Type()
	oidType := reflect.TypeOf(oid.value)

	for i := 0; i < mirrorType.NumField(); i++ {
		field := mirrorType.Field(i)
		if field.Type != oidType {
			continue
		}

		val, ok := field.Tag.Lookup(tag_name)
		if !ok || !strings.Contains(val, tag_value) {
			continue
		}

		oid.field = mirror.Field(i)
		oid.value = oid.field.Interface().(bson.ObjectId)
		return oid
	}

	panic(`No ID field found! Use tag bson:"_id,omitempty" to define one!`)
}

func (oid *objectId) Value() (value bson.ObjectId, valid bool) {
	if oid.value.Valid() {
		return oid.value, true
	}

	return bson.NewObjectId(), false
}

func (oid *objectId) SetValue(value bson.ObjectId) {
	oid.value = value
	oid.field.Set(reflect.ValueOf(oid.value))
}
