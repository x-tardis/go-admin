package log

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/x-tardis/go-admin/app/admin/models"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/paginator"
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
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/operloglist [get]
// @Security Bearer
func GetOperLogList(c *gin.Context) {
	var data models.SysOperLog
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

	data.OperName = c.Request.FormValue("operName")
	data.Status = c.Request.FormValue("status")
	data.OperIp = c.Request.FormValue("operIp")
	result, count, err := data.GetPage(pageSize, pageIndex)
	tools.HasError(err, "", -1)

	servers.Success(c, servers.WithData(&paginator.Page{
		List:      result,
		Count:     count,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}))
}

// @Summary 通过编码获取登录日志
// @Description 获取JSON
// @Tags 登录日志
// @Param infoId path int true "infoId"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/operlog/{infoId} [get]
// @Security Bearer
func GetOperLog(c *gin.Context) {
	var OperLog models.SysOperLog
	OperLog.OperId, _ = strconv.Atoi(c.Param("operId"))
	result, err := OperLog.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	servers.Success(c, servers.WithData(result))
}

// @Summary 添加操作日志
// @Description 获取JSON
// @Tags 操作日志
// @Accept  application/json
// @Product application/json
// @Param data body models.SysOperLog true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/operlog [post]
// @Security Bearer
func InsertOperLog(c *gin.Context) {
	var data models.SysOperLog
	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	servers.Success(c, servers.WithData(result))
}

// @Summary 批量删除操作日志
// @Description 删除数据
// @Tags 操作日志
// @Param operId path string true "以逗号（,）分割的operId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/operlog/{operId} [delete]
func DeleteOperLog(c *gin.Context) {
	var data models.SysOperLog
	data.UpdateBy = jwtauth.UserIdStr(c)
	IDS := tools.IdsStrToIdsIntGroup(c.Param("operId"))
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, "删除失败", 500)
	servers.Success(c, servers.WithMessage("删除成功"))
}
