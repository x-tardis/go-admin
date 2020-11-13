package migrate

import (
	"runtime"

	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	Register(GetFilename(fileName), _1599190683680Test)
}

func _1599190683680Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		list := []models.Menu{
			{MenuId: 496, MenuName: "Sources", Title: "资源管理", Icon: "network", Path: "/sources", Paths: "/0/496", MenuType: models.MenuTypeToc, Method: "无", Permission: "", ParentId: 0, NoCache: true, Breadcrumb: "", Component: "Layout", Sort: 3, Visible: "0", Creator: "1", Updator: "1", IsFrame: "1"},
			{MenuId: 497, MenuName: "File", Title: "文件管理", Icon: "documentation", Path: "file-manage", Paths: "/0/496/497", MenuType: models.MenuTypeMenu, Method: "", Permission: "", ParentId: 496, NoCache: true, Breadcrumb: "", Component: "/fileManage/index", Sort: 0, Visible: "0", Creator: "1", Updator: "1", IsFrame: "1"},
			{MenuId: 498, MenuName: "", Title: "内容管理", Icon: "pass", Path: "/content", Paths: "/0/498", MenuType: models.MenuTypeToc, Method: "无", Permission: "", ParentId: 0, NoCache: true, Breadcrumb: "", Component: "Layout", Sort: 4, Visible: "0", Creator: "1", Updator: "1", IsFrame: "1"},
			{MenuId: 499, MenuName: "Category", Title: "分类", Icon: "pass", Path: "syscategory", Paths: "/0/498/499", MenuType: models.MenuTypeMenu, Method: "无", Permission: "syscategory:syscategory:list", ParentId: 498, NoCache: true, Breadcrumb: "", Component: "/syscategory/index", Sort: 0, Visible: "0", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 500, MenuName: "", Title: "分页获取分类", Icon: "pass", Path: "", Paths: "/0/498/499/500", MenuType: models.MenuTypeBtn, Method: "无", Permission: "syscategory:syscategory:query", ParentId: 499, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 501, MenuName: "", Title: "创建分类", Icon: "pass", Path: "", Paths: "/0/498/499/501", MenuType: models.MenuTypeBtn, Method: "无", Permission: "syscategory:syscategory:add", ParentId: 499, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 502, MenuName: "", Title: "修改分类", Icon: "pass", Path: "", Paths: "/0/498/499/502", MenuType: models.MenuTypeBtn, Method: "无", Permission: "syscategory:syscategory:edit", ParentId: 499, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 503, MenuName: "", Title: "删除分类", Icon: "pass", Path: "", Paths: "/0/498/499/503", MenuType: models.MenuTypeBtn, Method: "无", Permission: "syscategory:syscategory:remove", ParentId: 499, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 504, MenuName: "Category", Title: "分类", Icon: "bug", Path: "category", Paths: "/0/63/504", MenuType: models.MenuTypeToc, Method: "无", Permission: "", ParentId: 63, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 505, MenuName: "", Title: "分页获取分类", Icon: "bug", Path: "/api/v1/categories", Paths: "/0/63/504/505", MenuType: models.MenuTypeIfc, Method: "GET", Permission: "", ParentId: 504, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 506, MenuName: "", Title: "根据id获取分类", Icon: "bug", Path: "/api/v1/categories/:id", Paths: "/0/63/504/506", MenuType: models.MenuTypeIfc, Method: "GET", Permission: "", ParentId: 504, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 507, MenuName: "", Title: "创建分类", Icon: "bug", Path: "/api/v1/categories", Paths: "/0/63/504/507", MenuType: models.MenuTypeIfc, Method: "POST", Permission: "", ParentId: 504, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 508, MenuName: "", Title: "修改分类", Icon: "bug", Path: "/api/v1/categories", Paths: "/0/63/504/508", MenuType: models.MenuTypeIfc, Method: "PUT", Permission: "", ParentId: 504, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 509, MenuName: "", Title: "删除分类", Icon: "bug", Path: "/api/v1/categories/:id", Paths: "/0/63/504/509", MenuType: models.MenuTypeIfc, Method: "DELETE", Permission: "", ParentId: 504, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 511, MenuName: "Content", Title: "内容管理", Icon: "pass", Path: "syscontent", Paths: "/0/498/511", MenuType: models.MenuTypeMenu, Method: "无", Permission: "syscontent:syscontent:list", ParentId: 498, NoCache: true, Breadcrumb: "", Component: "/syscontent/index", Sort: 0, Visible: "0", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 512, MenuName: "", Title: "分页获取内容管理", Icon: "pass", Path: "", Paths: "/0/510/511/512", MenuType: models.MenuTypeBtn, Method: "无", Permission: "syscontent:syscontent:query", ParentId: 511, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 513, MenuName: "", Title: "创建内容管理", Icon: "pass", Path: "", Paths: "/0/510/511/513", MenuType: models.MenuTypeBtn, Method: "无", Permission: "syscontent:syscontent:add", ParentId: 511, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 514, MenuName: "", Title: "修改内容管理", Icon: "pass", Path: "", Paths: "/0/510/511/514", MenuType: models.MenuTypeBtn, Method: "无", Permission: "syscontent:syscontent:edit", ParentId: 511, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 515, MenuName: "", Title: "删除内容管理", Icon: "pass", Path: "", Paths: "/0/510/511/515", MenuType: models.MenuTypeBtn, Method: "无", Permission: "syscontent:syscontent:remove", ParentId: 511, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "0", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 516, MenuName: "Content", Title: "内容管理", Icon: "bug", Path: "content", Paths: "/0/63/516", MenuType: models.MenuTypeToc, Method: "无", Permission: "", ParentId: 63, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 517, MenuName: "", Title: "分页获取内容管理", Icon: "bug", Path: "/api/v1/contents", Paths: "/0/63/516/517", MenuType: models.MenuTypeIfc, Method: "GET", Permission: "", ParentId: 516, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 518, MenuName: "", Title: "根据id获取内容管理", Icon: "bug", Path: "/api/v1/contents/:id", Paths: "/0/63/516/518", MenuType: models.MenuTypeIfc, Method: "GET", Permission: "", ParentId: 516, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 519, MenuName: "", Title: "创建内容管理", Icon: "bug", Path: "/api/v1/contents", Paths: "/0/63/516/519", MenuType: models.MenuTypeIfc, Method: "POST", Permission: "", ParentId: 516, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 520, MenuName: "", Title: "修改内容管理", Icon: "bug", Path: "/api/v1/contents", Paths: "/0/63/516/520", MenuType: models.MenuTypeIfc, Method: "PUT", Permission: "", ParentId: 516, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", Creator: "1", Updator: "1", IsFrame: "0"},
			{MenuId: 521, MenuName: "", Title: "删除内容管理", Icon: "bug", Path: "/api/v1/contents/:id", Paths: "/0/63/516/521", MenuType: models.MenuTypeIfc, Method: "DELETE", Permission: "", ParentId: 516, NoCache: true, Breadcrumb: "", Component: "", Sort: 0, Visible: "1", Creator: "1", Updator: "1", IsFrame: "0"},
		}

		err := tx.Create(list).Error
		if err != nil {
			return err
		}
		return tx.Create(&models.Migration{Version: version}).Error
	})
}
