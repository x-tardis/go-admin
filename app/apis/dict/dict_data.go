package dict

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/thinkgos/sharp/core/paginator"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
)

type DictData struct{}

// @Summary 字典数据列表
// @Description 获取JSON
// @Tags 字典数据
// @Param dictCode query int false "dictCode"
// @Param dictType query string false "dictType"
// @Param dictLabel query string false "dictLabel"
// @Param status query string false "status"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/data [get]
// @Security Bearer
func (DictData) QueryPage(c *gin.Context) {
	qp := models.DictDataQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	qp.Inspect()

	result, ifc, err := new(models.CallDictData).QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}

	servers.JSON(c, http.StatusOK, servers.WithData(&paginator.Pages{
		Info: ifc,
		List: result,
	}))
}

// @Summary 通过编码获取字典数据
// @Description 获取JSON
// @Tags 字典数据
// @Param dictCode path int true "字典编码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/data/{dictCode} [get]
// @Security Bearer
func (DictData) Get(c *gin.Context) {
	dictLabel := c.Query("dictLabel")
	dictCode := cast.ToInt(c.Param("dictCode"))
	result, err := new(models.CallDictData).Get(gcontext.Context(c), dictCode, dictLabel)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(result))
}

// @Summary 通过字典类型获取字典数据
// @Description 获取JSON
// @Tags 字典数据
// @Param dictType path int true "dictType"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dict/databyType/{dictType} [get]
// @Security Bearer
func (DictData) GetWithType(c *gin.Context) {
	dictType := c.Param("dictType")
	result, err := new(models.CallDictData).GetWithType(gcontext.Context(c), dictType)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(result))
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
func (DictData) Create(c *gin.Context) {
	item := models.DictData{}
	if err := c.ShouldBindJSON(&item); err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}

	result, err := new(models.CallDictData).Create(gcontext.Context(c), item)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(result))
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
func (DictData) Update(c *gin.Context) {
	var up models.DictData

	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}

	result, err := new(models.CallDictData).Update(gcontext.Context(c), up.DictCode, up)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(result))
}

// @Summary 删除字典数据
// @Description 删除数据
// @Tags 字典数据
// @Param dictCode path int true "dictCode"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/dict/data/{dictCode} [delete]
func (DictData) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("dictCode"))
	err := new(models.CallDictData).BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithMsg("删除成功"))
}
