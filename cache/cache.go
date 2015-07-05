package cache

import (
	"encoding/json"
	"reflect"
)

type Cache interface {
	Dial(addr string) error
	SetSerializer(s Serializer)
	SetSeparator(sep string)
	SetTypePrefixRegistry(r map[reflect.Type]string)
	GetBytes(t reflect.Type, id string) ([]byte, error)
	Get(v interface{}, id string) error
	Set(id string, v interface{}) error
	SelectDb(id int64) error
	FlushDb() error
	Incr(key string) (int64, error)
}

type Serializer interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type JsonSerialzer struct{}

func (JsonSerialzer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func (JsonSerialzer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func GetBaseType(v interface{}) reflect.Type {
	orig := reflect.ValueOf(v)
	switch orig.Kind() {
	case reflect.Ptr:
		orig = reflect.Indirect(orig)
	}
	return orig.Type()
}
