package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed/dao"
)

func InitDb(db *gorm.DB) (err error) {
	filePath := "script/db.sql"
	err = ExecSql(db, filePath)
	if dao.DbConfig.Dialect == "postgres" {
		filePath = "script/pg.sql"
		err = ExecSql(db, filePath)
	}
	return err
}

func ExecSql(db *gorm.DB, filePath string) error {
	sql, err := Ioutil(filePath)
	if err != nil {
		fmt.Println("数据库基础数据初始化脚本读取失败！原因:", err.Error())
		return err
	}
	sqlList := strings.Split(sql, ";")
	for i := 0; i < len(sqlList)-1; i++ {
		if strings.Contains(sqlList[i], "--") {
			fmt.Println(sqlList[i])
			continue
		}
		sql := strings.Replace(sqlList[i]+";", "\n", "", -1)
		sql = strings.TrimSpace(sql)
		if err = db.Exec(sql).Error; err != nil {
			log.Printf("error sql: %s", sql)
			if !strings.Contains(err.Error(), "Query was empty") {
				return err
			}
		}
	}
	return nil
}

func Ioutil(filePath string) (string, error) {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	// 因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
	return strings.Replace(string(contents), "\n", "", 1), nil
}
