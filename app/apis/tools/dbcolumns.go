package tools

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/thinkgos/sharp/core/paginator"

	"github.com/x-tardis/go-admin/models/tools"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// @Summary 分页列表数据 / page list data
// @Description 数据库表列分页列表 / database table column page list
// @Tags 工具 / Tools
// @Param tableName query string false "tableName / 数据表名称"
// @Param pageSize query int false "pageSize / 页条数"
// @Param pageIndex query int false "pageIndex / 页码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/db/columns/page [get]
func GetDBColumnList(c *gin.Context) {
	var data tools.DBColumns

	param := paginator.Param{
		PageIndex: cast.ToInt(c.Query("pageIndex")),
		PageSize:  cast.ToInt(c.Query("pageSize")),
	}
	param.Inspect()

	data.TableName = c.Request.FormValue("tableName")
	if data.TableName == "" {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithMsg("table name cannot be empty"))
		return
	}
	result, info, err := data.GetPage(param)
	if err != nil {
		servers.Fail(c, http.StatusOK, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(paginator.Pages{
		Info: info,
		List: result,
	}))
}
