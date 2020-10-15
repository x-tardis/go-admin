package tools

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/x-tardis/go-admin/app/models/tools"
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
	if data.TableName == "" {
		servers.Fail(c, 500, "table name cannot be empty")
		return
	}
	result, count, err := data.GetPage(pageSize, pageIndex)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(&paginator.Page{
		List:      result,
		Total:     count,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}))
}
