package log

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

type OperLog struct{}

// @Tags 操作日志
// @Summary 登录日志列表
// @Description 获取JSON
// @Param status query string false "status"
// @Param dictCode query string false "dictCode"
// @Param dictType query string false "dictType"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/operlog [get]
// @Security Bearer
func (OperLog) QueryPage(c *gin.Context) {
	qp := models.SysOperLogQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	qp.Inspect()

	result, ifc, err := new(models.CallSysOperLog).QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(&paginator.Pages{
		Info: ifc,
		List: result,
	}))
}

// @Tags 操作日志
// @Summary 通过编码获取登录日志
// @Description 获取JSON
// @Param infoId path int true "infoId"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/operlog/{id} [get]
// @Security Bearer
func (OperLog) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := new(models.CallSysOperLog).Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(item))
}

// @Tags 操作日志
// @Summary 添加操作日志
// @Description 获取JSON
// @Accept  application/json
// @Product application/json
// @Param data body models.SysOperLog true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/operlog [post]
// @Security Bearer
func (OperLog) Create(c *gin.Context) {
	var item models.SysOperLog

	if err := c.ShouldBindJSON(&item); err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}
	result, err := new(models.CallSysOperLog).Create(gcontext.Context(c), item)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(result))
}

// @Tags 操作日志
// @Summary 批量删除操作日志
// @Description 删除数据
// @Param operId path string true "以逗号（,）分割的operId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/operlog/{ids} [delete]
func (OperLog) BatchDelete(c *gin.Context) {
	var err error

	action := c.Param("ids")
	switch action {
	case "clean":
		err = new(models.CallSysOperLog).Clean(gcontext.Context(c))
	default: // ids
		ids := infra.ParseIdsGroup(action)
		err = new(models.CallSysOperLog).BatchDelete(gcontext.Context(c), ids)
	}

	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithMsg(codes.DeletedSuccess))
}
