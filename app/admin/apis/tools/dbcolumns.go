package tools

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/admin/models/tools"
	"github.com/x-tardis/go-admin/pkg/paginator"
	"github.com/x-tardis/go-admin/pkg/servers"
	tools2 "github.com/x-tardis/go-admin/tools"
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
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, err = strconv.Atoi(size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, err = strconv.Atoi(index)
	}

	data.TableName = c.Request.FormValue("tableName")
	tools2.Assert(data.TableName == "", "table name cannot be empty！", 500)
	result, count, err := data.GetPage(pageSize, pageIndex)
	tools2.HasError(err, "", -1)

	servers.Success(c, servers.WithData(&paginator.Page{
		List:      result,
		Count:     count,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}))
}
