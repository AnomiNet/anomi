package env

import (
	log "github.com/Sirupsen/logrus"
	"github.com/anominet/anomi/cache"
	"github.com/anominet/anomi/env/internal"
)

const (
	DEFAULT_REDIS_HOST  = "172.17.0.26:6379"
	DEFAULT_SEPARATOR   = ":"
	DEFAULT_AUTH_HEADER = "X-USER-TOKEN"
)

var DEFAULT_SERIALIZER = cache.JsonSerialzer{}

type Env struct {
	C          cache.Cache
	AuthHeader string
	Log        internal.Logger
}

func New(debug bool) *Env {
	e := Env{}
	e.Log = internal.Logger{log.New()}
	e.Log.Formatter = &log.TextFormatter{
		ForceColors:      true,
		DisableColors:    false,
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "",
		DisableSorting:   false,
	}
	if debug {
		e.Log.Level = log.DebugLevel
	}

	e.C = &cache.RedisCache{}
	e.C.Dial(DEFAULT_REDIS_HOST)
	e.C.SetSerializer(DEFAULT_SERIALIZER)
	e.C.SetSeparator(DEFAULT_SEPARATOR)
	e.C.SetLogger(e.Log)
	e.AuthHeader = DEFAULT_AUTH_HEADER
	return &e
}
