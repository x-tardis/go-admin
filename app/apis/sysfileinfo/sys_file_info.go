package sysfileinfo

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

type FileInfo struct{}

func (FileInfo) QueryPage(c *gin.Context) {
	qp := models.FileInfoQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}
	qp.Inspect()

	items, info, err := models.CFileInfo.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.JSON(c, http.StatusOK, servers.WithData(paginator.Pages{
		Info: info,
		List: items,
	}))
}

func (FileInfo) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	result, err := models.CFileInfo.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, -1, "抱歉未找到相关信息")
		return
	}
	servers.OKWithRequestID(c, result, "")
}

func (FileInfo) Create(c *gin.Context) {
	newItem := models.FileInfo{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, 500, err.Error())
		return
	}

	item, err := models.CFileInfo.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, item, "")
}

func (FileInfo) Update(c *gin.Context) {
	up := models.FileInfo{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, -1, codes.DataParseFailed)
		return
	}

	result, err := models.CFileInfo.Update(gcontext.Context(c), up.Id, up)
	if err != nil {
		servers.Fail(c, -1, err.Error())
		return
	}
	servers.OKWithRequestID(c, result, "")
}

func (FileInfo) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := models.CFileInfo.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, 500, codes.DeletedFail)
		return
	}
	servers.OKWithRequestID(c, nil, codes.DeletedSuccess)
}
