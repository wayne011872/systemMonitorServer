package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/wayne011872/goSterna/db"
	"github.com/wayne011872/goSterna/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/wayne011872/goSterna/api"
	apiErr "github.com/wayne011872/goSterna/api/err"
	"github.com/wayne011872/systemMonitorServer/model/sysInfo"
	"github.com/wayne011872/systemMonitorServer/input"
)

func NewSysInfoAPI(service string) api.GinAPI {
	return &sysInfoAPI{
		ErrorOutputAPI: api.NewErrorOutputAPI(service),
	}
}

type sysInfoAPI struct {
	api.ErrorOutputAPI
}

func (a *sysInfoAPI) GetName() string{
	return "sysInfo"
}

func (a *sysInfoAPI) GetAPIs() []*api.GinApiHandler {
	return [] *api.GinApiHandler{
		{Path: "/v1/sysInfo",Handler: a.postEndpoint,Method: "POST"},
	}
}

func(a *sysInfoAPI) postEndpoint(c *gin.Context) {
	in := &input.SysInfoInput{}
	err := c.BindJSON(in)
	if err != nil {
		error := apiErr.New(http.StatusBadRequest, err.Error())
		a.GinOutputErr(c, error)
		return
	}
	logger := log.GetLogByReq(c.Request)
	mgoClient := db.GetMgoDBClientByReq(c.Request)
	crud := sysInfo.NewCRUD(c.Request.Context(),mgoClient.GetCoreDB(),logger)
	err = crud.Save(primitive.NilObjectID,in.SysInfo)
	if err != nil {
		a.GinOutputErr(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"result": "ok",
	})
}