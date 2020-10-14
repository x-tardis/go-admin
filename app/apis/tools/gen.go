package tools

import (
	"bytes"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/go-core-package/extos"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/app/models/tools"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/servers"
)

func Preview(c *gin.Context) {
	table := tools.SysTables{}
	id, err := strconv.Atoi(c.Param("tableId"))
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	table.TableId = id
	t1, err := template.ParseFiles("template/v3/model.go.template")
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	t2, err := template.ParseFiles("template/v3/no_actions/api.go.template")
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	t3, err := template.ParseFiles("template/v3/js.go.template")
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	t4, err := template.ParseFiles("template/v3/vue.go.template")
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	t5, err := template.ParseFiles("template/v3/no_actions/router_check_role.go.template")
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	t6, err := template.ParseFiles("template/v3/no_actions/dto.go.template")
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	tab, _ := table.Get()
	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	var b2 bytes.Buffer
	err = t2.Execute(&b2, tab)
	var b3 bytes.Buffer
	err = t3.Execute(&b3, tab)
	var b4 bytes.Buffer
	err = t4.Execute(&b4, tab)
	var b5 bytes.Buffer
	err = t5.Execute(&b5, tab)
	var b6 bytes.Buffer
	err = t6.Execute(&b6, tab)

	mp := make(map[string]interface{})
	mp["template/model.go.template"] = b1.String()
	mp["template/api.go.template"] = b2.String()
	mp["template/js.go.template"] = b3.String()
	mp["template/vue.go.template"] = b4.String()
	mp["template/router.go.template"] = b5.String()
	mp["template/dto.go.template"] = b6.String()

	servers.Success(c, servers.WithData(mp))
}

func GenCodeV3(c *gin.Context) {
	table := tools.SysTables{}
	id, err := strconv.Atoi(c.Param("tableId"))
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	table.TableId = id
	tab, _ := table.Get()

	if tab.IsActions == 1 {
		err = ActionsGenV3(tab)
	} else {
		err = NOActionsGenV3(tab)
	}
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, "", "Code generated successfully！")
}

func NOActionsGenV3(tab tools.SysTables) error {
	basePath := "template/v3/"
	routerFile := basePath + "no_actions/router_check_role.go.template"

	if tab.IsAuth == 2 {
		routerFile = basePath + "no_actions/router_no_check_role.go.template"
	}

	t1, err := template.ParseFiles(basePath + "model.go.template")
	if err != nil {
		return err
	}
	t2, err := template.ParseFiles(basePath + "no_actions/apis.go.template")
	if err != nil {
		return err
	}
	t3, err := template.ParseFiles(routerFile)
	if err != nil {
		return err
	}
	t4, err := template.ParseFiles(basePath + "js.go.template")
	if err != nil {
		return err
	}
	t5, err := template.ParseFiles(basePath + "vue.go.template")
	if err != nil {
		return err
	}
	t6, err := template.ParseFiles(basePath + "dto.go.template")
	if err != nil {
		return err
	}
	t7, err := template.ParseFiles(basePath + "no_actions/service.go.template")
	if err != nil {
		return err
	}

	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	var b2 bytes.Buffer
	err = t2.Execute(&b2, tab)
	var b3 bytes.Buffer
	err = t3.Execute(&b3, tab)
	var b4 bytes.Buffer
	err = t4.Execute(&b4, tab)
	var b5 bytes.Buffer
	err = t5.Execute(&b5, tab)
	var b6 bytes.Buffer
	err = t6.Execute(&b6, tab)
	var b7 bytes.Buffer
	err = t7.Execute(&b7, tab)
	extos.WriteFile("./app/"+tab.PackageName+"/models/"+tab.BusinessName+".go", b1.Bytes())
	extos.WriteFile("./app/"+tab.PackageName+"/apis/"+tab.ModuleName+"/"+tab.BusinessName+".go", b2.Bytes())
	extos.WriteFile("./app/"+tab.PackageName+"/router/"+tab.BusinessName+".go", b3.Bytes())
	extos.WriteFile(deployed.GenConfig.FrontPath+"/api/"+tab.BusinessName+".js", b4.Bytes())
	extos.WriteFile(deployed.GenConfig.FrontPath+"/views/"+tab.BusinessName+"/index.vue", b5.Bytes())
	extos.WriteFile("./app/"+tab.PackageName+"/service/dto/"+tab.BusinessName+".go", b6.Bytes())
	extos.WriteFile("./app/"+tab.PackageName+"/service/"+tab.BusinessName+".go", b7.Bytes())
	return nil
}

