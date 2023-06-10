package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
)

const tag_name = "bson"
const tag_value_id = "_id"
const tag_value_inline = "inline"

type objectId struct {
	value       primitive.ObjectID
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
		oid.value = oid.field.Interface().(primitive.ObjectID)
		oid.initialized = true
	}
}

func (oid *objectId) Value() (value primitive.ObjectID, valid bool) {
	if primitive.IsValidObjectID(oid.value.Hex()) {
		return oid.value, true
	}

	return primitive.NewObjectID(), false
}

func (oid *objectId) SetValue(value primitive.ObjectID) {
	oid.value = value
	oid.field.Set(reflect.ValueOf(oid.value))
}
