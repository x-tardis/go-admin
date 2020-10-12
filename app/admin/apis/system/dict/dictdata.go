package dict

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/cast"

	"github.com/x-tardis/go-admin/app/admin/models"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/paginator"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/tools"
)

// @Summary 字典数据列表
// @Description 获取JSON
// @Tags 字典数据
// @Param status query string false "status"
// @Param dictCode query string false "dictCode"
// @Param dictType query string false "dictType"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/data/list [get]
// @Security Bearer
func GetDictDataList(c *gin.Context) {
	var data models.DictData
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, err = strconv.Atoi(size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, err = strconv.Atoi(index)
	}

	data.DictLabel = c.Request.FormValue("dictLabel")
	data.Status = c.Request.FormValue("status")
	data.DictType = c.Request.FormValue("dictType")
	id := c.Request.FormValue("dictCode")
	data.DictCode = cast.ToInt(id)
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

// @Summary 通过编码获取字典数据
// @Description 获取JSON
// @Tags 字典数据
// @Param dictCode path int true "字典编码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/data/{dictCode} [get]
// @Security Bearer
func GetDictData(c *gin.Context) {
	var DictData models.DictData
	DictData.DictLabel = c.Request.FormValue("dictLabel")
	DictData.DictCode = cast.ToInt(c.Param("dictCode"))
	result, err := DictData.GetByCode()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	servers.Success(c, servers.WithData(result))
}

// @Summary 通过字典类型获取字典数据
// @Description 获取JSON
// @Tags 字典数据
// @Param dictType path int true "dictType"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/databyType/{dictType} [get]
// @Security Bearer
func GetDictDataByDictType(c *gin.Context) {
	var DictData models.DictData
	DictData.DictType = c.Param("dictType")
	result, err := DictData.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	servers.Success(c, servers.WithData(result))
}

// @Summary 添加字典数据
// @Description 获取JSON
// @Tags 字典数据
// @Accept  application/json
// @Product application/json
// @Param data body models.DictType true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dict/data [post]
// @Security Bearer
func InsertDictData(c *gin.Context) {
	var data models.DictData
	err := c.ShouldBindJSON(&data)
	data.CreateBy = jwtauth.UserIdStr(c)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	servers.Success(c, servers.WithData(result))
}

// @Summary 修改字典数据
// @Description 获取JSON
// @Tags 字典数据
// @Accept  application/json
// @Product application/json
// @Param data body models.DictType true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dict/data [put]
// @Security Bearer
func UpdateDictData(c *gin.Context) {
	var data models.DictData
	err := c.BindWith(&data, binding.JSON)
	data.UpdateBy = jwtauth.UserIdStr(c)
	tools.HasError(err, "", -1)
	result, err := data.Update(data.DictCode)
	tools.HasError(err, "", -1)
	servers.Success(c, servers.WithData(result))
}

// @Summary 删除字典数据
// @Description 删除数据
// @Tags 字典数据
// @Param dictCode path int true "dictCode"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/dict/data/{dictCode} [delete]
func DeleteDictData(c *gin.Context) {
	var data models.DictData
	data.UpdateBy = jwtauth.UserIdStr(c)
	IDS := tools.IdsStrToIdsIntGroup(c.Param("dictCode"))
	result, err := data.BatchDelete(IDS)
	tools.HasError(err, "修改失败", 500)
	servers.OKWithRequestID(c, result, "删除成功")
}