func ActionsGenV3(tab tools.SysTables) error {
	basePath := "template/v3/"
	routerFile := basePath + "actions/router_check_role.go.template"

	if tab.IsAuth == 2 {
		routerFile = basePath + "actions/router_no_check_role.go.template"
	}

	t1, err := template.ParseFiles(basePath + "model.go.template")
	if err != nil {
		return err
	}
	t3, err := template.ParseFiles(routerFile)
	if err != nil {
		return err
	}
	t4, err := template.ParseFiles(basePath + "js.go.template")
	if err != nil {
		return err
	}
	t5, err := template.ParseFiles(basePath + "vue.go.template")
	if err != nil {
		return err
	}
	t6, err := template.ParseFiles(basePath + "dto.go.template")
	if err != nil {
		return err
	}

	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	var b3 bytes.Buffer
	err = t3.Execute(&b3, tab)
	var b4 bytes.Buffer
	err = t4.Execute(&b4, tab)
	var b5 bytes.Buffer
	err = t5.Execute(&b5, tab)
	var b6 bytes.Buffer
	err = t6.Execute(&b6, tab)

	extos.WriteFile("./app/"+tab.PackageName+"/models/"+tab.BusinessName+".go", b1.Bytes())
	extos.WriteFile("./app/"+tab.PackageName+"/router/"+tab.BusinessName+".go", b3.Bytes())
	extos.WriteFile(deployed.GenConfig.FrontPath+"/api/"+tab.BusinessName+".js", b4.Bytes())
	extos.WriteFile(deployed.GenConfig.FrontPath+"/views/"+tab.BusinessName+"/index.vue", b5.Bytes())
	extos.WriteFile("./app/"+tab.PackageName+"/service/dto/"+tab.BusinessName+".go", b6.Bytes())
	return nil
}

