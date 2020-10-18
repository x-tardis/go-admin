package dict

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
)

type DictType struct{}

// @Summary 字典类型列表数据
// @Description 获取JSON
// @Tags 字典类型
// @Param dictName query string false "dictName"
// @Param dictId query string false "dictId"
// @Param dictType query string false "dictType"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} paginator.Page "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/type [get]
// @Security Bearer
func (DictType) QueryPage(c *gin.Context) {
	qp := models.DictTypeQueryParam{}
	if err := c.ShouldBindQuery(qp); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	qp.Inspect()

	items, ifc, err := new(models.CallDictType).QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(paginator.Pages{
		Info: ifc,
		List: items,
	}))
}

// @Summary 字典类型列表数据
// @Description 获取JSON
// @Tags 字典类型
// @Param dictName query string false "dictName"
// @Param dictId query string false "dictId"
// @Param dictType query string false "dictType"
// @Success 200 {object} paginator.Page "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/typeoptionselect [get]
// @Security Bearer
func (DictType) GetOptionSelect(c *gin.Context) {
	qp := models.DictTypeQueryParam{}
	if err := c.ShouldBindQuery(qp); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}

	items, err := new(models.CallDictType).Query(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(items))
}

// @Summary 通过字典id获取字典类型
// @Description 获取JSON
// @Tags 字典类型
// @Param dictId path int true "字典类型编码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/type/{dictId} [get]
// @Security Bearer
func (DictType) Get(c *gin.Context) {
	dictName := c.Query("dictName")
	dictId := cast.ToInt(c.Param("dictId"))
	item, err := new(models.CallDictType).Get(gcontext.Context(c), dictId, dictName)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(item))
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
func (DictType) Create(c *gin.Context) {
	item := models.DictType{}
	if err := c.ShouldBindJSON(&item); err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}

	result, err := new(models.CallDictType).Create(gcontext.Context(c), item)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(result))
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
func (DictType) Update(c *gin.Context) {
	up := models.DictType{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	item, err := new(models.CallDictType).Update(gcontext.Context(c), up.DictId, up)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(item))
}

// @Summary 删除字典类型
// @Description 删除数据
// @Tags 字典类型
// @Param dictId path int true "dictId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/dict/type/{dictId} [delete]
func (DictType) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("dictId"))
	err := new(models.CallDictType).BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithMsg(codes.DeletedSuccess))
}