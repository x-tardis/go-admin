package tools

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/x-tardis/go-admin/app/admin/models/tools"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
	tools2 "github.com/x-tardis/go-admin/tools"
)

// @Summary 分页列表数据
// @Description 生成表分页列表
// @Tags 工具 - 生成表
// @Param tableName query string false "tableName / 数据表名称"
// @Param pageSize query int false "pageSize / 页条数"
// @Param pageIndex query int false "pageIndex / 页码"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sys/tables/page [get]
func GetSysTableList(c *gin.Context) {
	var data tools.SysTables
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, err = strconv.Atoi(size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, err = strconv.Atoi(index)
	}

	data.TBName = c.Request.FormValue("tableName")
	data.TableComment = c.Request.FormValue("tableComment")
	result, count, err := data.GetPage(pageSize, pageIndex)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	var mp = make(map[string]interface{}, 3)
	mp["list"] = result
	mp["count"] = count
	mp["pageIndex"] = pageIndex
	mp["pageSize"] = pageSize

	servers.Success(c, servers.WithData(mp))
}

// @Summary 获取配置
// @Description 获取JSON
// @Tags 工具 - 生成表
// @Param configKey path int true "configKey"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sys/tables/info/{tableId} [get]
// @Security Bearer
func GetSysTables(c *gin.Context) {
	var data tools.SysTables
	data.TableId = cast.ToInt(c.Param("tableId"))
	result, err := data.Get()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}

	mp := make(map[string]interface{})
	mp["list"] = result.Columns
	mp["info"] = result
	servers.Success(c, servers.WithData(mp))
}

func GetSysTablesInfo(c *gin.Context) {
	var data tools.SysTables
	if c.Request.FormValue("tableName") != "" {
		data.TBName = c.Request.FormValue("tableName")
	}
	result, err := data.Get()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	mp := make(map[string]interface{})
	mp["list"] = result.Columns
	mp["info"] = result
	servers.Success(c, servers.WithData(mp))
}

func GetSysTablesTree(c *gin.Context) {
	var data tools.SysTables
	result, err := data.GetTree()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.Success(c, servers.WithData(result))
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
func InsertSysTable(c *gin.Context) {
	tablesList := strings.Split(c.Request.FormValue("tables"), ",")
	for i := 0; i < len(tablesList); i++ {
		data, err := genTableInit(tablesList, i, c)
		if err != nil {
			servers.Fail(c, -1, err.Error())
			return
		}
		_, err = data.Create()
		if err != nil {
			servers.Fail(c, -1, err.Error())
			return
		}
	}
	servers.Success(c, servers.WithMessage("添加成功"))
}

func genTableInit(tablesList []string, i int, c *gin.Context) (tools.SysTables, error) {
	var data tools.SysTables
	var dbTable tools.DBTables
	var dbColumn tools.DBColumns
	data.TBName = tablesList[i]
	data.CreateBy = jwtauth.UserIdStr(c)

	dbTable.TableName = data.TBName
	dbtable, err := dbTable.Get()
	if err != nil {
		return tools.SysTables{}, err
	}
	dbColumn.TableName = data.TBName
	tablenamelist := strings.Split(dbColumn.TableName, "_")
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

	dbcolumn, err := dbColumn.GetList()
	if err != nil {
		return tools.SysTables{}, err
	}
	data.CreateBy = jwtauth.UserIdStr(c)
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
func UpdateSysTable(c *gin.Context) {
	var data tools.SysTables
	err := c.Bind(&data)
	if err != nil {
		servers.Fail(c, 500, codes.DataParseFailed)
		return
	}

	data.UpdateBy = jwtauth.UserIdStr(c)
	result, err := data.Update()
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.Success(c, servers.WithData(result), servers.WithMessage(codes.UpdatedSuccess))
}

// @Summary 删除表结构
// @Description 删除表结构
// @Tags 工具 - 生成表
// @Param tableId path int true "tableId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/sys/tables/info/{tableId} [delete]
func DeleteSysTables(c *gin.Context) {
	var data tools.SysTables
	IDS := tools2.IdsStrToIdsIntGroup(c.Param("tableId"))
	_, err := data.BatchDelete(IDS)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.Success(c, servers.WithMessage(codes.DeletedSuccess))
}