func GenMenuAndApi(c *gin.Context) {
	table := tools.SysTables{}
	timeNow := time.Now()
	id, err := strconv.Atoi(c.Param("tableId"))
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	table.TableId = id
	tab, _ := table.Get()
	Mmenu := models.Menu{}
	Mmenu.MenuName = tab.TBName + "Manage"
	Mmenu.Title = tab.TableComment
	Mmenu.Icon = "pass"
	Mmenu.Path = "/" + tab.TBName
	Mmenu.MenuType = "M"
	Mmenu.Action = "无"
	Mmenu.ParentId = 0
	Mmenu.NoCache = false
	Mmenu.Component = "Layout"
	Mmenu.Sort = 0
	Mmenu.Visible = "0"
	Mmenu.IsFrame = "0"
	Mmenu.CreateBy = "1"
	Mmenu.UpdateBy = "1"
	Mmenu.CreatedAt = timeNow
	Mmenu.UpdatedAt = timeNow
	Mmenu.MenuId, err = Mmenu.Create()

	Cmenu := models.Menu{}
	Cmenu.MenuName = tab.TBName
	Cmenu.Title = tab.TableComment
	Cmenu.Icon = "pass"
	Cmenu.Path = tab.TBName
	Cmenu.MenuType = "C"
	Cmenu.Action = "无"
	Cmenu.Permission = tab.ModuleName + ":" + tab.BusinessName + ":list"
	Cmenu.ParentId = Mmenu.MenuId
	Cmenu.NoCache = false
	Cmenu.Component = "/" + tab.BusinessName + "/index"
	Cmenu.Sort = 0
	Cmenu.Visible = "0"
	Cmenu.IsFrame = "0"
	Cmenu.CreateBy = "1"
	Cmenu.UpdateBy = "1"
	Cmenu.CreatedAt = timeNow
	Cmenu.UpdatedAt = timeNow
	Cmenu.MenuId, err = Cmenu.Create()

	MList := models.Menu{}
	MList.MenuName = ""
	MList.Title = "分页获取" + tab.TableComment
	MList.Icon = ""
	MList.Path = tab.TBName
	MList.MenuType = "F"
	MList.Action = "无"
	MList.Permission = tab.ModuleName + ":" + tab.BusinessName + ":query"
	MList.ParentId = Cmenu.MenuId
	MList.NoCache = false
	MList.Sort = 0
	MList.Visible = "0"
	MList.IsFrame = "0"
	MList.CreateBy = "1"
	MList.UpdateBy = "1"
	MList.CreatedAt = timeNow
	MList.UpdatedAt = timeNow
	MList.MenuId, err = MList.Create()

	MCreate := models.Menu{}
	MCreate.MenuName = ""
	MCreate.Title = "创建" + tab.TableComment
	MCreate.Icon = ""
	MCreate.Path = tab.TBName
	MCreate.MenuType = "F"
	MCreate.Action = "无"
	MCreate.Permission = tab.ModuleName + ":" + tab.BusinessName + ":add"
	MCreate.ParentId = Cmenu.MenuId
	MCreate.NoCache = false
	MCreate.Sort = 0
	MCreate.Visible = "0"
	MCreate.IsFrame = "0"
	MCreate.CreateBy = "1"
	MCreate.UpdateBy = "1"
	MCreate.CreatedAt = timeNow
	MCreate.UpdatedAt = timeNow
	MCreate.MenuId, err = MCreate.Create()

	MUpdate := models.Menu{}
	MUpdate.MenuName = ""
	MUpdate.Title = "修改" + tab.TableComment
	MUpdate.Icon = ""
	MUpdate.Path = tab.TBName
	MUpdate.MenuType = "F"
	MUpdate.Action = "无"
	MUpdate.Permission = tab.ModuleName + ":" + tab.BusinessName + ":edit"
	MUpdate.ParentId = Cmenu.MenuId
	MUpdate.NoCache = false
	MUpdate.Sort = 0
	MUpdate.Visible = "0"
	MUpdate.IsFrame = "0"
	MUpdate.CreateBy = "1"
	MUpdate.UpdateBy = "1"
	MUpdate.CreatedAt = timeNow
	MUpdate.UpdatedAt = timeNow
	MUpdate.MenuId, err = MUpdate.Create()

	MDelete := models.Menu{}
	MDelete.MenuName = ""
	MDelete.Title = "删除" + tab.TableComment
	MDelete.Icon = ""
	MDelete.Path = tab.TBName
	MDelete.MenuType = "F"
	MDelete.Action = "无"
	MDelete.Permission = tab.ModuleName + ":" + tab.BusinessName + ":remove"
	MDelete.ParentId = Cmenu.MenuId
	MDelete.NoCache = false
	MDelete.Sort = 0
	MDelete.Visible = "0"
	MDelete.IsFrame = "0"
	MDelete.CreateBy = "1"
	MDelete.UpdateBy = "1"
	MDelete.CreatedAt = timeNow
	MDelete.UpdatedAt = timeNow
	MDelete.MenuId, err = MDelete.Create()

	var InterfaceId = 63
	Amenu := models.Menu{}
	Amenu.MenuName = tab.TBName
	Amenu.Title = tab.TableComment
	Amenu.Icon = "bug"
	Amenu.Path = tab.TBName
	Amenu.MenuType = "M"
	Amenu.Action = "无"
	Amenu.ParentId = InterfaceId
	Amenu.NoCache = false
	Amenu.Sort = 0
	Amenu.Visible = "1"
	Amenu.IsFrame = "0"
	Amenu.CreateBy = "1"
	Amenu.UpdateBy = "1"
	Amenu.CreatedAt = timeNow
	Amenu.UpdatedAt = timeNow
	Amenu.MenuId, err = Amenu.Create()

	AList := models.Menu{}
	AList.MenuName = ""
	AList.Title = "分页获取" + tab.TableComment
	AList.Icon = "bug"
	AList.Path = "/api/v1/" + tab.ModuleName
	AList.MenuType = "A"
	AList.Action = "GET"
	AList.ParentId = Amenu.MenuId
	AList.NoCache = false
	AList.Sort = 0
	AList.Visible = "1"
	AList.IsFrame = "0"
	AList.CreateBy = "1"
	AList.UpdateBy = "1"
	AList.CreatedAt = timeNow
	AList.UpdatedAt = timeNow
	AList.MenuId, err = AList.Create()

	AGet := models.Menu{}
	AGet.MenuName = ""
	AGet.Title = "根据id获取" + tab.TableComment
	AGet.Icon = "bug"
	AGet.Path = "/api/v1/" + tab.ModuleName + "/:id"
	AGet.MenuType = "A"
	AGet.Action = "GET"
	AGet.ParentId = Amenu.MenuId
	AGet.NoCache = false
	AGet.Sort = 0
	AGet.Visible = "1"
	AGet.IsFrame = "0"
	AGet.CreateBy = "1"
	AGet.UpdateBy = "1"
	AGet.CreatedAt = timeNow
	AGet.UpdatedAt = timeNow
	AGet.MenuId, err = AGet.Create()

	ACreate := models.Menu{}
	ACreate.MenuName = ""
	ACreate.Title = "创建" + tab.TableComment
	ACreate.Icon = "bug"
	ACreate.Path = "/api/v1/" + tab.ModuleName
	ACreate.MenuType = "A"
	ACreate.Action = "POST"
	ACreate.ParentId = Amenu.MenuId
	ACreate.NoCache = false
	ACreate.Sort = 0
	ACreate.Visible = "1"
	ACreate.IsFrame = "0"
	ACreate.CreateBy = "1"
	ACreate.UpdateBy = "1"
	ACreate.CreatedAt = timeNow
	ACreate.UpdatedAt = timeNow
	ACreate.MenuId, err = ACreate.Create()

	AUpdate := models.Menu{}
	AUpdate.MenuName = ""
	AUpdate.Title = "修改" + tab.TableComment
	AUpdate.Icon = "bug"
	AUpdate.Path = "/api/v1/" + tab.ModuleName + "/:id"
	AUpdate.MenuType = "A"
	AUpdate.Action = "PUT"
	AUpdate.ParentId = Amenu.MenuId
	AUpdate.NoCache = false
	AUpdate.Sort = 0
	AUpdate.Visible = "1"
	AUpdate.IsFrame = "0"
	AUpdate.CreateBy = "1"
	AUpdate.UpdateBy = "1"
	AUpdate.CreatedAt = timeNow
	AUpdate.UpdatedAt = timeNow
	AUpdate.MenuId, err = AUpdate.Create()

	ADelete := models.Menu{}
	ADelete.MenuName = ""
	ADelete.Title = "删除" + tab.TableComment
	ADelete.Icon = "bug"
	ADelete.Path = "/api/v1/" + tab.ModuleName
	ADelete.MenuType = "A"
	ADelete.Action = "DELETE"
	ADelete.ParentId = Amenu.MenuId
	ADelete.NoCache = false
	ADelete.Sort = 0
	ADelete.Visible = "1"
	ADelete.IsFrame = "0"
	ADelete.CreateBy = "1"
	ADelete.UpdateBy = "1"
	ADelete.CreatedAt = timeNow
	ADelete.UpdatedAt = timeNow
	ADelete.MenuId, err = ADelete.Create()

	servers.OKWithRequestID(c, "", "数据生成成功！")
}