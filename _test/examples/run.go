package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/common/global"
	mycasbin "github.com/x-tardis/go-admin/pkg/casbin"
	"github.com/x-tardis/go-admin/pkg/logger"
)

func main() {
	var err error
	global.Eloquent, err = gorm.Open(mysql.Open("root:123456@tcp/inmg?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	mycasbin.Setup()
	logger.Setup()
	global.GinEngine = gin.Default()
	//router.InitRouter()
	log.Fatal(global.GinEngine.Run(":8000"))
}
