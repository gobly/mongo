package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
)

const tag_name = "bson"
const tag_value_id = "_id"
const tag_value_inline = "inline"

type objectId struct {
	value       bson.ObjectId
	field       reflect.Value
	initialized bool
}

func NewObjectId(v interface{}) *objectId {
	oid := &objectId{}
	mirror := reflect.Indirect(reflect.ValueOf(v))
	oid.scanFields(mirror.Type(), reflect.TypeOf(oid.value), mirror)

	if !oid.initialized {
		panic(`No ID field found! Use tag bson:"_id,omitempty" to define one! ` +
			`If this is a compound object, use tag bson:",inline" to enable scanning of embedded struct.`)
	}

	return oid
}

func (oid *objectId) scanFields(haystack reflect.Type, needle reflect.Type, value reflect.Value) {
	for i := 0; i < haystack.NumField(); i++ {
		field := haystack.Field(i)
		if field.Type != needle {
			val, ok := field.Tag.Lookup(tag_name)
			if ok && strings.Contains(val, tag_value_inline) {
				oid.scanFields(field.Type, needle, value.Field(i))
				continue
			}
		}

		val, ok := field.Tag.Lookup(tag_name)
		if !ok || !strings.Contains(val, tag_value_id) {
			continue
		}

		oid.field = value.Field(i)
		oid.value = oid.field.Interface().(bson.ObjectId)
		oid.initialized = true
	}
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
