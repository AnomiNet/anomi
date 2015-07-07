package cache

import (
	"github.com/anominet/anomi/env/internal"
	"encoding/json"
	"reflect"
)

type Cache interface {
	Dial(addr string) error
	SetSerializer(s Serializer)
	SetSeparator(sep string)
	SetTypePrefixRegistry(r map[reflect.Type]string)
	SetLogger(logger internal.Logger)
	GetBytes(t reflect.Type, id string) ([]byte, error)
	Get(v interface{}, id string) error
	Set(id string, v interface{}) error
	GetList(v interface{}, id string) error
	Append(id string, v interface{}) error
	ZAdd(set string, score int64, v interface{}) error
	ZScore(set string, v interface{}) (int64, error)
	ZRangeByScore(v interface{}, set string, dir bool, limit int64) ([]int64, error)
	SelectDb(id int64) error
	FlushDb() error
	Incr(key string) (int64, error)
}

type Serializer interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type JsonSerialzer struct{}

const (
	HIGH_TO_LOW = true
	LOW_TO_HIGH = false
)

func (JsonSerialzer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func (JsonSerialzer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func GetBaseType(v interface{}) reflect.Type {
	orig := reflect.ValueOf(v)
	for {
		switch orig.Kind() {
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			return orig.Type().Elem()
		case reflect.Ptr:
			orig = reflect.Indirect(orig)
		default:
			return orig.Type()
		}
	}
}
