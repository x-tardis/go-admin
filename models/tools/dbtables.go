package tools

import (
	"errors"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/deployed/dao"
)

type DBTables struct {
	TableName      string `gorm:"column:TABLE_NAME" json:"tableName"`
	Engine         string `gorm:"column:ENGINE" json:"engine"`
	TableRows      string `gorm:"column:TABLE_ROWS" json:"tableRows"`
	TableCollation string `gorm:"column:TABLE_COLLATION" json:"tableCollation"`
	CreateTime     string `gorm:"column:CREATE_TIME" json:"createTime"`
	UpdateTime     string `gorm:"column:UPDATE_TIME" json:"updateTime"`
	TableComment   string `gorm:"column:TABLE_COMMENT" json:"tableComment"`
}

type DbTablesQueryParam struct {
	TableName string `form:"tableName"`
	paginator.Param
}

type cDBTables struct{}

var CDBTables = cDBTables{}

func (cDBTables) QueryPage(qp DbTablesQueryParam) ([]DBTables, paginator.Info, error) {
	var items []DBTables

	if deployed.DbConfig.Driver != "mysql" {
		return nil, paginator.Info{}, errors.New("目前只支持mysql数据库")
	}

	db := dao.DB.Table("information_schema.tables").
		Where("TABLE_NAME not in (select table_name from "+deployed.GenConfig.DBName+".sys_tables) ").
		Where("table_schema= ? ", deployed.GenConfig.DBName)
	if qp.TableName != "" {
		db = db.Where("TABLE_NAME = ?", qp.TableName)
	}
	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

func (cDBTables) Get(tableName string) (item DBTables, err error) {
	if deployed.DbConfig.Driver != "mysql" {
		return item, errors.New("目前只支持mysql数据库")
	}
	if tableName == "" {
		return item, errors.New("table name cannot be empty！")
	}
	err = dao.DB.Table("information_schema.tables").
		Where("table_schema= ? ", deployed.GenConfig.DBName).
		Where("TABLE_NAME = ?", tableName).
		First(&item).Error
	return
}
