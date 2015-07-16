package env

import (
	log "github.com/Sirupsen/logrus"
	"github.com/anominet/anomi/cache"
	"github.com/anominet/anomi/env/internal"
	"os"
)

type Env struct {
	ApiPort     string
	C           cache.Cache
	AuthHeader  string
	Log         internal.Logger
	SwaggerPath string
}

const (
	VERSION             = "0.1"
	DEFAULT_API_PORT    = "8080"
	DEFAULT_REDIS_HOST  = "127.0.0.1"
	DEFAULT_REDIS_PORT  = "6379"
	DEFAULT_SEPARATOR   = ":"
	DEFAULT_AUTH_HEADER = "X-USER-TOKEN"
	REDIS_HOST_ENV_VAR  = "REDIS_PORT_" + DEFAULT_REDIS_PORT + "_TCP_ADDR"
)

var DEFAULT_SERIALIZER = cache.JsonSerialzer{}

func New(redis_host, api_port string, debug bool, swagger_path string) (*Env, error) {
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

	e.ApiPort = api_port
	e.SwaggerPath = swagger_path

	if redis_host == DEFAULT_REDIS_HOST {
		if redis_host = os.Getenv(REDIS_HOST_ENV_VAR); redis_host == "" {
			redis_host = DEFAULT_REDIS_HOST
		}
	}

	e.Log.Debug("[env] Using redis host: " + redis_host)

	e.C = &cache.RedisCache{}

	err := e.C.Dial(redis_host + ":" + DEFAULT_REDIS_PORT)
	if err != nil {
		return &e, err
	}

	e.C.SetSerializer(DEFAULT_SERIALIZER)
	e.C.SetSeparator(DEFAULT_SEPARATOR)
	e.C.SetLogger(e.Log)
	e.AuthHeader = DEFAULT_AUTH_HEADER

	return &e, nil
}
