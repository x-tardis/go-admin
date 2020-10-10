package sysfileinfo

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/x-tardis/go-admin/app/admin/models"
	"github.com/x-tardis/go-admin/tools"
	"github.com/x-tardis/go-admin/tools/app"
	"github.com/x-tardis/go-admin/tools/app/msg"
)

func GetSysFileInfoList(c *gin.Context) {
	var data models.SysFileInfo
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, err = strconv.Atoi(size)
	}
	tools.HasError(err, "", -1)
	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, err = strconv.Atoi(index)
	}
	tools.HasError(err, "", -1)
	if pid := c.Request.FormValue("pId"); pid != "" {
		data.PId, err = strconv.Atoi(pid)
	}
	tools.HasError(err, "", -1)

	data.DataScope = tools.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	tools.HasError(err, "", -1)

	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

func GetSysFileInfo(c *gin.Context) {
	var data models.SysFileInfo
	data.Id, _ = strconv.Atoi(c.Param("id"))
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.OK(c, result, "")
}

func InsertSysFileInfo(c *gin.Context) {
	var data models.SysFileInfo
	err := c.ShouldBindJSON(&data)
	data.CreateBy = tools.GetUserIdStr(c)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

func UpdateSysFileInfo(c *gin.Context) {
	var data models.SysFileInfo
	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	data.UpdateBy = tools.GetUserIdStr(c)
	result, err := data.Update(data.Id)
	tools.HasError(err, "", -1)

	app.OK(c, result, "")
}

func DeleteSysFileInfo(c *gin.Context) {
	var data models.SysFileInfo
	data.UpdateBy = tools.GetUserIdStr(c)

	IDS := tools.IdsStrToIdsIntGroup("id", c)
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}
