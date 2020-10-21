package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thinkgos/sharp/gin/gcontext"

	"github.com/thinkgos/sharp/core/paginator"

	"github.com/x-tardis/go-admin/app/models"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
	"github.com/x-tardis/go-admin/pkg/servers/prompt"
)

type FileInfo struct{}

func (FileInfo) QueryPage(c *gin.Context) {
	qp := models.FileInfoQueryParam{}
	if err := c.ShouldBindQuery(&qp); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}
	qp.Inspect()

	items, info, err := models.CFileInfo.QueryPage(gcontext.Context(c), qp)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(paginator.Pages{
		Info: info,
		List: items,
	}))
}

func (FileInfo) Get(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	item, err := models.CFileInfo.Get(gcontext.Context(c), id)
	if err != nil {
		servers.Fail(c, http.StatusNotFound, servers.WithPrompt(prompt.NotFound))
		return
	}
	servers.OK(c, servers.WithData(item))
}

func (FileInfo) Create(c *gin.Context) {
	newItem := models.FileInfo{}
	if err := c.ShouldBindJSON(&newItem); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CFileInfo.Create(gcontext.Context(c), newItem)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithError(err))
		return
	}
	servers.OK(c, servers.WithData(item))
}

func (FileInfo) Update(c *gin.Context) {
	up := models.FileInfo{}
	if err := c.ShouldBindJSON(&up); err != nil {
		servers.Fail(c, http.StatusBadRequest, servers.WithError(err))
		return
	}

	item, err := models.CFileInfo.Update(gcontext.Context(c), up.Id, up)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.UpdateFailed))
		return
	}
	servers.OK(c, servers.WithData(item))
}

func (FileInfo) BatchDelete(c *gin.Context) {
	ids := infra.ParseIdsGroup(c.Param("ids"))
	err := models.CFileInfo.BatchDelete(gcontext.Context(c), ids)
	if err != nil {
		servers.Fail(c, http.StatusInternalServerError, servers.WithPrompt(prompt.DeleteFailed))
		return
	}
	servers.OK(c, servers.WithPrompt(prompt.DeleteSuccess))
}
