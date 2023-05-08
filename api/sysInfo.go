package api

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/wayne011872/goSterna/db"
	"github.com/wayne011872/goSterna/log"
	"github.com/wayne011872/goSterna/api"
	apiErr "github.com/wayne011872/goSterna/api/err"
	"github.com/wayne011872/systemMonitorServer/model/sysInfo"
	"github.com/wayne011872/systemMonitorServer/input"
	"github.com/wayne011872/systemMonitorServer/libs"
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
	isError := libs.DetectError(in.SysInfo)
	if isError {
		in.SysInfo.SendTime = time.Now().Format("2006-01-02 15:04:05")
	}
	logger := log.GetLogByGin(c)
	mgoClient := db.GetMgoDBClientByGin(c)
	crud := sysInfo.NewCRUD(c.Request.Context(),mgoClient.GetCoreDB(),logger)
	_,err = crud.Save(in.SysInfo)
	if err != nil {
		a.GinOutputErr(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"result": "ok",
	})
}