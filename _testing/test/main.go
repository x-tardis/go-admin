package main

import (
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/izap"
)

func main() {
	deployed.SetupLogger()
	izap.Sugar.Errorf("Logger")
	deployed.JobLogger.Errorf("JobLogger")
	deployed.RequestLogger.Errorf("RequestLogger")
}
