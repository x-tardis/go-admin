package syscategory

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

func GetSysCategoryList(c *gin.Context) {
	var data models.SysCategory
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, err = strconv.Atoi(size)
	}
	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, err = strconv.Atoi(index)
	}

	data.Name = c.Request.FormValue("name")
	data.Status = c.Request.FormValue("status")

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

func GetSysCategory(c *gin.Context) {
	var data models.SysCategory
	data.Id = cast.ToInt(c.Param("id"))
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)

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
	err := c.ShouldBindJSON(&data)
	data.CreateBy = jwtauth.UserIdStr(c)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	servers.OKWithRequestID(c, result, "")
}

func UpdateSysCategory(c *gin.Context) {
	var data models.SysCategory
	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	data.UpdateBy = jwtauth.UserIdStr(c)
	result, err := data.Update(data.Id)
	tools.HasError(err, "", -1)

	servers.OKWithRequestID(c, result, "")
}

func DeleteSysCategory(c *gin.Context) {
	var data models.SysCategory
	data.UpdateBy = jwtauth.UserIdStr(c)

	IDS := tools.IdsStrToIdsIntGroup(c.Param("id"))
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, codes.DeletedFail, 500)
	servers.OKWithRequestID(c, nil, codes.DeletedSuccess)
}
