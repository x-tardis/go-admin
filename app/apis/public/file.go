package public

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/go-core-package/extos"

	"github.com/x-tardis/go-admin/pkg/extimg"
	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
)

type FileResponse struct {
	Size     int64  `json:"size"`
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

// @Summary 上传图片
// @Description 获取JSON
// @Tags 公共接口
// @Accept multipart/form-data
// @Param type query string true "type" (1：单图，2：多图, 3：base64图片)
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/public/uploadFile [post]
func UploadFile(c *gin.Context) {
	var fileResponse FileResponse

	tag, _ := c.GetPostForm("type")
	urlPrefix := fmt.Sprintf("http://%s/", c.Request.Host)

	switch tag {
	case "":
		servers.Fail(c, http.StatusOK, servers.WithMsg("缺少标识"))
	case "1": // 单图
		files, err := c.FormFile("file")
		if err != nil {
			servers.Fail(c, http.StatusOK, servers.WithMsg("图片不能为空"))
			return
		}
		// 上传文件至指定目录
		guid := infra.GenerateUUID()
		singleFile := "static/uploadfile/" + guid + path.Ext(files.Filename)
		_ = c.SaveUploadedFile(files, singleFile)
		fileType, _ := extimg.GetType(singleFile)
		fileSize, _ := extos.FileSize(singleFile)
		fileResponse = FileResponse{
			Size:     fileSize,
			Path:     singleFile,
			FullPath: urlPrefix + singleFile,
			Name:     files.Filename,
			Type:     fileType,
		}
		servers.OK(c, servers.WithData(fileResponse), servers.WithMsg("上传成功"))
		return
	case "2": // 多图
		files := c.Request.MultipartForm.File["file"]
		var multipartFile []FileResponse
		for _, f := range files {
			guid := infra.GenerateUUID()
			multipartFileName := "static/uploadfile/" + guid + path.Ext(f.Filename)
			e := c.SaveUploadedFile(f, multipartFileName)
			fileType, _ := extimg.GetType(multipartFileName)
			if e == nil {
				fileSize, _ := extos.FileSize(multipartFileName)
				multipartFile = append(multipartFile, FileResponse{
					Size:     fileSize,
					Path:     multipartFileName,
					FullPath: urlPrefix + multipartFileName,
					Name:     f.Filename,
					Type:     fileType,
				})
			}
		}

		servers.OK(c, servers.WithData(multipartFile), servers.WithMsg("上传成功"))
		return
	case "3": // base64
		files, _ := c.GetPostForm("file")
		typeStr, ddd, _ := extimg.DecodeBase64(files)
		guid := infra.GenerateUUID()
		base64File := "static/uploadfile/" + guid + ".jpg"
		_ = ioutil.WriteFile(base64File, ddd, 0666)
		fileSize, _ := extos.FileSize(base64File)
		fileResponse = FileResponse{
			Size:     fileSize,
			Path:     base64File,
			FullPath: urlPrefix + base64File,
			Name:     "",
			Type:     typeStr,
		}
		servers.OK(c, servers.WithData(fileResponse), servers.WithMsg("上传成功"))
	}
}
