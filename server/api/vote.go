package api

import (
	"github.com/anominet/anomi/model"
	"github.com/emicklei/go-restful"
)

func (e ApiEnv) registerVoteApis(c *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/votes").
		Doc("Vote Management").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("").To(e.createVote).
		Doc("Create a vote").
		Operation("createVote").
		Param(ws.HeaderParameter(e.AuthHeader, "Authorization Token").DataType("string").Required(true)).
		Reads(model.User{}))

	c.Add(ws)
}

func (e ApiEnv) createVote(request *restful.Request, response *restful.Response) {
	tok := request.HeaderParameter(e.AuthHeader)
	user_handle, err := e.Model().GetActiveUser(tok)
	if err != nil {
		response.WriteErrorString(400, "No valid user session")
		return
	}

	vote := model.Vote{}
	err = request.ReadEntity(&vote)
	if err != nil {
		response.WriteErrorString(400, err.Error())
		return
	}

	if vote.Value != 0 && vote.Value != 1 && vote.Value != -1 {
		response.WriteErrorString(400, "Invalid vote vector specified")
		return
	}
	vote.UserHandle = user_handle

	err = e.Model().CreateVote(&vote)
	if err != nil {
		response.WriteErrorString(500, err.Error())
		return
	}

	response.WriteEntity(vote)
}
