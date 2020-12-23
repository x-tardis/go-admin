package tools

import (
	"errors"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"

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

type DBColumnsQueryParam struct {
	TableName string `form:"tableName"`
	paginator.Param
}

type cDBColumns struct{}

var CDBColumns = cDBColumns{}

func (cDBColumns) Query(tbName string) (items []DBColumns, err error) {
	if tbName == "" {
		return nil, errors.New("table name cannot be empty！")
	}
	if dao.DbConfig.Dialect != "mysql" {
		return nil, errors.New("目前只支持mysql数据库")
	}

	err = dao.DB.Table("information_schema.columns").
		Where("table_schema= ? ", deployed.GenConfig.DBName).
		Where("TABLE_NAME = ?", tbName).
		Order("ORDINAL_POSITION asc").
		Find(&items).Error
	return items, err
}

func (cDBColumns) QueryPage(qp DBColumnsQueryParam) ([]DBColumns, paginator.Info, error) {
	var items []DBColumns

	db := dao.DB
	if dao.DbConfig.Dialect == "mysql" {
		if qp.TableName != "" {
			return nil, paginator.Info{}, errors.New("table name cannot be empty！")
		}
		db = db.Table("information_schema.`COLUMNS`").
			Where("table_schema= ? ", deployed.GenConfig.DBName).
			Where("TABLE_NAME = ?", qp.TableName)
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}
