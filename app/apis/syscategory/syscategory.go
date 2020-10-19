package syscategory

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
)

func GetSysCategoryList(c *gin.Context) {
	var data models.SysCategory

	param := paginator.Param{
		PageIndex: cast.ToInt(c.Query("pageIndex")),
		PageSize:  cast.ToInt(c.Query("pageSize")),
	}
	param.Inspect()

	data.Name = c.Request.FormValue("name")
	data.Status = c.Request.FormValue("status")

	data.DataScope = jwtauth.UserIdStr(c)
	result, ifc, err := data.GetPage(param)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(paginator.Pages{
		Info: ifc,
		List: result,
	}))
}

func GetSysCategory(c *gin.Context) {
	var data models.SysCategory
	data.Id = cast.ToInt(c.Param("id"))
	result, err := data.Get()
	if err != nil {
		servers.Fail(c, -1, "抱歉未找到相关信息")
		return
	}
	servers.OKWithRequestID(c, result, "")
}

// @Summary 添加分类
// @Description 获取JSON
// @Tags 分类
// @Accept  application/json
// @Product application/json
// @Param data body models.SysCategory true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/syscategory [post]
func InsertSysCategory(c *gin.Context) {
	var data models.SysCategory

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

func UpdateSysCategory(c *gin.Context) {
	var data models.SysCategory

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

func DeleteSysCategory(c *gin.Context) {
	var data models.SysCategory
	data.Updator = jwtauth.UserIdStr(c)

	IDS := infra.ParseIdsGroup(c.Param("id"))
	_, err := data.BatchDelete(IDS)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.OKWithRequestID(c, nil, codes.DeletedSuccess)
}
