package sysfiledir

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/codes"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
)

type FileDir struct{}

func (FileDir) QueryTree(c *gin.Context) {
	qp := models.FileDirQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}

	items, err := models.CFileDir.QueryTree(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, items, "")
}

func (FileDir) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CFileDir.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, -1, codes.NotFoundRelatedInfo)
		return
	}
	servers.OKWithRequestID(c, item, "")
}

// @Summary 添加SysFileDir
// @Description 获取JSON
// @Tags SysFileDir
// @Accept  application/json
// @Product application/json
// @Param data body models.FileDir true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/sysfiledir [post]
func (FileDir) Create(c *gin.Context) {
	newItem := models.FileDir{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}

	item, err := models.CFileDir.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, item, "")
}

func (FileDir) Update(c *gin.Context) {
	up := models.FileDir{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}

	item, err := models.CFileDir.Update(gcontext.Context(c), up.Id, up)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, item, "")
}

func (FileDir) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	if err := models.CFileDir.BatchDelete(gcontext.Context(c), ids); err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.OKWithRequestID(c, nil, codes.DeletedSuccess)
}
