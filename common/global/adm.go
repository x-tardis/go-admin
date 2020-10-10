package global

import (
	"github.com/gogf/gf/os/glog"

	"github.com/x-tardis/go-admin/common/config"
)

var Cfg config.Conf = config.DefaultConfig()

var Driver string

var (
	Logger        *glog.Logger
	JobLogger     *glog.Logger
	RequestLogger *glog.Logger
)
