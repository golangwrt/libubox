package libubox

import (
	"testing"
	"encoding/json"
)

func TestJSONObject(t *testing.T) {
	var param struct {
		Name string `json:"name"`
		Age int `json:"age"`
	}

	obj, err := NewJSONObject(&param)
	if err != nil {
		t.Fatalf("new json object fail: %s\n", err)
	}
	data := json.RawMessage(`{"name":"test", "age":22}`)
	err = obj.UnmarshalJSON([]byte(data))
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	t.Logf("val is %+v\n", obj.Value)
}
