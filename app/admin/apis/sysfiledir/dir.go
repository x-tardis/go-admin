package sysfiledir

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/x-tardis/go-admin/app/admin/models"
	"github.com/x-tardis/go-admin/tools"
	"github.com/x-tardis/go-admin/tools/app"
	"github.com/x-tardis/go-admin/tools/app/msg"
)

func GetSysFileDirList(c *gin.Context) {
	var SysFileDir models.SysFileDir
	SysFileDir.Label = c.Request.FormValue("label")
	SysFileDir.PId, _ = strconv.Atoi(c.Request.FormValue("pid"))
	SysFileDir.Id, _ = strconv.Atoi(c.Request.FormValue("id"))
	SysFileDir.DataScope = tools.GetUserIdStr(c)
	result, err := SysFileDir.SetSysFileDir()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

func GetSysFileDir(c *gin.Context) {
	var data models.SysFileDir
	data.Id, _ = strconv.Atoi(c.Param("id"))
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.OK(c, result, "")
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
	data.CreateBy = tools.GetUserIdStr(c)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

func UpdateSysFileDir(c *gin.Context) {
	var data models.SysFileDir
	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	data.UpdateBy = tools.GetUserIdStr(c)
	result, err := data.Update(data.Id)
	tools.HasError(err, "", -1)

	app.OK(c, result, "")
}

func DeleteSysFileDir(c *gin.Context) {
	var data models.SysFileDir
	data.UpdateBy = tools.GetUserIdStr(c)

	IDS := tools.IdsStrToIdsIntGroup(c.Param("id"))
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}
