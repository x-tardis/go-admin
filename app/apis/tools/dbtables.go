package tools

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thinkgos/sharp/core/paginator"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/models/tools"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// @Summary 分页列表数据 / page list data
// @Description 数据库表分页列表 / database table page list
// @Tags 工具 / Tools
// @Param tableName query string false "tableName / 数据表名称"
// @Param pageSize query int false "pageSize / 页条数"
// @Param pageIndex query int false "pageIndex / 页码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/db/tables/page [get]
func GetDBTableList(c *gin.Context) {
	if deployed.DbConfig.Driver == "sqlite3" || deployed.DbConfig.Driver == "postgres" {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithMsg("对不起，sqlite3 或 postgres 不支持代码生成"))
		return
	}
	qp := tools.DbTablesQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest,
			servers.WithError(err))
		return
	}
	qp.Inspect()

	items, info, err := tools.CDBTables.QueryPage(qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(&paginator.Pages{
		Info: info,
		List: items,
	}))
}
