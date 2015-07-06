package env

import (
	"github.com/anominet/anomi/cache"
)

const (
	DEFAULT_REDIS_HOST  = "172.17.0.26:6379"
	DEFAULT_SEPARATOR   = ":"
	DEFAULT_AUTH_HEADER = "HTTP_X_USER_TOKEN"
)

var DEFAULT_SERIALIZER = cache.JsonSerialzer{}

type Env struct {
	C          cache.Cache
	AuthHeader string
}

func (e *Env) Initialize() {
	e.C = &cache.RedisCache{}
	e.C.Dial(DEFAULT_REDIS_HOST)
	e.C.SetSerializer(DEFAULT_SERIALIZER)
	e.C.SetSeparator(DEFAULT_SEPARATOR)
	e.AuthHeader = DEFAULT_AUTH_HEADER
}
