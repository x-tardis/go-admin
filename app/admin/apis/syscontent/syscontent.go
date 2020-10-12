package syscontent

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/cast"

	"github.com/x-tardis/go-admin/app/admin/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/paginator"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/tools"
)

func GetSysContentList(c *gin.Context) {
	var data models.SysContent
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, err = strconv.Atoi(size)
	}
	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, err = strconv.Atoi(index)
	}

	data.CateId = c.Request.FormValue("cateId")
	data.Name = c.Request.FormValue("name")
	data.Status = c.Request.FormValue("status")

	data.DataScope = jwtauth.UserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}

	servers.Success(c, servers.WithData(&paginator.Page{
		List:      result,
		Count:     count,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}))
}

func GetSysContent(c *gin.Context) {
	var data models.SysContent
	data.Id = cast.ToInt(c.Param("id"))
	result, err := data.Get()
	if err != nil {
		servers.Fail(c, -1, "抱歉未找到相关信息")
		return
	}
	servers.OKWithRequestID(c, result, "")
}

// @Summary 添加内容管理
// @Description 获取JSON
// @Tags 内容管理
// @Accept  application/json
// @Product application/json
// @Param data body models.SysContent true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/syscontent [post]
func InsertSysContent(c *gin.Context) {
	var data models.SysContent
	err := c.ShouldBindJSON(&data)
	data.CreateBy = jwtauth.UserIdStr(c)

	if err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}
	result, err := data.Create()
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, result, "")
}

func UpdateSysContent(c *gin.Context) {
	var data models.SysContent
	err := c.BindWith(&data, binding.JSON)
	if err != nil {
		servers.Fail(c, -1, "数据解析失败")
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

func DeleteSysContent(c *gin.Context) {
	var data models.SysContent
	data.UpdateBy = jwtauth.UserIdStr(c)

	IDS := tools.IdsStrToIdsIntGroup(c.Param("id"))
	_, err := data.BatchDelete(IDS)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.OKWithRequestID(c, nil, codes.DeletedSuccess)
}
