package api

import (
	"github.com/anominet/anomi/model"
	"github.com/emicklei/go-restful"
)

func (e ApiEnv) registerUserApis(c *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/users").
		Doc("User Management").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("").To(e.createUser).
		Doc("Create a user").
		Operation("createUser").
		Reads(model.User{}))

	c.Add(ws)
}

func (e ApiEnv) createUser(request *restful.Request, response *restful.Response) {
	usr := model.User{}
	err := request.ReadEntity(&usr)
	if err != nil {
		response.WriteErrorString(500, err.Error())
		return
	}
	if usr.Handle == "" {
		response.WriteErrorString(400, "No handle specified")
		return
	}
	err = e.ModelEnv().CreateUser(&usr)
	if err != nil {
		response.WriteErrorString(500, err.Error())
		return
	}
	response.WriteEntity(usr)
}
