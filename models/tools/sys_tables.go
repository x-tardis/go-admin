package tools

import (
	"context"
	"strings"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/thinkgos/sharp/iorm"
	"github.com/thinkgos/sharp/iorm/trans"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
)

type SysTables struct {
	TableId             int    `gorm:"primary_key;auto_increment" json:"tableId"`    // 表编码
	TBName              string `gorm:"column:table_name;size:255;" json:"tableName"` // 表名称
	TableComment        string `gorm:"size:255;" json:"tableComment"`                // 表备注
	ClassName           string `gorm:"size:255;" json:"className"`                   // 类名
	TplCategory         string `gorm:"size:255;" json:"tplCategory"`                 //
	PackageName         string `gorm:"size:255;" json:"packageName"`                 // 包名
	ModuleName          string `gorm:"size:255;" json:"moduleName"`                  // 模块名
	BusinessName        string `gorm:"size:255;" json:"businessName"`                //
	FunctionName        string `gorm:"size:255;" json:"functionName"`                // 功能名称
	FunctionAuthor      string `gorm:"size:255;" json:"functionAuthor"`              // 功能作者
	PkColumn            string `gorm:"size:255;" json:"pkColumn"`
	PkGoField           string `gorm:"size:255;" json:"pkGoField"`
	PkJsonField         string `gorm:"size:255;" json:"pkJsonField"`
	Options             string `gorm:"size:255;" json:"options"`
	TreeCode            string `gorm:"size:255;" json:"treeCode"`
	TreeParentCode      string `gorm:"size:255;" json:"treeParentCode"`
	TreeName            string `gorm:"size:255;" json:"treeName"`
	Tree                bool   `gorm:"size:1;" json:"tree"`
	Crud                bool   `gorm:"size:1;" json:"crud"`
	Remark              string `gorm:"size:255;" json:"remark"`
	IsDataScope         int    `gorm:"size:1;" json:"isDataScope"`
	IsActions           int    `gorm:"size:1;" json:"isActions"`
	IsAuth              int    `gorm:"size:1;" json:"isAuth"`
	IsLogicalDelete     string `gorm:"size:1;" json:"isLogicalDelete"`
	LogicalDelete       bool   `gorm:"size:1;" json:"logicalDelete"`
	LogicalDeleteColumn string `gorm:"size:128;" json:"logicalDeleteColumn"`
	Creator             string `gorm:"size:128;" json:"creator"`
	Updator             string `gorm:"size:128;" json:"updator"`
	models.Model

	DataScope string       `gorm:"-" json:"dataScope"`
	Params    Params       `gorm:"-" json:"params"`
	Columns   []SysColumns `gorm:"-" json:"columns"`
}

type Params struct {
	TreeCode       string `gorm:"-" json:"treeCode"`
	TreeParentCode string `gorm:"-" json:"treeParentCode"`
	TreeName       string `gorm:"-" json:"treeName"`
}

func (SysTables) TableName() string {
	return "sys_tables"
}

func TablesDB(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(trans.CtxDB(ctx)).Model(SysTables{})
	}
}

type TablesQueryParam struct {
	TableName    string `form:"tableName"`
	TableComment string `form:"tableComment"`
	paginator.Param
}

type cTables struct{}

var CTables = cTables{}

func (cTables) QueryPage(ctx context.Context, qp TablesQueryParam) ([]SysTables, paginator.Info, error) {
	var items []SysTables

	db := dao.DB.Scopes(TablesDB(ctx))
	if qp.TableName != "" {
		db = db.Where("table_name=?", qp.TableName)
	}
	if qp.TableComment != "" {
		db = db.Where("table_comment=?", qp.TableComment)
	}

	info, err := iorm.QueryPages(db, qp.Param, &items)
	return items, info, err
}

func (cTables) Get(ctx context.Context, id int) (item SysTables, err error) {
	err = dao.DB.Scopes(TablesDB(ctx)).
		Where("table_id=?", id).First(&item).Error
	if err != nil {
		return
	}
	item.Columns, err = CColumns.QueryWithTableId(ctx, id)
	return
}

