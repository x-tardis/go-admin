package tools

import (
	"errors"

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

func (e *DBTables) GetPage(pageSize int, pageIndex int) ([]DBTables, int, error) {
	var doc []DBTables
	table := new(gorm.DB)
	var count int64

	if deployed.DatabaseConfig.Driver == "mysql" {
		table = deployed.DB.Table("information_schema.tables")
		table = table.Where("TABLE_NAME not in (select table_name from " + deployed.GenConfig.DBName + ".sys_tables) ")
		table = table.Where("table_schema= ? ", deployed.GenConfig.DBName)

		if e.TableName != "" {
			table = table.Where("TABLE_NAME = ?", e.TableName)
		}
		if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
			return nil, 0, err
		}
	} else {
		return nil, 0, errors.New("目前只支持mysql数据库")
	}

	//table.Count(&count)
	return doc, int(count), nil
}

func (e *DBTables) Get() (DBTables, error) {
	var doc DBTables
	table := new(gorm.DB)
	if deployed.DatabaseConfig.Driver == "mysql" {
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
