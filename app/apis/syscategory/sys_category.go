package syscategory

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

type Category struct{}

func (Category) QueryPage(c *gin.Context) {
	qp := models.CategoryQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	qp.Inspect()

	items, info, err := models.CCategory.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(paginator.Pages{
		Info: info,
		List: items,
	}))
}

func (Category) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CCategory.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, item, "")
}

// @Summary 添加分类
// @Description 获取JSON
// @Tags 分类
// @Accept  application/json
// @Product application/json
// @Param data body models.Category true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/syscategory [post]
func (Category) Create(c *gin.Context) {
	newItem := models.Category{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}

	item, err := models.CCategory.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, item, "")
}

func (Category) Update(c *gin.Context) {
	up := models.Category{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}

	result, err := models.CCategory.Update(gcontext.Context(c), up.Id, up)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, result, "")
}

func (Category) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("id"))
	err := models.CCategory.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.OKWithRequestID(c, nil, codes.DeletedSuccess)
}
