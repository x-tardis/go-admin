package log

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/x-tardis/go-admin/app/admin/models"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/tools"
)

// @Summary 登录日志列表
// @Description 获取JSON
// @Tags 登录日志
// @Param status query string false "status"
// @Param dictCode query string false "dictCode"
// @Param dictType query string false "dictType"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/loginloglist [get]
// @Security Bearer
func GetLoginLogList(c *gin.Context) {
	var data models.LoginLog
	var err error
	var pageSize = 10
	var pageIndex = 1

	size := c.Request.FormValue("pageSize")
	if size != "" {
		pageSize, err = strconv.Atoi(size)
	}

	index := c.Request.FormValue("pageIndex")
	if index != "" {
		pageIndex, err = strconv.Atoi(index)
	}

	data.Username = c.Request.FormValue("username")
	data.Status = c.Request.FormValue("status")
	data.Ipaddr = c.Request.FormValue("ipaddr")
	result, count, err := data.GetPage(pageSize, pageIndex)
	tools.HasError(err, "", -1)

	var mp = make(map[string]interface{}, 3)
	mp["list"] = result
	mp["count"] = count
	mp["pageIndex"] = pageIndex
	mp["pageSize"] = pageSize

	servers.Success(c, servers.WithData(mp))
}

// @Summary 通过编码获取登录日志
// @Description 获取JSON
// @Tags 登录日志
// @Param infoId path int true "infoId"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/loginlog/{infoId} [get]
// @Security Bearer
func GetLoginLog(c *gin.Context) {
	var LoginLog models.LoginLog
	LoginLog.InfoId = cast.ToInt(c.Param("infoId"))
	result, err := LoginLog.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	servers.Success(c, servers.WithData(result))
}

// @Summary 添加登录日志
// @Description 获取JSON
// @Tags 登录日志
// @Accept  application/json
// @Product application/json
// @Param data body models.LoginLog true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/loginlog [post]
// @Security Bearer
func InsertLoginLog(c *gin.Context) {
	var data models.LoginLog
	err := c.BindJSON(&data)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	servers.Success(c, servers.WithData(result))
}

// @Summary 修改登录日志
// @Description 获取JSON
// @Tags 登录日志
// @Accept  application/json
// @Product application/json
// @Param data body models.LoginLog true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/loginlog [put]
// @Security Bearer
func UpdateLoginLog(c *gin.Context) {
	var data models.LoginLog
	err := c.BindJSON(&data)
	tools.HasError(err, "", -1)
	result, err := data.Update(data.InfoId)
	tools.HasError(err, "", -1)
	servers.Success(c, servers.WithData(result))
}

// @Summary 批量删除登录日志
// @Description 删除数据
// @Tags 登录日志
// @Param infoId path string true "以逗号（,）分割的infoId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/loginlog/{infoId} [delete]
func DeleteLoginLog(c *gin.Context) {
	var data models.LoginLog
	data.UpdateBy = tools.GetUserIdStr(c)
	IDS := tools.IdsStrToIdsIntGroup(c.Param("infoId"))
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, "修改失败", 500)
	servers.Success(c, servers.WithMessage("删除成功"))
}
