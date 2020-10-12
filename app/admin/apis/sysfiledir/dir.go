package sysfiledir

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/cast"

	"github.com/x-tardis/go-admin/app/admin/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/tools"
)

func GetSysFileDirList(c *gin.Context) {
	var SysFileDir models.SysFileDir
	SysFileDir.Label = c.Request.FormValue("label")
	SysFileDir.PId = cast.ToInt(c.Request.FormValue("pid"))
	SysFileDir.Id = cast.ToInt(c.Request.FormValue("id"))
	SysFileDir.DataScope = jwtauth.UserIdStr(c)
	result, err := SysFileDir.SetSysFileDir()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	servers.OKWithRequestID(c, result, "")
}

func GetSysFileDir(c *gin.Context) {
	var data models.SysFileDir
	data.Id = cast.ToInt(c.Param("id"))
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)

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
	err := c.ShouldBindJSON(&data)
	data.CreateBy = jwtauth.UserIdStr(c)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	servers.OKWithRequestID(c, result, "")
}

func UpdateSysFileDir(c *gin.Context) {
	var data models.SysFileDir
	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	data.UpdateBy = jwtauth.UserIdStr(c)
	result, err := data.Update(data.Id)
	tools.HasError(err, "", -1)

	servers.OKWithRequestID(c, result, "")
}

func DeleteSysFileDir(c *gin.Context) {
	var data models.SysFileDir
	data.UpdateBy = jwtauth.UserIdStr(c)

	IDS := tools.IdsStrToIdsIntGroup(c.Param("id"))
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, codes.DeletedFail, 500)
	servers.OKWithRequestID(c, nil, codes.DeletedSuccess)
}
