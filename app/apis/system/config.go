package system

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
func GetConfigList(c *gin.Context) {
	var data models.SysConfig

	param := paginator.Param{
		PageIndex: cast.ToInt(c.Query("pageIndex")),
		PageSize:  cast.ToInt(c.Query("pageSize")),
	}
	param.Inspect()

	data.ConfigKey = c.Request.FormValue("configKey")
	data.ConfigName = c.Request.FormValue("configName")
	data.ConfigType = c.Request.FormValue("configType")
	data.DataScope = jwtauth.UserIdStr(c)
	result, ifc, err := data.GetPage(param)
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
func GetConfig(c *gin.Context) {
	var Config models.SysConfig

	Config.ConfigId = cast.ToInt(c.Param("id"))
	result, err := Config.Get()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(result))
}

// @Summary 获取配置
// @Description 获取JSON
// @Tags 配置
// @Param configKey path int true "configKey"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/configKey/{configKey} [get]
// @Security Bearer
func GetConfigByConfigKey(c *gin.Context) {
	var Config models.SysConfig
	Config.ConfigKey = c.Param("configKey")
	result, err := Config.Get()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}

	servers.OKWithRequestID(c, result, result.ConfigValue)
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
func InsertConfig(c *gin.Context) {
	var data models.SysConfig

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
func UpdateConfig(c *gin.Context) {
	var data models.SysConfig

	if err := c.ShouldBindJSON(&data); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	data.UpdateBy = jwtauth.UserIdStr(c)
	result, err := data.Update(data.ConfigId)
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
func DeleteConfig(c *gin.Context) {
	var data models.SysConfig

	data.UpdateBy = jwtauth.UserIdStr(c)
	ids := infra.ParseIdsGroup(c.Param("ids"))
	result, err := data.BatchDelete(ids)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.OKWithRequestID(c, result, codes.DeletedSuccess)
}
