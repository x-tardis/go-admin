package tools

import (
	"bytes"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/go-core-package/extos"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/app/models/tools"
	"github.com/x-tardis/go-admin/pkg/deployed"
	"github.com/x-tardis/go-admin/pkg/servers"
)

func Preview(c *gin.Context) {
	table := tools.SysTables{}
	id, err := strconv.Atoi(c.Param("tableId"))
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	table.TableId = id
	t1, err := template.ParseFiles("template/v3/model.go.template")
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	t2, err := template.ParseFiles("template/v3/no_actions/api.go.template")
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	t3, err := template.ParseFiles("template/v3/js.go.template")
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	t4, err := template.ParseFiles("template/v3/vue.go.template")
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	t5, err := template.ParseFiles("template/v3/no_actions/router_check_role.go.template")
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	t6, err := template.ParseFiles("template/v3/no_actions/dto.go.template")
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
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

	servers.JSON(c, http.StatusOK, servers.WithData(mp))
}

func GenCodeV3(c *gin.Context) {
	table := tools.SysTables{}
	id, err := strconv.Atoi(c.Param("tableId"))
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
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
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithMsg("Code generated successfully！"))
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
	id, err := strconv.Atoi(c.Param("tableId"))
	if err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	now := time.Now()
	table := tools.SysTables{
		TableId: id,
	}
	tab, _ := table.Get()

	Mmenu := models.Menu{
		MenuName:  tab.TBName + "Manage",
		Title:     tab.TableComment,
		Icon:      "pass",
		Path:      "/" + tab.TBName,
		MenuType:  "M",
		Action:    "无",
		ParentId:  0,
		NoCache:   false,
		Component: "Layout",
		Sort:      0,
		Visible:   "0",
		IsFrame:   "0",
		Creator:   "1",
		Updator:   "1",
		Model: models.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	Mmenu, err = models.CMenu.Create(gcontext.Context(c), Mmenu)

	Cmenu := models.Menu{
		MenuName:   tab.TBName,
		Title:      tab.TableComment,
		Icon:       "pass",
		Path:       tab.TBName,
		MenuType:   "C",
		Action:     "无",
		Permission: tab.ModuleName + ":" + tab.BusinessName + ":list",
		ParentId:   Mmenu.MenuId,
		NoCache:    false,
		Component:  "/" + tab.BusinessName + "/index",
		Sort:       0,
		Visible:    "0",
		IsFrame:    "0",
		Creator:    "1",
		Updator:    "1",
		Model: models.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	Cmenu, err = models.CMenu.Create(gcontext.Context(c), Cmenu)

	MList := models.Menu{
		MenuName:   "",
		Title:      "分页获取" + tab.TableComment,
		Icon:       "",
		Path:       tab.TBName,
		MenuType:   "F",
		Action:     "无",
		Permission: tab.ModuleName + ":" + tab.BusinessName + ":query",
		ParentId:   Cmenu.MenuId,
		NoCache:    false,
		Sort:       0,
		Visible:    "0",
		IsFrame:    "0",
		Creator:    "1",
		Updator:    "1",
		Model: models.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	MList, err = models.CMenu.Create(gcontext.Context(c), MList)

	MCreate := models.Menu{
		MenuName:   "",
		Title:      "创建" + tab.TableComment,
		Icon:       "",
		Path:       tab.TBName,
		MenuType:   "F",
		Action:     "无",
		Permission: tab.ModuleName + ":" + tab.BusinessName + ":add",
		ParentId:   Cmenu.MenuId,
		NoCache:    false,
		Sort:       0,
		Visible:    "0",
		IsFrame:    "0",
		Creator:    "1",
		Updator:    "1",
		Model: models.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	MCreate, err = models.CMenu.Create(gcontext.Context(c), MCreate)

	MUpdate := models.Menu{
		MenuName:   "",
		Title:      "修改" + tab.TableComment,
		Icon:       "",
		Path:       tab.TBName,
		MenuType:   "F",
		Action:     "无",
		Permission: tab.ModuleName + ":" + tab.BusinessName + ":edit",
		ParentId:   Cmenu.MenuId,
		NoCache:    false,
		Sort:       0,
		Visible:    "0",
		IsFrame:    "0",
		Creator:    "1",
		Updator:    "1",
		Model: models.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	MUpdate, err = models.CMenu.Create(gcontext.Context(c), MUpdate)

	MDelete := models.Menu{
		MenuName:   "",
		Title:      "删除" + tab.TableComment,
		Icon:       "",
		Path:       tab.TBName,
		MenuType:   "F",
		Action:     "无",
		Permission: tab.ModuleName + ":" + tab.BusinessName + ":remove",
		ParentId:   Cmenu.MenuId,
		NoCache:    false,
		Sort:       0,
		Visible:    "0",
		IsFrame:    "0",
		Creator:    "1",
		Updator:    "1",
		Model: models.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	MDelete, err = models.CMenu.Create(gcontext.Context(c), MDelete)

	var InterfaceId = 63
	Amenu := models.Menu{
		MenuName: tab.TBName,
		Title:    tab.TableComment,
		Icon:     "bug",
		Path:     tab.TBName,
		MenuType: "M",
		Action:   "无",
		ParentId: InterfaceId,
		NoCache:  false,
		Sort:     0,
		Visible:  "1",
		IsFrame:  "0",
		Creator:  "1",
		Updator:  "1",
		Model: models.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	Amenu, err = models.CMenu.Create(gcontext.Context(c), Amenu)

	AList := models.Menu{
		MenuName: "",
		Title:    "分页获取" + tab.TableComment,
		Icon:     "bug",
		Path:     "/api/v1/" + tab.ModuleName,
		MenuType: "A",
		Action:   "GET",
		ParentId: Amenu.MenuId,
		NoCache:  false,
		Sort:     0,
		Visible:  "1",
		IsFrame:  "0",
		Creator:  "1",
		Updator:  "1",
		Model: models.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	AList, err = models.CMenu.Create(gcontext.Context(c), AList)

	AGet := models.Menu{
		MenuName: "",
		Title:    "根据id获取" + tab.TableComment,
		Icon:     "bug",
		Path:     "/api/v1/" + tab.ModuleName + "/:id",
		MenuType: "A",
		Action:   "GET",
		ParentId: Amenu.MenuId,
		NoCache:  false,
		Sort:     0,
		Visible:  "1",
		IsFrame:  "0",
		Creator:  "1",
		Updator:  "1",
		Model: models.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	AGet, err = models.CMenu.Create(gcontext.Context(c), AGet)

	ACreate := models.Menu{
		MenuName: "",
		Title:    "创建" + tab.TableComment,
		Icon:     "bug",
		Path:     "/api/v1/" + tab.ModuleName,
		MenuType: "A",
		Action:   "POST",
		ParentId: Amenu.MenuId,
		NoCache:  false,
		Sort:     0,
		Visible:  "1",
		IsFrame:  "0",
		Creator:  "1",
		Updator:  "1",
		Model: models.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	ACreate, err = models.CMenu.Create(gcontext.Context(c), ACreate)

	AUpdate := models.Menu{
		MenuName: "",
		Title:    "修改" + tab.TableComment,
		Icon:     "bug",
		Path:     "/api/v1/" + tab.ModuleName + "/:id",
		MenuType: "A",
		Action:   "PUT",
		ParentId: Amenu.MenuId,
		NoCache:  false,
		Sort:     0,
		Visible:  "1",
		IsFrame:  "0",
		Creator:  "1",
		Updator:  "1",
		Model: models.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	AUpdate, err = models.CMenu.Create(gcontext.Context(c), AUpdate)

	ADelete := models.Menu{
		MenuName: "",
		Title:    "删除" + tab.TableComment,
		Icon:     "bug",
		Path:     "/api/v1/" + tab.ModuleName,
		MenuType: "A",
		Action:   "DELETE",
		ParentId: Amenu.MenuId,
		NoCache:  false,
		Sort:     0,
		Visible:  "1",
		IsFrame:  "0",
		Creator:  "1",
		Updator:  "1",
		Model: models.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	ADelete, err = models.CMenu.Create(gcontext.Context(c), ADelete)

	servers.JSON(c, http.StatusOK, servers.WithMsg("数据生成成功！"))
}
