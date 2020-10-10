package logger

import (
	"github.com/gogf/gf/os/glog"

	"github.com/x-tardis/go-admin/common/global"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/textcolor"
)

var Logger *glog.Logger
var JobLogger *glog.Logger
var RequestLogger *glog.Logger

func Setup() {
	Logger = glog.New()
	_ = Logger.SetPath(deployed.LoggerConfig.Path + "/bus")
	Logger.SetStdoutPrint(deployed.LoggerConfig.EnabledBUS && deployed.LoggerConfig.Stdout)
	Logger.SetFile("bus-{Ymd}.log")
	_ = Logger.SetLevelStr(deployed.LoggerConfig.Level)

	JobLogger = glog.New()
	_ = JobLogger.SetPath(deployed.LoggerConfig.Path + "/job")
	JobLogger.SetStdoutPrint(false)
	JobLogger.SetFile("db-{Ymd}.log")
	_ = JobLogger.SetLevelStr(deployed.LoggerConfig.Level)

	RequestLogger = glog.New()
	_ = RequestLogger.SetPath(deployed.LoggerConfig.Path + "/request")
	RequestLogger.SetStdoutPrint(false)
	RequestLogger.SetFile("access-{Ymd}.log")
	_ = RequestLogger.SetLevelStr(deployed.LoggerConfig.Level)

	Logger.Info(textcolor.Green("Logger init success!"))

	global.Logger = Logger.Line()
	global.JobLogger = JobLogger.Line()
	global.RequestLogger = RequestLogger.Line()
}
