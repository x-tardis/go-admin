package main

import (
	"github.com/x-tardis/go-admin/pkg/deployed"
)

func main() {
	deployed.SetupLogger()
	deployed.Logger.Errorf("Logger")
	deployed.JobLogger.Errorf("JobLogger")
	deployed.RequestLogger.Errorf("RequestLogger")
}
