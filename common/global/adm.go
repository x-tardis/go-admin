package global

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/glog"

	"github.com/x-tardis/go-admin/common/config"
)

var Cfg config.Conf = config.DefaultConfig()

var GinEngine *gin.Engine

var (
	Source string
	Driver string
	DBName string
)

var (
	Logger        *glog.Logger
	JobLogger     *glog.Logger
	RequestLogger *glog.Logger
)
