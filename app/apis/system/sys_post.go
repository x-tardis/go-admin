package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/thinkgos/sharp/core/paginator"
	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
)

type Post struct{}

// @Summary 岗位列表数据
// @Description 获取JSON
// @Tags 岗位
// @Param postName query string false "postName"
// @Param postCode query string false "postCode"
// @Param postId query string false "postId"
// @Param status query string false "status"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/posts [get]
// @Security Bearer
func (Post) QueryPage(c *gin.Context) {
	qp := models.PostQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, codes.DataParseFailed)
		return
	}
	qp.Inspect()

	items, info, err := new(models.CallPost).QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(&paginator.Pages{
		Info: info,
		List: items,
	}))
}

// @Summary 获取岗位信息
// @Description 获取JSON
// @Tags 岗位
// @Param id path int true "post id"
// @Success 200 {object} servers.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/posts/{id} [get]
// @Security Bearer
func (Post) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := new(models.CallPost).Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, item, "")
}

// @Summary 添加岗位
// @Description 获取JSON
// @Tags 岗位
// @Accept  application/json
// @Product application/json
// @Param data body models.Post true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/posts [post]
// @Security Bearer
func (Post) Create(c *gin.Context) {
	var item models.Post

	err := c.ShouldBindJSON(&item)
	if err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}

	result, err := new(models.CallPost).Create(gcontext.Context(c), item)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, result, "")
}

// @Summary 修改岗位
// @Description 获取JSON
// @Tags 岗位
// @Accept  application/json
// @Product application/json
// @Param data body models.Post true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/posts [put]
// @Security Bearer
func (Post) Update(c *gin.Context) {
	var up models.Post

	err := c.ShouldBindJSON(&up)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	item, err := new(models.CallPost).Update(gcontext.Context(c), up.PostId, up)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, item, codes.UpdatedSuccess)
}

// @Summary 删除岗位
// @Description 删除数据
// @Tags 岗位
// @Param id path int true "id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 500 {string} string	"{"code": 500, "message": "删除失败"}"
// @Router /api/v1/posts/{ids} [delete]
func (Post) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := new(models.CallPost).BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithMsg(codes.DeletedSuccess))
}
