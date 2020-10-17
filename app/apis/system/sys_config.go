package system

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

type Config struct{}

// @Summary 配置列表数据
// @Description 获取JSON
// @Tags 配置
// @Param configKey query string false "configKey"
// @Param configName query string false "configName"
// @Param configType query string false "configType"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/configs [get]
// @Security Bearer
func (Config) QueryPage(c *gin.Context) {
	qp := models.SysConfigQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	qp.Inspect()

	result, ifc, err := new(models.CallSysConfig).QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}

	servers.JSON(c, http.StatusOK, servers.WithData(&paginator.Pages{
		Info: ifc,
		List: result,
	}))
}

// @Summary 获取配置
// @Description 获取JSON
// @Tags 配置
// @Param configId path int true "配置编码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/configs/{id} [get]
// @Security Bearer
func (Config) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := new(models.CallSysConfig).Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(item))
}

// @Summary 获取配置
// @Description 获取JSON
// @Tags 配置
// @Param configKey path int true "configKey"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/configKey/{key} [get]
// @Security Bearer
func (Config) GetWithKey(c *gin.Context) {
	key := c.Param("key")
	item, err := new(models.CallSysConfig).GetWithKey(gcontext.Context(c), key)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, item, item.ConfigValue)
}

// @Summary 添加配置
// @Description 获取JSON
// @Tags 配置
// @Accept  application/json
// @Product application/json
// @Param data body models.SysConfig true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/configs [post]
// @Security Bearer
func (Config) Create(c *gin.Context) {
	var data models.SysConfig

	if err := c.ShouldBindJSON(&data); err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}

	result, err := new(models.CallSysConfig).Create(gcontext.Context(c), data)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, result, "")
}

// @Summary 修改配置
// @Description 获取JSON
// @Tags 配置
// @Accept  application/json
// @Product application/json
// @Param data body models.SysConfig true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/configs [put]
// @Security Bearer
func (Config) Update(c *gin.Context) {
	var data models.SysConfig

	if err := c.ShouldBindJSON(&data); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}

	result, err := new(models.CallSysConfig).Update(gcontext.Context(c), data.ConfigId, data)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, result, "")
}

// @Summary 删除配置
// @Description 删除数据
// @Tags 配置
// @Param configId path int true "configId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/configs/{ids} [delete]
func (Config) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := new(models.CallSysConfig).BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithMsg(codes.DeletedSuccess))
}