func (cTables) GetWithName(ctx context.Context, tbName string) (item SysTables, err error) {
	err = dao.DB.Scopes(TablesDB(ctx)).
		Where("table_name=?", tbName).
		First(&item).Error
	if err != nil {
		return item, err
	}
	item.Columns, err = CColumns.QueryWithTableId(ctx, item.TableId)
	return item, err
}

func (cTables) QueryTree(ctx context.Context) (items []SysTables, err error) {
	err = dao.DB.Scopes(TablesDB(ctx)).Find(&items).Error
	if err != nil {
		return
	}
	for i := 0; i < len(items); i++ {
		items[i].Columns, err = CColumns.QueryWithTableId(ctx, items[i].TableId)
		if err != nil {
			return items, err
		}
	}
	return items, nil
}

func (cTables) Create(ctx context.Context, item SysTables) (SysTables, error) {
	err := dao.DB.Scopes(TablesDB(ctx)).Create(&item).Error
	if err != nil {
		return item, err
	}
	for i := 0; i < len(item.Columns); i++ {
		item.Columns[i].TableId = item.TableId
		item.Columns[i], _ = CColumns.Create(ctx, item.Columns[i])
	}
	return item, nil
}

func (cTables) Update(ctx context.Context, id int, up SysTables) (item SysTables, err error) {
	// if err = dao.DB.Scopes(TablesDB(ctx)).First(&item, id).Error; err != nil {
	// 	return
	// }

	up.Updator = jwtauth.FromUserIdStr(ctx)
	err = dao.DB.Scopes(TablesDB(ctx)).
		Where("table_id = ?", id).Updates(&up).Error
	if err != nil {
		return
	}

	tableNames := make([]string, 0)
	for i := range up.Columns {
		if up.Columns[i].FkTableName != "" {
			tableNames = append(tableNames, up.Columns[i].FkTableName)
		}
	}

	tables := make([]SysTables, 0)
	tableMap := make(map[string]*SysTables)
	if len(tableNames) > 0 {
		err = dao.DB.Scopes(TablesDB(ctx)).
			Where("table_name in (?)", tableNames).Find(&tables).Error
		if err != nil {
			return
		}
		for i := range tables {
			tableMap[tables[i].TBName] = &tables[i]
		}
	}

	for i := 0; i < len(up.Columns); i++ {
		if up.Columns[i].FkTableName != "" {
			t, ok := tableMap[up.Columns[i].FkTableName]
			if ok {
				up.Columns[i].FkTableNameClass = t.ClassName
				up.Columns[i].FkTableNamePackage = t.ModuleName
			} else {
				tableNameList := strings.Split(up.Columns[i].FkTableName, "_")
				up.Columns[i].FkTableNameClass = ""
				up.Columns[i].FkTableNamePackage = ""
				for a := 0; a < len(tableNameList); a++ {
					strStart := string([]byte(tableNameList[a])[:1])
					strEnd := string([]byte(tableNameList[a])[1:])
					up.Columns[i].FkTableNameClass += strings.ToUpper(strStart) + strEnd
					up.Columns[i].FkTableNamePackage += strings.ToLower(strStart) + strings.ToLower(strEnd)
				}
			}
		}
		up.Columns[i], _ = CColumns.Update(ctx, up.Columns[i].ColumnId, up.Columns[i])
	}
	return
}

func (cTables) Delete(ctx context.Context, id int) error {
	return trans.Exec(ctx, dao.DB, func(ctx context.Context) error {
		err := dao.DB.Scopes(TablesDB(ctx)).
			Delete(SysTables{}, "table_id=?", id).Error
		if err != nil {
			return err
		}
		return dao.DB.Scopes(ColumnsDB(ctx)).
			Delete(SysColumns{}, "table_id=?", id).Error
	})
}

func (cTables) BatchDelete(ctx context.Context, ids []int) error {
	return dao.DB.Unscoped().Scopes(TablesDB(ctx)).
		Where(" table_id in (?)", ids).Delete(&SysColumns{}).Error
}
