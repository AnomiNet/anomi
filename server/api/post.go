package api

import (
	"github.com/anominet/anomi/model"
	"github.com/emicklei/go-restful"
)

func (e ApiEnv) registerPostApis(c *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/posts").
		Doc("Post Management").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("").To(e.createPost).
		Doc("Create a post").
		Operation("createPost").
		Param(ws.HeaderParameter(e.AuthHeader, "Authorization Token").DataType("string").Required(true)).
		Reads(model.Post{}).
		Writes(model.Post{}))

	ws.Route(ws.GET("").To(e.getPosts).
		Doc("Get posts").
		Operation("getPosts").
		Param(ws.HeaderParameter(e.AuthHeader, "Authorization Token").DataType("string").Required(true)).
		Writes([]model.Post{}))

	c.Add(ws)
}

func (e ApiEnv) createPost(request *restful.Request, response *restful.Response) {
	post := model.Post{}
	err := request.ReadEntity(&post)
	if err != nil {
		response.WriteErrorString(400, err.Error())
		return
	}
	if post.Url == "" && post.Body == "" {
		response.WriteErrorString(400, "Neither url or body specified")
		return
	}

	tok := request.HeaderParameter(e.AuthHeader)
	if tok == "" {
		response.WriteErrorString(400, "No valid user session")
		return
	}

	var userhandle string
	err = e.C.Get(&userhandle, ACTIVE_USERS+"."+tok)
	if err != nil {
		response.WriteErrorString(400, "No valid user session")
		return
	}

	post.UserHandle = userhandle

	err = e.Model().CreatePost(&post)
	if err != nil {
		response.WriteErrorString(500, err.Error())
		return
	}
	response.WriteEntity(post)
}

func (e ApiEnv) getPosts(request *restful.Request, response *restful.Response) {
	p, err := e.Model().GetTopPosts(DEFAULT_TOP_POST_LIMIT)
	if err != nil {
		response.WriteErrorString(500, err.Error())
		return
	}
	response.WriteEntity(p)
}
