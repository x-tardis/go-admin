package sysfiledir

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
)

func GetSysFileDirList(c *gin.Context) {
	var SysFileDir models.SysFileDir
	SysFileDir.Label = c.Request.FormValue("label")
	SysFileDir.PId = cast.ToInt(c.Request.FormValue("pid"))
	SysFileDir.Id = cast.ToInt(c.Request.FormValue("id"))
	SysFileDir.DataScope = jwtauth.UserIdStr(c)
	result, err := SysFileDir.SetSysFileDir()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, result, "")
}

func GetSysFileDir(c *gin.Context) {
	var data models.SysFileDir
	data.Id = cast.ToInt(c.Param("id"))
	result, err := data.Get()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, result, "")
}

// @Summary 添加SysFileDir
// @Description 获取JSON
// @Tags SysFileDir
// @Accept  application/json
// @Product application/json
// @Param data body models.SysFileDir true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/sysfiledir [post]
func InsertSysFileDir(c *gin.Context) {
	var data models.SysFileDir

	if err := c.ShouldBindJSON(&data); err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}
	data.Creator = jwtauth.UserIdStr(c)
	result, err := data.Create()
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, result, "")
}

func UpdateSysFileDir(c *gin.Context) {
	var data models.SysFileDir

	if err := c.ShouldBindJSON(&data); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	data.Updator = jwtauth.UserIdStr(c)
	result, err := data.Update(data.Id)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, result, "")
}

func DeleteSysFileDir(c *gin.Context) {
	var data models.SysFileDir
	data.Updator = jwtauth.UserIdStr(c)

	IDS := infra.ParseIdsGroup(c.Param("ids"))
	_, err := data.BatchDelete(IDS)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.OKWithRequestID(c, nil, codes.DeletedSuccess)
}
