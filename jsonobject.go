package libubox

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// JsonObject
type JSONObject struct {
	Value interface {}
	// underline struct type
	T reflect.Type
}
// NewJSONObject create a new json object, then
// bind the type/value with i
func NewJSONObject(i interface{}) (*JSONObject, error) {
	obj, err := NewJSONObjectWith(reflect.TypeOf(i))
	if err != nil {
		return nil, err
	}
	obj.Value = i
	return obj, nil
}

// NewJSONObjectWith create a new json object, then
// bind the object with concrete struct type of i
func NewJSONObjectWith(t reflect.Type) (*JSONObject, error) {
	if t == nil {
		return nil, fmt.Errorf("nil type")
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("parameter type is %s, require struct", t.Kind().String())
	}
	return &JSONObject{
		T: t,
	}, nil
}

func (obj *JSONObject) UnmarshalBlobAttr(attr *BlobAttr) error  {
	return obj.UnmarshalJSON([]byte(attr.FormatJSON(true)))
}

func (obj *JSONObject) UnmarshalJSON(data []byte) error  {
	val := reflect.New(obj.T)
	err := json.Unmarshal(data, val.Interface())
	if err != nil {
		return err
	}
	obj.Value = val.Interface()
	return nil
}

func (obj *JSONObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(obj.Value)
}
