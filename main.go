package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
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
	fmt.Println("0.0.1")
	os.Exit(0)
	return nil
}

var (
	app     = kingpin.New("anomi", "The breakdown of bonds between the individual and community")
	help    = app.Flag("help", "Show help.").Short('h').Dispatch(onHelp).Hidden().Bool()
	version = app.Flag("version", "Show application version.").Short('v').Dispatch(onVersion).Bool()
	debug   = app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	port    = app.Flag("port", "Api server port.").Short('p').Default("8080").Int()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	log.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		DisableColors:    false,
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "",
		DisableSorting:   false,
	})
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	e := &env.Env{}
	e.Initialize()
	e.C.SetTypePrefixRegistry(model.TypePrefixRegistry)

	api.StartServer(strconv.Itoa(*port), e)
}
