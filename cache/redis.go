package cache

import (
	"github.com/anominet/anomi/env/internal"
	"github.com/garyburd/redigo/redis"
	"reflect"
	"strconv"
)

type RedisCache struct {
	C   redis.Conn
	S   Serializer
	Sep string
	P   map[reflect.Type]string
	Log internal.Logger
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

func (r *RedisCache) SetLogger(logger internal.Logger) {
	r.Log = logger
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

func (r *RedisCache) GetListBytes(t reflect.Type, id string) ([]byte, error) {
	b, err := redis.Bytes(r.C.Do("GET", "[]"+r.GetTypePrefix(t)+r.Sep+id))
	if err != nil {
		return nil, err
	} else {
		// FIXME
		s := append([]byte("["), b[:len(b)-1]...)
		s = append(s, []byte("]")...)
		return s, nil
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

func (r *RedisCache) GetList(v interface{}, id string) error {
	b, err := r.GetListBytes(GetBaseType(v), id)

	if err != nil || b == nil {
		if err == redis.ErrNil {
			return nil
		} else {
			return err
		}
	}

	return r.S.Unmarshal(b, v)
}

func (r *RedisCache) Append(id string, v interface{}) error {
	if v == nil {
		panic("can't append nil")
	}
	b, err := r.S.Marshal(v)
	if err != nil {
		return err
	}
	b = append(b, []byte(",")...)

	t := GetBaseType(v)
	_, err = r.C.Do("APPEND", "[]"+r.GetTypePrefix(t)+r.Sep+id, b)
	return err
}

func (r *RedisCache) ZAdd(set string, score int64, v interface{}) error {
	if v == nil {
		panic("can't add nil")
	}
	b, err := r.S.Marshal(v)
	if err != nil {
		return err
	}

	t := GetBaseType(v)
	_, err = r.C.Do("ZADD", "set"+r.Sep+r.GetTypePrefix(t)+r.Sep+set, score, b)
	return err

}

func (r *RedisCache) ZScore(set string, v interface{}) (int64, error) {
	if v == nil {
		panic("can't get score of nil")
	}
	b, err := r.S.Marshal(v)
	if err != nil {
		return 0, err
	}

	t := GetBaseType(v)

	return redis.Int64(r.C.Do("ZSCORE", "set"+r.Sep+r.GetTypePrefix(t)+r.Sep+set, b))
}

func (r *RedisCache) ZRangeByScore(v interface{}, set string, dir bool, limit int64) ([]int64, error) {

	var values []interface{}
	var err error
	t := GetBaseType(v)
	// Using LIMIT offset for pagination has bad performance
	//  so always use 0 offset and implement elsewhere
	if dir == HIGH_TO_LOW {
		values, err = redis.Values(r.C.Do(
			"ZREVRANGEBYSCORE",
			"set"+r.Sep+r.GetTypePrefix(t)+r.Sep+set,
			"+inf",
			"-inf",
			"WITHSCORES",
			"LIMIT",
			"0",
			strconv.FormatInt(limit, 10)))
	} else {
		// LOW_TO_HIGH
		values, err = redis.Values(r.C.Do(
			"ZRANGEBYSCORE",
			"set"+r.Sep+r.GetTypePrefix(t)+r.Sep+set,
			"-inf",
			"+inf",
			"WITHSCORES",
			"LIMIT",
			"0",
			strconv.FormatInt(limit, 10)))
	}

	if err != nil {
		return nil, err
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Ptr:
		val = reflect.Indirect(val)
	}

	reply_len := len(values) / 2
	list := reflect.MakeSlice(val.Type(), reply_len, reply_len)
	scores := make([]int64, reply_len)

	for i := 0; i < reply_len; i++ {
		scores[i], err = redis.Int64(values[(i*2)+1], nil)
		if err != nil {
			return nil, err
		}
		b, err := redis.Bytes(values[i*2], nil)
		if err != nil {
			return nil, err
		}
		r.S.Unmarshal(b, list.Index(i).Addr().Interface())
	}
	val.Set(list)
	return scores, nil
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
