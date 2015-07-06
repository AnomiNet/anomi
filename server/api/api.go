package api

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/anominet/anomi/env"
	"github.com/anominet/anomi/model"
	"github.com/emicklei/go-restful"
	//"github.com/emicklei/go-restful/swagger"
	"net/http"
	"strings"
)

const ACTIVE_USERS = "active.users"
const DEFAULT_TOP_POST_LIMIT = 100

type ApiEnv struct {
	*env.Env
}

func (e ApiEnv) Model() model.ModelEnv {
	return model.ModelEnv{e.Env}
}

func ReqLogger(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	chain.ProcessFilter(req, resp)

	log.Info(fmt.Sprintf(
		"[anomi/servers/api] %s %s %s %s %d",
		strings.Split(req.Request.RemoteAddr, ":")[0],
		req.Request.Method,
		req.Request.URL.RequestURI(),
		req.Request.Header.Get("Content-Type"),
		resp.StatusCode(),
	))

	var temp interface{}
	err := req.ReadEntity(&temp)
	if err != nil {
		return
	}
	content, err := json.Marshal(temp)
	if err != nil {
		return
	}
	log.Debug("[chromaticity/servers/api] " + string(content))
}

func StartServer(port string, e *env.Env) {
	// FIXME, logger from env?
	restful.SetLogger(log.StandardLogger())

	wsContainer := restful.NewContainer()
	// Enable gzip encoding
	//wsContainer.EnableContentEncoding(true)
	wsContainer.Filter(ReqLogger)

	aenv := ApiEnv{e}

	// Register apis
	aenv.registerUserApis(wsContainer)
	aenv.registerPostApis(wsContainer)

	// Uncomment to add some swagger
	//config := swagger.Config{
	//WebServices:    wsContainer.RegisteredWebServices(),
	//WebServicesUrl: "/",
	//ApiPath:        "/swagger/apidocs.json",
	//SwaggerPath:    "/swagger/apidocs/",
	//}

	//Container just for swagger
	//swContainer := restful.NewContainer()
	//swagger.RegisterSwaggerService(config, swContainer)
	//http.Handle("/swagger/", swContainer)
	//http.Handle("/apidocs/", &AssetHandler{})

	// FIXME
	//http.Handle("/api", wsContainer)
	//http.Handle("/api/", wsContainer)

	log.Info("[anomi/servers/api] start listening on localhost:" + port)
	//log.Fatal(http.ListenAndServe(":"+port, nil))

	server := &http.Server{Addr: ":" + port, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
