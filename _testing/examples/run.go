package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/deployed"
	"github.com/x-tardis/go-admin/deployed/dao"
)

func main() {
	var err error
	dao.DB, err = gorm.Open(mysql.Open("root:123456@tcp/inmg?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	deployed.SetupCasbin()
	deployed.SetupLogger()
	engine := gin.Default()
	//router.InitRouter()
	log.Fatal(engine.Run(":8000"))
}
