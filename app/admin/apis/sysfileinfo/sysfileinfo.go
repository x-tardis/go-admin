package sysfileinfo

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/x-tardis/go-admin/app/admin/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/paginator"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/tools"
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

	data.DataScope = jwtauth.UserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	tools.HasError(err, "", -1)

	servers.Success(c, servers.WithData(&paginator.Page{
		List:      result,
		Count:     count,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}))
}

func GetSysFileInfo(c *gin.Context) {
	var data models.SysFileInfo
	data.Id, _ = strconv.Atoi(c.Param("id"))
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	servers.OKWithRequestID(c, result, "")
}

func InsertSysFileInfo(c *gin.Context) {
	var data models.SysFileInfo
	err := c.ShouldBindJSON(&data)
	data.CreateBy = jwtauth.UserIdStr(c)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	servers.OKWithRequestID(c, result, "")
}

func UpdateSysFileInfo(c *gin.Context) {
	var data models.SysFileInfo
	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	data.UpdateBy = jwtauth.UserIdStr(c)
	result, err := data.Update(data.Id)
	tools.HasError(err, "", -1)

	servers.OKWithRequestID(c, result, "")
}

func DeleteSysFileInfo(c *gin.Context) {
	var data models.SysFileInfo
	data.UpdateBy = jwtauth.UserIdStr(c)

	IDS := tools.IdsStrToIdsIntGroup(c.Param("id"))
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, codes.DeletedFail, 500)
	servers.OKWithRequestID(c, nil, codes.DeletedSuccess)
}
