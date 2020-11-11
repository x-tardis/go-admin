package main

import (
	"github.com/x-tardis/go-admin/cmd"

	_ "github.com/x-tardis/go-admin/docs"
)

// @title go-admin API
// @version 1.0.1
// @description 基于Gin + Vue + Element UI的前后端分离权限管理系统的接口文档
// @license.name MIT
// @license.url https://github.com/x-tardis/go-admin/blob/master/LICENSE.md

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
