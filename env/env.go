package env

import (
	"github.com/anominet/anomi/cache"
)

const (
	DEFAULT_REDIS_HOST = "172.17.0.26:6379"
	DEFAULT_SEPARATOR  = ":"
)

var DEFAULT_SERIALIZER = cache.JsonSerialzer{}

type Env struct {
	C cache.Cache
}

func (e *Env) Initialize() {
	e.C = &cache.RedisCache{}
	e.C.Dial(DEFAULT_REDIS_HOST)
	e.C.SetSerializer(DEFAULT_SERIALIZER)
	e.C.SetSeparator(DEFAULT_SEPARATOR)
}
