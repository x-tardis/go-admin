package syscontent

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

type Content struct{}

func (Content) QueryPage(c *gin.Context) {
	qp := models.ContentQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	qp.Inspect()

	items, info, err := models.CContent.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(paginator.Pages{
		Info: info,
		List: items,
	}))
}

func (Content) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	result, err := models.CContent.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, result, "")
}

// @Summary 添加内容管理
// @Description 获取JSON
// @Tags 内容管理
// @Accept  application/json
// @Product application/json
// @Param data body models.Content true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/syscontents [post]
func (Content) Create(c *gin.Context) {
	newItem := models.Content{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}

	item, err := models.CContent.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, item, "")
}

func (Content) Update(c *gin.Context) {
	up := models.Content{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}

	item, err := models.CContent.Update(gcontext.Context(c), up.Id, up)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, item, "")
}

func (Content) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := models.CContent.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.OKWithRequestID(c, nil, codes.DeletedSuccess)
}
