package main

import (
	"fmt"
	"github.com/anominet/anomi/env"
	"github.com/anominet/anomi/model"
	"github.com/anominet/anomi/server/api"
	"gopkg.in/alecthomas/kingpin.v1"
	"os"
	"strconv"
)

func onHelp(context *kingpin.ParseContext) error {
	app.Usage(os.Stderr)
	os.Exit(0)
	return nil
}

func onVersion(context *kingpin.ParseContext) error {
	fmt.Println(env.VERSION)
	os.Exit(0)
	return nil
}

var (
	app        = kingpin.New("anomi", "The breakdown of bonds between the individual and community")
	help       = app.Flag("help", "Show help.").Short('h').Dispatch(onHelp).Hidden().Bool()
	version    = app.Flag("version", "Show application version.").Short('v').Dispatch(onVersion).Bool()
	debug      = app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	port       = app.Flag("port", "Api server port.").Short('p').Default(env.DEFAULT_API_PORT).Int()
	redis_host = app.Flag("redis", "Redis server host.").Short('r').Default(env.DEFAULT_REDIS_HOST).String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	e := env.New(*redis_host, *debug)
	e.C.SetTypePrefixRegistry(model.TypePrefixRegistry)

	api.StartServer(strconv.Itoa(*port), e)
}
