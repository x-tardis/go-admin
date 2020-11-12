package tools

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/models/tools"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/codes"
)

type Tables struct{}

// @Summary 分页列表数据
// @Description 生成表分页列表
// @Tags 工具 - 生成表
// @Param tableName query string false "tableName / 数据表名称"
// @Param pageSize query int false "pageSize / 页条数"
// @Param pageIndex query int false "pageIndex / 页码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sys/tables/page [get]
func (Tables) QueryPage(c *gin.Context) {
	qp := tools.TablesQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	items, info, err := tools.CTables.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithMsg(codes.QueryFailed),
			servers.WithError(err))
		return
	}

	servers.OK(c, servers.WithData(paginator.Pages{
		Info: info,
		List: items,
	}))
}

func (Tables) QueryTree(c *gin.Context) {
	items, err := tools.CTables.QueryTree(gcontext.Context(c))
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithMsg(codes.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(items))
}

// @Summary 获取配置
// @Description 获取JSON
// @Tags 工具 - 生成表
// @Param configKey path int true "configKey"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sys/tables/info/{id} [get]
// @Security Bearer
func (Tables) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	result, err := tools.CTables.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithMsg(codes.NotFound),
			servers.WithError(err))
		return
	}

	mp := make(map[string]interface{})
	mp["list"] = result.Columns
	mp["info"] = result
	servers.OK(c, servers.WithData(mp))
}

func (Tables) GetWithName(c *gin.Context) {
	tbName := c.Request.FormValue("tableName")
	item, err := tools.CTables.GetWithName(gcontext.Context(c), tbName)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithMsg(codes.NotFound),
			servers.WithError(err))
		return
	}
	mp := map[string]interface{}{
		"list": item.Columns,
		"info": item,
	}
	servers.OK(c, servers.WithData(mp))
}

// @Summary 添加表结构
// @Description 添加表结构
// @Tags 工具 - 生成表
// @Accept  application/json
// @Product application/json
// @Param tables query string false "tableName / 数据表名称"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/sys/tables/info [post]
// @Security Bearer
func (Tables) Create(c *gin.Context) {
	tablesList := strings.Split(c.Request.FormValue("tables"), ",")
	for i := 0; i < len(tablesList); i++ {
		data, err := genTableInit(tablesList, i, c)
		if err != nil {
			servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
			return
		}
		_, err = tools.CTables.Create(gcontext.Context(c), data)
		if err != nil {
			servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
			return
		}
	}
	servers.OK(c, servers.WithMsg("添加成功"))
}

func genTableInit(tablesList []string, i int, c *gin.Context) (tools.SysTables, error) {
	var data tools.SysTables

	data.TBName = tablesList[i]
	data.Creator = jwtauth.FromUserIdStr(gcontext.Context(c))

	dbtable, err := tools.CDBTables.Get(data.TBName)
	if err != nil {
		return tools.SysTables{}, err
	}

	tablenamelist := strings.Split(data.TBName, "_")
	for i := 0; i < len(tablenamelist); i++ {
		strStart := string([]byte(tablenamelist[i])[:1])
		strend := string([]byte(tablenamelist[i])[1:])
		data.ClassName += strings.ToUpper(strStart) + strend
		//data.PackageName += strings.ToLower(strStart) + strings.ToLower(strend)
		data.ModuleName += strings.ToLower(strStart) + strings.ToLower(strend)
	}
	data.PackageName = "admin"
	data.TplCategory = "crud"
	data.Crud = true

	dbcolumn, err := tools.CDBColumns.Query(data.TBName)
	if err != nil {
		return tools.SysTables{}, err
	}
	data.Creator = jwtauth.FromUserIdStr(gcontext.Context(c))
	data.TableComment = dbtable.TableComment
	if dbtable.TableComment == "" {
		data.TableComment = data.ClassName
	}

	data.FunctionName = data.TableComment
	data.BusinessName = data.ModuleName
	data.IsLogicalDelete = "1"
	data.LogicalDelete = true
	data.LogicalDeleteColumn = "is_del"
	data.IsActions = 1
	data.IsDataScope = 1
	data.IsAuth = 1

	data.FunctionAuthor = "wenjianzhang"
	for i := 0; i < len(dbcolumn); i++ {
		var column tools.SysColumns
		column.ColumnComment = dbcolumn[i].ColumnComment
		column.ColumnName = dbcolumn[i].ColumnName
		column.ColumnType = dbcolumn[i].ColumnType
		column.Sort = i + 1
		column.Insert = true
		column.IsInsert = "1"
		column.QueryType = "EQ"
		column.IsPk = "0"

		namelist := strings.Split(dbcolumn[i].ColumnName, "_")
		for i := 0; i < len(namelist); i++ {
			strStart := string([]byte(namelist[i])[:1])
			strend := string([]byte(namelist[i])[1:])
			column.GoField += strings.ToUpper(strStart) + strend
			if i == 0 {
				column.JsonField = strings.ToLower(strStart) + strend
			} else {
				column.JsonField += strings.ToUpper(strStart) + strend
			}
		}
		if strings.Contains(dbcolumn[i].ColumnKey, "PR") {
			column.IsPk = "1"
			column.Pk = true
			data.PkColumn = dbcolumn[i].ColumnName
			column.GoField = strings.ToUpper(column.GoField)
			column.JsonField = strings.ToUpper(column.JsonField)
			data.PkGoField = column.GoField
			data.PkJsonField = column.JsonField
		}
		column.IsRequired = "0"
		if strings.Contains(dbcolumn[i].IsNullable, "NO") {
			column.IsRequired = "1"
			column.Required = true
		}

		if strings.Contains(dbcolumn[i].ColumnType, "int") {
			if strings.Contains(dbcolumn[i].ColumnKey, "PR") {
				column.GoType = "uint"
			} else if strings.Contains(dbcolumn[i].ColumnType, "unsigned") {
				column.GoType = "uint"
			} else {
				column.GoType = "string"
			}
			column.HtmlType = "input"
		} else if strings.Contains(dbcolumn[i].ColumnType, "timestamp") {
			column.GoType = "time.Time"
			column.HtmlType = "datetime"
		} else if strings.Contains(dbcolumn[i].ColumnType, "datetime") {
			column.GoType = "time.Time"
			column.HtmlType = "datetime"
		} else {
			column.GoType = "string"
			column.HtmlType = "input"
		}

		data.Columns = append(data.Columns, column)
	}
	return data, err
}

// @Summary 修改表结构
// @Description 修改表结构
// @Tags 工具 - 生成表
// @Accept  application/json
// @Product application/json
// @Param data body tools.SysTables true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/sys/tables/info [put]
// @Security Bearer
func (Tables) Update(c *gin.Context) {
	up := tools.SysTables{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := tools.CTables.Update(gcontext.Context(c), up.TableId, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(codes.UpdateFailed))
		return
	}
	servers.OK(c, servers.WithData(item), servers.WithMsg(codes.UpdatedSuccess))
}

// @Summary 删除表结构
// @Description 删除表结构
// @Tags 工具 - 生成表
// @Param tableId path int true "tableId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/sys/tables/info/{ids} [delete]
func (Tables) DeleteSysTables(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := tools.CTables.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(codes.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithMsg(codes.DeleteSuccess))
}
