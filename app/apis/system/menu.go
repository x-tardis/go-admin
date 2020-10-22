package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/models"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

type Menu struct{}

// @Summary Menu列表数据
// @Description 获取Menu tree
// @Tags 菜单
// @Param menuName query string false "menuName"
// @Param visible query string false "visible"
// @Param title query string false "title"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menus [get]
// @Security Bearer
func (Menu) QueryTree(c *gin.Context) {
	qp := models.MenuQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	tree, err := models.CMenu.QueryTree(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(tree))
}

// @Summary Menu列表数据
// @Description 获取JSON
// @Tags 菜单
// @Param menuName query string false "menuName"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menus/{id} [get]
// @Security Bearer
func (Menu) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CMenu.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

// @Tags 菜单
// @Summary 创建菜单
// @Description 获取JSON
// @Accept  application/x-www-form-urlencoded
// @Product application/x-www-form-urlencoded
// @Param menuName formData string true "menuName"
// @Param Path formData string false "Path"
// @Param Action formData string true "Action"
// @Param Permission formData string true "Permission"
// @Param ParentId formData string true "ParentId"
// @Param IsDel formData string true "IsDel"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/menus [post]
// @Security Bearer
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

// @Tags 菜单
// @Summary 修改菜单
// @Description 获取JSON
// @Accept  application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body models.Menu true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/menus [put]
// @Security Bearer
func (Menu) Update(c *gin.Context) {
	up := models.Menu{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	_, err := models.CMenu.Update(gcontext.Context(c), up.MenuId, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.UpdateFailed))
		return
	}
	servers.OK(c)
}

// @Tags 菜单
// @Summary 删除菜单
// @Description 删除菜单
// @Param id path int true "id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/menus/{id} [delete]
func (Menu) Delete(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	err := models.CMenu.Delete(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithPrompt(prompt.DeleteSuccess))
}

func GetMenuTreeRoleselect(c *gin.Context) {
	result, err := models.CMenu.QueryLabelTree(gcontext.Context(c))
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	id := cast.ToInt(c.Param("roleId"))
	menuIds := make([]int, 0)
	if id != 0 {
		menuIds, err = models.CRole.GetMenuIds(id)
		if err != nil {
			servers.Fail(c, http.StatusInternalServerError,
				servers.WithPrompt(prompt.QueryFailed),
				servers.WithError(err))
			return
		}
	}
	servers.JSON(c, http.StatusOK, gin.H{
		"code":        200,
		"menus":       result,
		"checkedKeys": menuIds,
	})
}

// @Summary 获取菜单树
// @Description 获取JSON
// @Tags 菜单
// @Accept  application/x-www-form-urlencoded
// @Product application/x-www-form-urlencoded
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/menuTreeoption [get]
// @Security Bearer
func GetMenuTreeOption(c *gin.Context) {
	items, err := models.CMenu.QueryLabelTree(gcontext.Context(c))
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError,
			servers.WithPrompt(prompt.QueryFailed),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(items))
}

// @Summary 根据角色名称获取菜单列表数据（左菜单使用）
// @Description 获取JSON
// @Tags 菜单
// @Param id path int true "id"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menurole [get]
// @Security Bearer
func GetMenuTreeWithRoleName(c *gin.Context) {
	items, err := models.CMenu.QueryTreeWithRoleName(gcontext.Context(c))
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(items))
}

// @Summary 获取角色对应的菜单id数组
// @Description 获取JSON
// @Tags 菜单
// @Param id path int true "id"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menuids/{id} [get]
// @Security Bearer
func GetMenuIDS(c *gin.Context) {
	var data models.RoleMenu

	data.RoleName = jwtauth.FromRoleName(gcontext.Context(c))
	data.Updator = jwtauth.FromUserIdStr(gcontext.Context(c))
	result, err := data.GetIDS()
	if err != nil {
		servers.Fail(c, http.StatusNotFound,
			servers.WithPrompt(prompt.NotFound),
			servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(result))
}
