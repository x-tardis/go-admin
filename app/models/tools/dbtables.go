package tools

import (
	"errors"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/pkg/deployed"
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

func (e *DBTables) GetPage(params paginator.Param) ([]DBTables, paginator.Info, error) {
	var doc []DBTables

	if deployed.DbConfig.Driver != "mysql" {
		return nil, paginator.Info{}, errors.New("目前只支持mysql数据库")
	}

	table := deployed.DB.Table("information_schema.tables")
	table = table.Where("TABLE_NAME not in (select table_name from " + deployed.GenConfig.DBName + ".sys_tables) ")
	table = table.Where("table_schema= ? ", deployed.GenConfig.DBName)

	if e.TableName != "" {
		table = table.Where("TABLE_NAME = ?", e.TableName)
	}
	info, err := iorm.QueryPages(table, params, &doc)
	if err != nil {
		return nil, info, err
	}
	return doc, info, err
}

func (e *DBTables) Get() (DBTables, error) {
	var doc DBTables
	table := new(gorm.DB)
	if deployed.DbConfig.Driver == "mysql" {
		table = deployed.DB.Table("information_schema.tables")
		table = table.Where("table_schema= ? ", deployed.GenConfig.DBName)
		if e.TableName == "" {
			return doc, errors.New("table name cannot be empty！")
		}
		table = table.Where("TABLE_NAME = ?", e.TableName)
	} else {
		return DBTables{}, errors.New("目前只支持mysql数据库")
	}
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}
