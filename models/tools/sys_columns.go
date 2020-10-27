package tools

import (
	"context"

	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/models"
)

type SysColumns struct {
	ColumnId           int          `gorm:"primary_key;auto_increment" json:"columnId"`
	TableId            int          `gorm:"" json:"tableId"`
	ColumnName         string       `gorm:"size:128;" json:"columnName"`
	ColumnComment      string       `gorm:"column:column_comment;size:128;" json:"columnComment"`
	ColumnType         string       `gorm:"column:column_type;size:128;" json:"columnType"`
	GoType             string       `gorm:"column:go_type;size:128;" json:"goType"`
	GoField            string       `gorm:"column:go_field;size:128;" json:"goField"`
	JsonField          string       `gorm:"column:json_field;size:128;" json:"jsonField"`
	IsPk               string       `gorm:"column:is_pk;size:4;" json:"isPk"`
	IsIncrement        string       `gorm:"column:is_increment;size:4;" json:"isIncrement"`
	IsRequired         string       `gorm:"column:is_required;size:4;" json:"isRequired"`
	IsInsert           string       `gorm:"column:is_insert;size:4;" json:"isInsert"`
	IsEdit             string       `gorm:"column:is_edit;size:4;" json:"isEdit"`
	IsList             string       `gorm:"column:is_list;size:4;" json:"isList"`
	IsQuery            string       `gorm:"column:is_query;size:4;" json:"isQuery"`
	QueryType          string       `gorm:"column:query_type;size:128;" json:"queryType"`
	HtmlType           string       `gorm:"column:html_type;size:128;" json:"htmlType"`
	DictType           string       `gorm:"column:dict_type;size:128;" json:"dictType"`
	Sort               int          `gorm:"column:sort;" json:"sort"`
	List               string       `gorm:"column:list;size:1;" json:"list"`
	Pk                 bool         `gorm:"column:pk;size:1;" json:"pk"`
	Required           bool         `gorm:"column:required;size:1;" json:"required"`
	SuperColumn        bool         `gorm:"column:super_column;size:1;" json:"superColumn"`
	UsableColumn       bool         `gorm:"column:usable_column;size:1;" json:"usableColumn"`
	Increment          bool         `gorm:"column:increment;size:1;" json:"increment"`
	Insert             bool         `gorm:"column:insert;size:1;" json:"insert"`
	Edit               bool         `gorm:"column:edit;size:1;" json:"edit"`
	Query              bool         `gorm:"column:query;size:1;" json:"query"`
	Remark             string       `gorm:"column:remark;size:255;" json:"remark"`
	FkTableName        string       `gorm:"" json:"fkTableName"`
	FkTableNameClass   string       `gorm:"" json:"fkTableNameClass"`
	FkTableNamePackage string       `gorm:"" json:"fkTableNamePackage"`
	FkCol              []SysColumns `gorm:"-" json:"fkCol"`
	FkLabelId          string       `gorm:"" json:"fkLabelId"`
	FkLabelName        string       `gorm:"size:255;" json:"fkLabelName"`
	Creator            string       `gorm:"column:creator;size:128;" json:"creator"`
	Updator            string       `gorm:"column:updator;size:128;" json:"updator"`

	models.Model
}

func (SysColumns) TableName() string {
	return "sys_columns"
}

func ColumnsDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(SysColumns{})
	}
}

type cColumns struct{}

var CColumns = cColumns{}

func (cColumns) QueryWithTableId(ctx context.Context, id int) (items []SysColumns, err error) {
	err = dao.DB.Scopes(ColumnsDB(ctx)).
		Where("table_id=?", id).
		Find(&items).Error
	return items, err
}

func (cColumns) Create(ctx context.Context, item SysColumns) (SysColumns, error) {
	err := dao.DB.Scopes(ColumnsDB(ctx)).Create(&item).Error
	return item, err
}

func (cColumns) Update(ctx context.Context, id int, up SysColumns) (item SysColumns, err error) {
	if err = dao.DB.Scopes(ColumnsDB(ctx)).First(&item, id).Error; err != nil {
		return
	}

	err = dao.DB.Table("sys_columns").Model(&item).Updates(&up).Error
	return
}
