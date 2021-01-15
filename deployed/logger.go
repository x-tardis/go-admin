package deployed

import (
	"github.com/thinkgos/x/lib/ternary"
	"go.uber.org/zap"

	"github.com/x-tardis/go-admin/pkg/izap"
)

var JobLogger *zap.SugaredLogger
var RequestLogger *zap.Logger

func SetupLogger() {
	c := ViperLogger()
	logger := izap.New(c)
	izap.ReplaceGlobals(logger)
	izap.Logger.Info("base logger init success")

	c.FileName = "job.log"
	c.Level = "info"
	c.Adapter = ternary.IfString(IsModeProd(), "file", "console")
	JobLogger = izap.New(c).Sugar()
	JobLogger.Info("job logger init success")

	c.FileName = "request.log"
	c.Level = "info"
	c.Adapter = ternary.IfString(IsModeProd(), "file", "multi")
	RequestLogger = izap.New(c)
	RequestLogger.Info("request logger init success")
}
