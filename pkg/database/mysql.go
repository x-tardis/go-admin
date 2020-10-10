package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newMsql(source string) gorm.Dialector {
	// dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	c.UserName, c.Password, c.Protocol, c.Addr, c.DbName) // DSN data source name
	return mysql.New(mysql.Config{
		DSN: source,
		// DefaultStringSize:         256,   // string 类型字段的默认长度
		// DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		// DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		// DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		// SkipInitializeWithVersion: false, // 根据版本自动配置
	})
}
