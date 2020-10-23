package tools

import (
	"errors"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/deployed/dao"
)

type DBColumns struct {
	TableSchema            string `gorm:"column:TABLE_SCHEMA" json:"tableSchema"`
	TableName              string `gorm:"column:TABLE_NAME" json:"tableName"`
	ColumnName             string `gorm:"column:COLUMN_NAME" json:"columnName"`
	ColumnDefault          string `gorm:"column:COLUMN_DEFAULT" json:"columnDefault"`
	IsNullable             string `gorm:"column:IS_NULLABLE" json:"isNullable"`
	DataType               string `gorm:"column:DATA_TYPE" json:"dataType"`
	CharacterMaximumLength string `gorm:"column:CHARACTER_MAXIMUM_LENGTH" json:"characterMaximumLength"`
	CharacterSetName       string `gorm:"column:CHARACTER_SET_NAME" json:"characterSetName"`
	ColumnType             string `gorm:"column:COLUMN_TYPE" json:"columnType"`
	ColumnKey              string `gorm:"column:COLUMN_KEY" json:"columnKey"`
	Extra                  string `gorm:"column:EXTRA" json:"extra"`
	ColumnComment          string `gorm:"column:COLUMN_COMMENT" json:"columnComment"`
}

func (e *DBColumns) GetPage(param paginator.Param) ([]DBColumns, paginator.Info, error) {
	var doc []DBColumns

	table := new(gorm.DB)

	if deployed.DbConfig.Driver == "mysql" {
		table = dao.DB.Table("information_schema.`COLUMNS`")
		table = table.Where("table_schema= ? ", deployed.GenConfig.DBName)

		if e.TableName != "" {
			return nil, paginator.Info{}, errors.New("table name cannot be empty！")
		}

		table = table.Where("TABLE_NAME = ?", e.TableName)
	}

	info, err := iorm.QueryPages(table, param, &doc)
	if err != nil {
		return nil, info, err
	}
	return doc, info, err
}

func (e *DBColumns) GetList() ([]DBColumns, error) {
	var doc []DBColumns
	table := new(gorm.DB)

	if e.TableName == "" {
		return nil, errors.New("table name cannot be empty！")
	}

	if deployed.DbConfig.Driver == "mysql" {
		table = dao.DB.Table("information_schema.columns")
		table = table.Where("table_schema= ? ", deployed.GenConfig.DBName)

		table = table.Where("TABLE_NAME = ?", e.TableName).Order("ORDINAL_POSITION asc")
	} else {
		return nil, errors.New("目前只支持mysql数据库")
	}
	if err := table.Find(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}