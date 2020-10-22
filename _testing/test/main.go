package main

import (
	"github.com/x-tardis/go-admin/deployed/dao"
)

func main() {
	dao.SetupDatabase("mysql", "root:catmaotu@tcp(127.0.0.1:3306)/goadmin?charset=utf8&parseTime=True&loc=Local&timeout=1000ms", true)

}
