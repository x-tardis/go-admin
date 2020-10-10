package deployed

import (
	"github.com/gogf/gf/os/glog"

	"github.com/x-tardis/go-admin/pkg/textcolor"
)

var (
	Logger        *glog.Logger
	JobLogger     *glog.Logger
	RequestLogger *glog.Logger
)

func SetupLogger() {
	log := glog.New()
	_ = log.SetPath(LoggerConfig.Path + "/bus")
	log.SetStdoutPrint(LoggerConfig.EnabledBUS && LoggerConfig.Stdout)
	log.SetFile("bus-{Ymd}.log")
	_ = log.SetLevelStr(LoggerConfig.Level)

	JobLog := glog.New()
	_ = JobLog.SetPath(LoggerConfig.Path + "/job")
	JobLog.SetStdoutPrint(false)
	JobLog.SetFile("db-{Ymd}.log")
	_ = JobLog.SetLevelStr(LoggerConfig.Level)

	RequestLog := glog.New()
	_ = RequestLog.SetPath(LoggerConfig.Path + "/request")
	RequestLog.SetStdoutPrint(false)
	RequestLog.SetFile("access-{Ymd}.log")
	_ = RequestLog.SetLevelStr(LoggerConfig.Level)

	log.Info(textcolor.Green("Logger init success!"))

	Logger = log.Line()
	JobLogger = JobLog.Line()
	RequestLogger = RequestLog.Line()
}
