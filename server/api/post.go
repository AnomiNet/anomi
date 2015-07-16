package api

import (
	"github.com/anominet/anomi/model"
	"github.com/emicklei/go-restful"
	"net/http"
	"strconv"
)

func (e ApiEnv) registerPostApis(c *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/api/posts").
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
		Param(ws.HeaderParameter(e.AuthHeader, "Authorization Token").DataType("string").Required(false)).
		Writes([]model.Post{}))

	ws.Route(ws.GET("/{post-id}").To(e.getPost).
		Doc("Get post").
		Operation("getPost").
		Param(ws.PathParameter("post-id", "identifier of the post").DataType("int64")).
		Param(ws.HeaderParameter(e.AuthHeader, "Authorization Token").DataType("string").Required(false)).
		Writes(model.Post{}))

	ws.Route(ws.GET("/{post-id}/context").To(e.getPostInContext).
		Doc("Get post in context").
		Operation("getPostInContext").
		Param(ws.PathParameter("post-id", "identifier of the post").DataType("int64")).
		Param(ws.HeaderParameter(e.AuthHeader, "Authorization Token").DataType("string").Required(false)).
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

	post.UserHandle, err = e.Model().GetActiveUser(tok)
	if err != nil {
		response.WriteErrorString(400, "No valid user session")
		return
	}

	err = e.Model().CreatePost(&post)
	if err != nil {
		response.WriteErrorString(500, err.Error())
		return
	}
	response.WriteEntity(post)
}

func (e ApiEnv) getPosts(request *restful.Request, response *restful.Response) {
	// TODO view based on query parameter, top view only for now
	posts, err := e.Model().GetTopPosts(DEFAULT_TOP_POST_LIMIT)
	if err != nil {
		response.WriteErrorString(500, err.Error())
		return
	}
	tok := request.HeaderParameter(e.AuthHeader)
	if tok != "" {
		for i := range posts {
			e.Model().PopulateUserVote(&posts[i], tok)
		}
	}
	response.WriteEntity(posts)
}

func (e ApiEnv) getPost(request *restful.Request, response *restful.Response) {
	id, err := strconv.ParseInt(request.PathParameter("post-id"), 10, 64)
	if err != nil {
		e.WriteErrorJsonString(response, http.StatusBadRequest, "The specified post id is not a number")
		return
	}
	post, err := e.Model().GetPostNormalized(id)
	if err != nil {
		e.WriteErrorJsonString(response, http.StatusNotFound, "The specified post does not exist")
		return
	}
	tok := request.HeaderParameter(e.AuthHeader)
	if tok != "" {
		e.Model().PopulateUserVote(post, tok)
	}
	response.WriteEntity(post)
}

func (e ApiEnv) getPostInContext(request *restful.Request, response *restful.Response) {
	id, err := strconv.ParseInt(request.PathParameter("post-id"), 10, 64)
	if err != nil {
		e.WriteErrorJsonString(response, http.StatusBadRequest, "The specified post id is not a number")
		return
	}
	posts, err := e.Model().GetPostInContext(id)
	if err != nil {
		e.WriteErrorJsonString(response, http.StatusNotFound, "The specified post does not exist")
		return
	}
	tok := request.HeaderParameter(e.AuthHeader)
	if tok != "" {
		for i := range posts {
			e.Model().PopulateUserVote(&posts[i], tok)
		}
	}
	response.WriteEntity(posts)
}
