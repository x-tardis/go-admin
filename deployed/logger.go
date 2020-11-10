package deployed

import (
	"go.uber.org/zap"

	"github.com/x-tardis/go-admin/pkg/izap"
)

var JobLogger *zap.SugaredLogger
var RequestLogger *zap.SugaredLogger

func SetupLogger() {
	c := ViperLogger()
	c.FileName = "bus.log"
	logger := izap.New(c)
	izap.ReplaceGlobals(logger)

	c.FileName = "job.log"
	c.InConsole = false
	JobLogger = izap.New(c).Sugar()

	c.FileName = "request.log"
	c.InConsole = true
	RequestLogger = izap.New(c).Sugar()
}
