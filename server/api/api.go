package api

import (
	"encoding/json"
	"fmt"
	"github.com/anominet/anomi/env"
	"github.com/anominet/anomi/model"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"net/http"
	"strings"
)

const (
	DEFAULT_TOP_POST_LIMIT   = 100
	DEFAULT_FORWARDED_HEADER = "X-Forwarded-For"
)

type ApiEnv struct {
	*env.Env
}

func (e ApiEnv) Model() model.ModelEnv {
	return model.ModelEnv{e.Env}
}

func (e ApiEnv) ReqLogger(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	chain.ProcessFilter(req, resp)

	var raddr string
	if list, ok := req.Request.Header[DEFAULT_FORWARDED_HEADER]; ok {
		// The last ip will always be the one added by haproxy
		raddr = list[len(list)-1]
	}
	if raddr == "" {
		raddr = strings.Split(req.Request.RemoteAddr, ":")[0]
	}

	e.Log.Info(fmt.Sprintf(
		"[server/api] %s %s %s %s %d",
		raddr,
		req.Request.Method,
		req.Request.URL.RequestURI(),
		req.Request.Header.Get("Content-Type"),
		resp.StatusCode(),
	))

	if e.Log.Level.String() == "debug" {
		content, err := json.Marshal(req.Request.Header)
		if err != nil {
			return
		}
		e.Log.Debug("[server/api] " + string(content))

		var temp interface{}
		err = req.ReadEntity(&temp)
		if err != nil {
			return
		}
		content, err = json.Marshal(temp)
		if err != nil {
			return
		}
		e.Log.Debug("[server/api] " + string(content))
	}
}

func (e ApiEnv) WriteServiceErrorJson(err restful.ServiceError, req *restful.Request, resp *restful.Response) {
	e.Log.Error(err)
	e.WriteErrorJsonString(resp, err.Code, err.Message)
}

func (e ApiEnv) WriteErrorJsonString(r *restful.Response, httpStatus int, err string) error {
	r.WriteHeader(httpStatus)
	return r.WriteAsJson(map[string][]string{"errors": []string{err}})
}

func (e ApiEnv) WriteErrorJson(r *restful.Response, httpStatus int, err error) error {
	return e.WriteErrorJsonString(r, httpStatus, err.Error())
}

func StartServer(e *env.Env) {
	restful.SetLogger(e.Log)
	aEnv := ApiEnv{e}

	wsContainer := restful.NewContainer()

	// Enable gzip encoding
	//wsContainer.EnableContentEncoding(true)

	// Request logging
	wsContainer.Filter(aEnv.ReqLogger)

	// Route error handling
	wsContainer.ServiceErrorHandler(aEnv.WriteServiceErrorJson)

	// Register apis
	aEnv.registerUserApis(wsContainer)
	aEnv.registerPostApis(wsContainer)
	aEnv.registerVoteApis(wsContainer)

	// Some bullshit so ServiceErrorHandler works...
	errWs := new(restful.WebService)
	errWs.Path("/")
	wsContainer.Add(errWs)

	if e.SwaggerPath != "" {
		config := swagger.Config{
			WebServices:     wsContainer.RegisteredWebServices(),
			ApiPath:         "/swagger/apidocs.json",
			SwaggerPath:     "/swagger/apidocs/",
			SwaggerFilePath: e.SwaggerPath,
		}

		//Container just for swagger
		swContainer := restful.NewContainer()
		swagger.RegisterSwaggerService(config, swContainer)
		http.Handle("/swagger/", swContainer)
	}

	http.Handle("/", wsContainer)
	e.Log.Info("[server/api] start listening on localhost:" + e.ApiPort)

	e.Log.Fatal(http.ListenAndServe(":"+e.ApiPort, nil))
}
