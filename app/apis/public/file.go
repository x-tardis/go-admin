package public

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/thinkgos/go-core-package/extimg"
	"github.com/thinkgos/go-core-package/extos"

	"github.com/x-tardis/go-admin/pkg/infra"
	"github.com/x-tardis/go-admin/pkg/servers"
)

// FileResponse
type FileResponse struct {
	Size     int64  `json:"size"`
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

// @tags 文件上传
// @summary 上传图片
// @description 上传图片
// @accept multipart/form-data
// @produce json
// @param type formData string true "type" (1：单图，2：多图, 3：base64图片)
// @success 200 {object} servers.Response "成功"
// @failure 400 {object} servers.Response "错误请求"
// @failure 401 {object} servers.Response "鉴权失败"
// @failure 500 {object} servers.Response "服务器内部错误"
// @router /api/v1/public/uploadFile [post]
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
