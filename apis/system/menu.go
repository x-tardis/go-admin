package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

type Menu struct{}

// @tags 菜单/Menu
// @summary 获取Menu树
// @description 获取Menu树
// @security Bearer
// @accept json
// @produce json
// @param title query string false "title"
// @param menuName query string false "menuName"
// @param path query string false "path"
// @param action query string false "action"
// @param menuType query string false "menuType"
// @param visible query string false "visible"
// @success 200 {string} string "{"code": 200, "data": [...]}"
// @success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @router /api/v1/menus [get]
func (Menu) QueryTree(c *gin.Context) {
	qp := models.MenuQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	tree, err := models.CMenu.QueryTree(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithMsg(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(tree))
}

// @tags 菜单/Menu
// @summary 获取Menu数据
// @description 获取Menu数据
// @security Bearer
// @accept json
// @produce json
// @param id path int true "id"
// @success 200 {string} string "{"code": 200, "data": [...]}"
// @success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @router /api/v1/menus/{id} [get]
func (Menu) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CMenu.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithMsg(prompt.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 菜单/Menu
// @summary 创建menu
// @description 创建menu
// @security Bearer
// @accept json
// @produce json
// @param newItem body models.Menu true "new item"
// @success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @router /api/v1/menus [post]
func (Menu) Create(c *gin.Context) {
	newItem := models.Menu{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CMenu.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @tags 菜单/Menu
// @summary 修改菜单
// @description 获取JSON
// @security Bearer
// @accept json
// @produce json
// @param up body models.Menu true "update item"
// @success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @router /api/v1/menus [put]
func (Menu) Update(c *gin.Context) {
	up := models.Menu{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	err := models.CMenu.Update(gcontext.Context(c), up.MenuId, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(prompt.UpdateFailed))
		return
	}
	servers.OK(c)
}

// @tags 菜单/Menu
// @summary 删除菜单
// @description 删除菜单
// @security Bearer
// @accept json
// @produce json
// @param id path int true "id"
// @success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @router /api/v1/menus/{id} [delete]
func (Menu) Delete(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	err := models.CMenu.Delete(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithMsg(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithMsg(prompt.DeleteSuccess))
}

// @tags 菜单/Menu
// @Summary 获取菜单树
// @Description 获取JSON
// @Accept  application/x-www-form-urlencoded
// @Product application/x-www-form-urlencoded
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/menuTree/option [get]
// @Security Bearer
func (Menu) GetMenuTreeOption(c *gin.Context) {
	items, err := models.CMenu.QueryTitleLabelTree(gcontext.Context(c))
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithMsg(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(items))
}

// @tags 菜单/Menu
// @Summary 根据角色名称获取菜单列表数据（左菜单使用）
// @Description 获取JSON
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menuTree/role [get]
// @Security Bearer
func (Menu) GetMenuTreeWithRoleName(c *gin.Context) {
	items, err := models.CMenu.QueryTreeWithRoleName(gcontext.Context(c))
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithMsg(prompt.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(items))
}
