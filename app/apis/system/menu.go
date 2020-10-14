package system

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/jwtauth"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// @Summary Menu列表数据
// @Description 获取JSON
// @Tags 菜单
// @Param menuName query string false "menuName"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menulist [get]
// @Security Bearer
func GetMenuList(c *gin.Context) {
	var Menu models.Menu
	Menu.MenuName = c.Request.FormValue("menuName")
	Menu.Visible = c.Request.FormValue("visible")
	Menu.Title = c.Request.FormValue("title")
	Menu.DataScope = jwtauth.UserIdStr(c)
	result, err := Menu.SetMenu()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, result, "")
}

// @Summary Menu列表数据
// @Description 获取JSON
// @Tags 菜单
// @Param menuName query string false "menuName"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menu [get]
// @Security Bearer
func GetMenu(c *gin.Context) {
	var data models.Menu
	id, err := strconv.Atoi(c.Param("id"))
	data.MenuId = id
	result, err := data.GetByMenuId()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, result, "")
}

func GetMenuTreeRoleselect(c *gin.Context) {
	var Menu models.Menu
	var SysRole models.SysRole
	id, err := strconv.Atoi(c.Param("roleId"))
	SysRole.RoleId = id
	result, err := Menu.SetMenuLable()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	menuIds := make([]int, 0)
	if id != 0 {
		menuIds, err = SysRole.GetRoleMeunId()
		if err != nil {
			servers.Fail(c, -1, codes.NotFoundRelatedInfo)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
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
// @Router /api/v1/menuTreeselect [get]
// @Security Bearer
func GetMenuTreeelect(c *gin.Context) {
	var data models.Menu
	result, err := data.SetMenuLable()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, result, "")
}

// @Summary 创建菜单
// @Description 获取JSON
// @Tags 菜单
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
// @Router /api/v1/menu [post]
// @Security Bearer
func InsertMenu(c *gin.Context) {
	var data models.Menu

	if err := c.ShouldBindJSON(&data); err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	data.CreateBy = jwtauth.UserIdStr(c)
	result, err := data.Create()
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, result, "")
}

// @Summary 修改菜单
// @Description 获取JSON
// @Tags 菜单
// @Accept  application/x-www-form-urlencoded
// @Product application/x-www-form-urlencoded
// @Param id path int true "id"
// @Param data body models.Menu true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/menu/{id} [put]
// @Security Bearer
func UpdateMenu(c *gin.Context) {
	var data models.Menu

	if err := c.ShouldBindJSON(&data); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	data.UpdateBy = jwtauth.UserIdStr(c)
	_, err := data.Update(data.MenuId)
	if err != nil {
		servers.Fail(c, 501, err.Error())
		return
	}
	servers.OKWithRequestID(c, "", "修改成功")

}

// @Summary 删除菜单
// @Description 删除数据
// @Tags 菜单
// @Param id path int true "id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/menu/{id} [delete]
func DeleteMenu(c *gin.Context) {
	var data models.Menu
	id, err := strconv.Atoi(c.Param("id"))
	data.UpdateBy = jwtauth.UserIdStr(c)
	_, err = data.Delete(id)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.OKWithRequestID(c, "", "删除成功")
}

// @Summary 根据角色名称获取菜单列表数据（左菜单使用）
// @Description 获取JSON
// @Tags 菜单
// @Param id path int true "id"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menurole [get]
// @Security Bearer
func GetMenuRole(c *gin.Context) {
	var Menu models.Menu
	result, err := Menu.SetMenuRole(jwtauth.RoleKey(c))
	if err != nil {
		servers.Fail(c, 500, codes.GetFail)
		return
	}
	servers.OKWithRequestID(c, result, "")
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
	data.RoleName = jwtauth.RoleName(c)
	data.UpdateBy = jwtauth.UserIdStr(c)
	result, err := data.GetIDS()
	if err != nil {
		servers.Fail(c, 500, codes.GetFail)
		return
	}
	servers.OKWithRequestID(c, result, codes.GetSuccess)
}
