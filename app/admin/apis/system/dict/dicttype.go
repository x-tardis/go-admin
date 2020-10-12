package dict

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

// @Summary 字典类型列表数据
// @Description 获取JSON
// @Tags 字典类型
// @Param dictName query string false "dictName"
// @Param dictId query string false "dictId"
// @Param dictType query string false "dictType"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} paginator.Page "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/type/list [get]
// @Security Bearer
func GetDictTypeList(c *gin.Context) {
	var data models.DictType
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, err = strconv.Atoi(size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, err = strconv.Atoi(index)
	}

	data.DictName = c.Request.FormValue("dictName")
	id := c.Request.FormValue("dictId")
	data.DictId = cast.ToInt(id)
	data.DictType = c.Request.FormValue("dictType")
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

// @Summary 通过字典id获取字典类型
// @Description 获取JSON
// @Tags 字典类型
// @Param dictId path int true "字典类型编码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/type/{dictId} [get]
// @Security Bearer
func GetDictType(c *gin.Context) {
	var DictType models.DictType
	DictType.DictName = c.Request.FormValue("dictName")
	DictType.DictId = cast.ToInt(c.Param("dictId"))
	result, err := DictType.Get()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.Success(c, servers.WithData(result))
}

func GetDictTypeOptionSelect(c *gin.Context) {
	var DictType models.DictType
	DictType.DictName = c.Request.FormValue("dictName")
	DictType.DictId = cast.ToInt(c.Param("dictId"))
	result, err := DictType.GetList()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.Success(c, servers.WithData(result))
}

// @Summary 添加字典类型
// @Description 获取JSON
// @Tags 字典类型
// @Accept  application/json
// @Product application/json
// @Param data body models.DictType true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dict/type [post]
// @Security Bearer
func InsertDictType(c *gin.Context) {
	var data models.DictType
	err := c.BindJSON(&data)
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
	servers.Success(c, servers.WithData(result))
}

// @Summary 修改字典类型
// @Description 获取JSON
// @Tags 字典类型
// @Accept  application/json
// @Product application/json
// @Param data body models.DictType true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dict/type [put]
// @Security Bearer
func UpdateDictType(c *gin.Context) {
	var data models.DictType
	err := c.BindWith(&data, binding.JSON)
	data.UpdateBy = jwtauth.UserIdStr(c)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	result, err := data.Update(data.DictId)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.Success(c, servers.WithData(result))
}

// @Summary 删除字典类型
// @Description 删除数据
// @Tags 字典类型
// @Param dictId path int true "dictId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/dict/type/{dictId} [delete]
func DeleteDictType(c *gin.Context) {
	var data models.DictType
	data.UpdateBy = jwtauth.UserIdStr(c)
	IDS := tools.IdsStrToIdsIntGroup(c.Param("dictId"))
	result, err := data.BatchDelete(IDS)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.OKWithRequestID(c, result, codes.DeletedSuccess)
}
