package sysfileinfo

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
)

func GetSysFileInfoList(c *gin.Context) {
	var data models.SysFileInfo
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, err = strconv.Atoi(size)
	}
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, err = strconv.Atoi(index)
	}
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	if pid := c.Request.FormValue("pId"); pid != "" {
		data.PId, err = strconv.Atoi(pid)
	}
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}

	data.DataScope = jwtauth.UserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.Success(c, servers.WithData(&paginator.Page{
		List:      result,
		Total:     count,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}))
}

func GetSysFileInfo(c *gin.Context) {
	var data models.SysFileInfo
	data.Id = cast.ToInt(c.Param("id"))
	result, err := data.Get()
	if err != nil {
		servers.Fail(c, -1, "抱歉未找到相关信息")
		return
	}
	servers.OKWithRequestID(c, result, "")
}

func InsertSysFileInfo(c *gin.Context) {
	var data models.SysFileInfo

	if err := c.ShouldBindJSON(&data); err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}
	data.CreateBy = jwtauth.UserIdStr(c)
	result, err := data.Create()
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, result, "")
}

func UpdateSysFileInfo(c *gin.Context) {
	var data models.SysFileInfo

	if err := c.ShouldBindJSON(&data); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}

	data.UpdateBy = jwtauth.UserIdStr(c)
	result, err := data.Update(data.Id)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, result, "")
}

func DeleteSysFileInfo(c *gin.Context) {
	var data models.SysFileInfo
	data.UpdateBy = jwtauth.UserIdStr(c)

	IDS := infra.ParseIdsGroup(c.Param("id"))
	_, err := data.BatchDelete(IDS)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.OKWithRequestID(c, nil, codes.DeletedSuccess)
}
