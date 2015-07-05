package cache

import (
	"github.com/garyburd/redigo/redis"
	"reflect"
	"strconv"
)

type RedisCache struct {
	C   redis.Conn
	S   Serializer
	Sep string
	P   map[reflect.Type]string
}

func (r *RedisCache) Dial(addr string) error {
	var err error
	r.C, err = redis.Dial("tcp", addr)
	return err
}

func (r *RedisCache) SetSerializer(s Serializer) {
	r.S = s
}

func (r *RedisCache) SetSeparator(s string) {
	r.Sep = s
}
func (r *RedisCache) SetTypePrefixRegistry(pr map[reflect.Type]string) {
	r.P = pr
}

func (r *RedisCache) GetTypePrefix(t reflect.Type) string {
	if s, ok := r.P[t]; ok {
		return s
	} else {
		return reflect.TypeOf(t).Name()
	}
}

func (r *RedisCache) GetBytes(t reflect.Type, id string) ([]byte, error) {
	b, err := redis.Bytes(r.C.Do("GET", r.GetTypePrefix(t)+r.Sep+id))
	if err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

func (r *RedisCache) Get(v interface{}, id string) error {
	b, err := r.GetBytes(GetBaseType(v), id)

	if err != nil || b == nil {
		return err
	}

	return r.S.Unmarshal(b, v)
}

func (r *RedisCache) Set(id string, v interface{}) error {
	if v == nil {
		panic("can't set to nil")
	}
	b, err := r.S.Marshal(v)
	if err != nil {
		return err
	}

	t := GetBaseType(v)

	_, err = r.C.Do("SET", r.GetTypePrefix(t)+r.Sep+id, b)
	return err
}

func (r *RedisCache) SelectDb(id int64) error {
	_, err := r.C.Do("SELECT", strconv.FormatInt(id, 10))
	return err
}

func (r *RedisCache) FlushDb() error {
	_, err := r.C.Do("FLUSHDB")
	return err
}

func (r *RedisCache) Incr(key string) (int64, error) {
	return redis.Int64(r.C.Do("INCR", key))
}
