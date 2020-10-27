package version

import (
	"runtime"

	"gorm.io/gorm"

	"github.com/x-tardis/go-admin/cmd/migrate/migration"
	"github.com/x-tardis/go-admin/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1599190683670Test)
}

func _1599190683670Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		list1 := []models.RoleMenu{
			{RoleId: 1, MenuId: 2, RoleName: "admin"},
			{RoleId: 1, MenuId: 3, RoleName: "admin"},
			{RoleId: 1, MenuId: 43, RoleName: "admin"},
			{RoleId: 1, MenuId: 44, RoleName: "admin"},
			{RoleId: 1, MenuId: 45, RoleName: "admin"},
			{RoleId: 1, MenuId: 46, RoleName: "admin"},
			{RoleId: 1, MenuId: 51, RoleName: "admin"},
			{RoleId: 1, MenuId: 52, RoleName: "admin"},
			{RoleId: 1, MenuId: 56, RoleName: "admin"},
			{RoleId: 1, MenuId: 57, RoleName: "admin"},
			{RoleId: 1, MenuId: 58, RoleName: "admin"},
			{RoleId: 1, MenuId: 59, RoleName: "admin"},
			{RoleId: 1, MenuId: 60, RoleName: "admin"},
			{RoleId: 1, MenuId: 61, RoleName: "admin"},
			{RoleId: 1, MenuId: 62, RoleName: "admin"},
			{RoleId: 1, MenuId: 63, RoleName: "admin"},
			{RoleId: 1, MenuId: 64, RoleName: "admin"},
			{RoleId: 1, MenuId: 66, RoleName: "admin"},
			{RoleId: 1, MenuId: 67, RoleName: "admin"},
			{RoleId: 1, MenuId: 68, RoleName: "admin"},
			{RoleId: 1, MenuId: 69, RoleName: "admin"},
			{RoleId: 1, MenuId: 70, RoleName: "admin"},
			{RoleId: 1, MenuId: 71, RoleName: "admin"},
			{RoleId: 1, MenuId: 72, RoleName: "admin"},
			{RoleId: 1, MenuId: 73, RoleName: "admin"},
			{RoleId: 1, MenuId: 74, RoleName: "admin"},
			{RoleId: 1, MenuId: 75, RoleName: "admin"},
			{RoleId: 1, MenuId: 76, RoleName: "admin"},
			{RoleId: 1, MenuId: 77, RoleName: "admin"},
			{RoleId: 1, MenuId: 78, RoleName: "admin"},
			{RoleId: 1, MenuId: 79, RoleName: "admin"},
			{RoleId: 1, MenuId: 80, RoleName: "admin"},
			{RoleId: 1, MenuId: 81, RoleName: "admin"},
			{RoleId: 1, MenuId: 82, RoleName: "admin"},
			{RoleId: 1, MenuId: 83, RoleName: "admin"},
			{RoleId: 1, MenuId: 84, RoleName: "admin"},
			{RoleId: 1, MenuId: 85, RoleName: "admin"},
			{RoleId: 1, MenuId: 86, RoleName: "admin"},
			{RoleId: 1, MenuId: 87, RoleName: "admin"},
			{RoleId: 1, MenuId: 89, RoleName: "admin"},
			{RoleId: 1, MenuId: 90, RoleName: "admin"},
			{RoleId: 1, MenuId: 91, RoleName: "admin"},
			{RoleId: 1, MenuId: 92, RoleName: "admin"},
			{RoleId: 1, MenuId: 93, RoleName: "admin"},
			{RoleId: 1, MenuId: 94, RoleName: "admin"},
			{RoleId: 1, MenuId: 95, RoleName: "admin"},
			{RoleId: 1, MenuId: 96, RoleName: "admin"},
			{RoleId: 1, MenuId: 97, RoleName: "admin"},
			{RoleId: 1, MenuId: 103, RoleName: "admin"},
			{RoleId: 1, MenuId: 104, RoleName: "admin"},
			{RoleId: 1, MenuId: 105, RoleName: "admin"},
			{RoleId: 1, MenuId: 106, RoleName: "admin"},
			{RoleId: 1, MenuId: 107, RoleName: "admin"},
			{RoleId: 1, MenuId: 108, RoleName: "admin"},
			{RoleId: 1, MenuId: 109, RoleName: "admin"},
			{RoleId: 1, MenuId: 110, RoleName: "admin"},
			{RoleId: 1, MenuId: 111, RoleName: "admin"},
			{RoleId: 1, MenuId: 112, RoleName: "admin"},
			{RoleId: 1, MenuId: 113, RoleName: "admin"},
			{RoleId: 1, MenuId: 114, RoleName: "admin"},
			{RoleId: 1, MenuId: 115, RoleName: "admin"},
			{RoleId: 1, MenuId: 116, RoleName: "admin"},
			{RoleId: 1, MenuId: 117, RoleName: "admin"},
			{RoleId: 1, MenuId: 118, RoleName: "admin"},
			{RoleId: 1, MenuId: 119, RoleName: "admin"},
			{RoleId: 1, MenuId: 120, RoleName: "admin"},
			{RoleId: 1, MenuId: 121, RoleName: "admin"},
			{RoleId: 1, MenuId: 122, RoleName: "admin"},
			{RoleId: 1, MenuId: 123, RoleName: "admin"},
			{RoleId: 1, MenuId: 138, RoleName: "admin"},
			{RoleId: 1, MenuId: 142, RoleName: "admin"},
			{RoleId: 1, MenuId: 201, RoleName: "admin"},
			{RoleId: 1, MenuId: 202, RoleName: "admin"},
			{RoleId: 1, MenuId: 203, RoleName: "admin"},
			{RoleId: 1, MenuId: 204, RoleName: "admin"},
			{RoleId: 1, MenuId: 205, RoleName: "admin"},
			{RoleId: 1, MenuId: 206, RoleName: "admin"},
			{RoleId: 1, MenuId: 211, RoleName: "admin"},
			{RoleId: 1, MenuId: 212, RoleName: "admin"},
			{RoleId: 1, MenuId: 213, RoleName: "admin"},
			{RoleId: 1, MenuId: 214, RoleName: "admin"},
			{RoleId: 1, MenuId: 215, RoleName: "admin"},
			{RoleId: 1, MenuId: 216, RoleName: "admin"},
			{RoleId: 1, MenuId: 217, RoleName: "admin"},
			{RoleId: 1, MenuId: 220, RoleName: "admin"},
			{RoleId: 1, MenuId: 221, RoleName: "admin"},
			{RoleId: 1, MenuId: 222, RoleName: "admin"},
			{RoleId: 1, MenuId: 223, RoleName: "admin"},
			{RoleId: 1, MenuId: 224, RoleName: "admin"},
			{RoleId: 1, MenuId: 225, RoleName: "admin"},
			{RoleId: 1, MenuId: 226, RoleName: "admin"},
			{RoleId: 1, MenuId: 227, RoleName: "admin"},
			{RoleId: 1, MenuId: 228, RoleName: "admin"},
			{RoleId: 1, MenuId: 229, RoleName: "admin"},
			{RoleId: 1, MenuId: 230, RoleName: "admin"},
			{RoleId: 1, MenuId: 231, RoleName: "admin"},
			{RoleId: 1, MenuId: 232, RoleName: "admin"},
			{RoleId: 1, MenuId: 233, RoleName: "admin"},
			{RoleId: 1, MenuId: 234, RoleName: "admin"},
			{RoleId: 1, MenuId: 235, RoleName: "admin"},
			{RoleId: 1, MenuId: 236, RoleName: "admin"},
			{RoleId: 1, MenuId: 237, RoleName: "admin"},
			{RoleId: 1, MenuId: 238, RoleName: "admin"},
			{RoleId: 1, MenuId: 239, RoleName: "admin"},
			{RoleId: 1, MenuId: 240, RoleName: "admin"},
			{RoleId: 1, MenuId: 241, RoleName: "admin"},
			{RoleId: 1, MenuId: 242, RoleName: "admin"},
			{RoleId: 1, MenuId: 243, RoleName: "admin"},
			{RoleId: 1, MenuId: 244, RoleName: "admin"},
			{RoleId: 1, MenuId: 245, RoleName: "admin"},
			{RoleId: 1, MenuId: 246, RoleName: "admin"},
			{RoleId: 1, MenuId: 247, RoleName: "admin"},
			{RoleId: 1, MenuId: 248, RoleName: "admin"},
			{RoleId: 1, MenuId: 249, RoleName: "admin"},
			{RoleId: 1, MenuId: 250, RoleName: "admin"},
			{RoleId: 1, MenuId: 251, RoleName: "admin"},
			{RoleId: 1, MenuId: 252, RoleName: "admin"},
			{RoleId: 1, MenuId: 253, RoleName: "admin"},
			{RoleId: 1, MenuId: 254, RoleName: "admin"},
			{RoleId: 1, MenuId: 255, RoleName: "admin"},
			{RoleId: 1, MenuId: 256, RoleName: "admin"},
			{RoleId: 1, MenuId: 257, RoleName: "admin"},
			{RoleId: 1, MenuId: 258, RoleName: "admin"},
			{RoleId: 1, MenuId: 259, RoleName: "admin"},
			{RoleId: 1, MenuId: 260, RoleName: "admin"},
			{RoleId: 1, MenuId: 261, RoleName: "admin"},
			{RoleId: 1, MenuId: 262, RoleName: "admin"},
			{RoleId: 1, MenuId: 263, RoleName: "admin"},
			{RoleId: 1, MenuId: 264, RoleName: "admin"},
			{RoleId: 1, MenuId: 267, RoleName: "admin"},
			{RoleId: 1, MenuId: 269, RoleName: "admin"},
			{RoleId: 1, MenuId: 459, RoleName: "admin"},
			{RoleId: 1, MenuId: 460, RoleName: "admin"},
			{RoleId: 1, MenuId: 461, RoleName: "admin"},
			{RoleId: 1, MenuId: 462, RoleName: "admin"},
			{RoleId: 1, MenuId: 463, RoleName: "admin"},
			{RoleId: 1, MenuId: 464, RoleName: "admin"},
			{RoleId: 1, MenuId: 465, RoleName: "admin"},
			{RoleId: 1, MenuId: 466, RoleName: "admin"},
			{RoleId: 1, MenuId: 467, RoleName: "admin"},
			{RoleId: 1, MenuId: 468, RoleName: "admin"},
			{RoleId: 1, MenuId: 469, RoleName: "admin"},
			{RoleId: 1, MenuId: 470, RoleName: "admin"},
			{RoleId: 1, MenuId: 471, RoleName: "admin"},
			{RoleId: 1, MenuId: 473, RoleName: "admin"},
			{RoleId: 1, MenuId: 474, RoleName: "admin"},
			{RoleId: 1, MenuId: 475, RoleName: "admin"},
			{RoleId: 1, MenuId: 476, RoleName: "admin"},
			{RoleId: 1, MenuId: 477, RoleName: "admin"},
			{RoleId: 1, MenuId: 478, RoleName: "admin"},
			{RoleId: 1, MenuId: 479, RoleName: "admin"},
			{RoleId: 1, MenuId: 480, RoleName: "admin"},
			{RoleId: 1, MenuId: 481, RoleName: "admin"},
			{RoleId: 1, MenuId: 482, RoleName: "admin"},
			{RoleId: 1, MenuId: 483, RoleName: "admin"},
		}
		list2 := []models.CasbinRule{
			{PType: "p", V0: "admin", V1: "/api/v1/menus", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/menus/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/menus", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/menus", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/menus/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/users", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/users/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/users/", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/users", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/users", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/users/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/user/profile", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/roles", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/roles/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/roles", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/roles", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/roles/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/configs", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/configs/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/configs", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/configs", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/configs/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/menurole", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/roleMenuTreeoption/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/menuTreeoption", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/rolemenu", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/rolemenu", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/rolemenu/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/depts", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/depts/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/depts", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/depts", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/depts/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/data", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/data/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/data", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/data", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/data/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/databytype/", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/databytype/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/type", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/type/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/type", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/type", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/type/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/dict/typeoption", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/posts", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/posts/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/posts", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/posts", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/posts/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/menuids", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/loginlog", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/loginlog/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/operlog", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/operlog/:id", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/getinfo", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/roledatascope", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/roleDeptTreeoption/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/deptTree", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/configKey/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/logout", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/user/avatar", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/user/pwd", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/jobs", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/jobs/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/jobs", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/jobs", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/jobs", V2: "DELETE"},
			{PType: "p", V0: "admin", V1: "/api/v1/system/setting", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/system/setting/:id", V2: "GET"},
			{PType: "p", V0: "admin", V1: "/api/v1/system/setting", V2: "POST"},
			{PType: "p", V0: "admin", V1: "/api/v1/system/setting", V2: "PUT"},
			{PType: "p", V0: "admin", V1: "/api/v1/system/setting/:id", V2: "DELETE"},
		}

		list3 := []models.Dept{
			{DeptId: 1, ParentId: 0, DeptPath: "/0/1", DeptName: "爱拓科技", Sort: 0, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "0", Creator: "1", Updator: "1"},
			{DeptId: 7, ParentId: 1, DeptPath: "/0/1/7", DeptName: "研发部", Sort: 1, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "0", Creator: "1", Updator: "1"},
			{DeptId: 8, ParentId: 1, DeptPath: "/0/1/8", DeptName: "运维部", Sort: 0, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "0", Creator: "1", Updator: "1"},
			{DeptId: 9, ParentId: 1, DeptPath: "/0/1/9", DeptName: "客服部", Sort: 0, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "0", Creator: "1", Updator: "1"},
			{DeptId: 10, ParentId: 1, DeptPath: "/0/1/10", DeptName: "人力资源", Sort: 3, Leader: "aituo", Phone: "13782218188", Email: "atuo@aituo.com", Status: "1", Creator: "1", Updator: "1"},
		}

		list4 := []models.Config{
			{ConfigId: 1, ConfigName: "主框架页-默认皮肤样式名称", ConfigKey: "sys_index_skinName", ConfigValue: "skin-blue", ConfigType: "Y", Remark: "蓝色 skin-blue、绿色 skin-green、紫色 skin-purple、红色 skin-red、黄色 skin-yellow", Creator: "1", Updator: "1"},
			{ConfigId: 2, ConfigName: "用户管理-账号初始密码", ConfigKey: "sys.user.initPassword", ConfigValue: "123456", ConfigType: "Y", Remark: "初始化密码 123456", Creator: "1", Updator: "1"},
			{ConfigId: 3, ConfigName: "主框架页-侧边栏主题", ConfigKey: "sys_index_sideTheme", ConfigValue: "theme-dark", ConfigType: "Y", Remark: "深色主题theme-dark，浅色主题theme-light", Creator: "1", Updator: "1"},
		}

		list5 := []models.Post{
			{PostId: 1, PostCode: "CEO", PostName: "首席执行官", Sort: 1, Status: "0", Remark: "首席执行官", Creator: "1", Updator: "1"},
			{PostId: 2, PostCode: "CTO", PostName: "首席技术执行官", Sort: 2, Status: "0", Remark: "首席技术执行官", Creator: "1", Updator: "1"},
			{PostId: 3, PostCode: "COO", PostName: "首席运营官", Sort: 3, Status: "0", Remark: "测试工程师", Creator: "1", Updator: "1"},
		}

		list6 := []models.Role{
			{1, "系统管理员", "admin", "0", 1, "", true, "", "", "1", "", models.Model{}, []int{}, []int{}, ""},
		}

		list7 := []models.DictType{
			{DictId: 1, DictName: "系统开关", DictType: "sys_normal_disable", Status: "0", Creator: "1", Updator: "1", Remark: "系统开关列表"},
			{DictId: 2, DictName: "用户性别", DictType: "sys_user_sex", Status: "0", Creator: "1", Updator: "", Remark: "用户性别列表"},
			{DictId: 3, DictName: "菜单状态", DictType: "sys_show_hide", Status: "0", Creator: "1", Updator: "", Remark: "菜单状态列表"},
			{DictId: 4, DictName: "系统是否", DictType: "sys_yes_no", Status: "0", Creator: "1", Updator: "", Remark: "系统是否列表"},
			{DictId: 5, DictName: "任务状态", DictType: "sys_job_status", Status: "0", Creator: "1", Updator: "", Remark: "任务状态列表"},
			{DictId: 6, DictName: "任务分组", DictType: "sys_job_group", Status: "0", Creator: "1", Updator: "", Remark: "任务分组列表"},
			{DictId: 7, DictName: "通知类型", DictType: "sys_notice_type", Status: "0", Creator: "1", Updator: "", Remark: "通知类型列表"},
			{DictId: 8, DictName: "系统状态", DictType: "sys_common_status", Status: "0", Creator: "1", Updator: "", Remark: "登录状态列表"},
			{DictId: 9, DictName: "操作类型", DictType: "sys_oper_type", Status: "0", Creator: "1", Updator: "", Remark: "操作类型列表"},
			{DictId: 10, DictName: "通知状态", DictType: "sys_notice_status", Status: "0", Creator: "1", Updator: "", Remark: "通知状态列表"},
		}

		list8 := []models.User{
			{1, "admin", "$2a$10$cKFFTCzGOvaIHHJY2K45Zuwt8TD6oPzYi4s5MzYIBAWCLL6ZhouP2", "zhangwj", "13818888888", "", "", "0", "1@qq.com", "0", "", 1, 1, 1, "1", "1", models.Model{}, "", ""},
		}

		list9 := []models.DictData{
			{DictId: 1, Sort: 0, DictLabel: "正常", DictValue: "0", DictType: "sys_normal_disable", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "系统正常"},
			{DictId: 2, Sort: 0, DictLabel: "停用", DictValue: "1", DictType: "sys_normal_disable", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "系统停用"},
			{DictId: 3, Sort: 0, DictLabel: "男", DictValue: "0", DictType: "sys_user_sex", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "性别男"},
			{DictId: 4, Sort: 0, DictLabel: "女", DictValue: "1", DictType: "sys_user_sex", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "性别女"},
			{DictId: 5, Sort: 0, DictLabel: "未知", DictValue: "2", DictType: "sys_user_sex", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "性别未知"},
			{DictId: 6, Sort: 0, DictLabel: "显示", DictValue: "0", DictType: "sys_show_hide", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "显示菜单"},
			{DictId: 7, Sort: 0, DictLabel: "隐藏", DictValue: "1", DictType: "sys_show_hide", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "隐藏菜单"},
			{DictId: 8, Sort: 0, DictLabel: "是", DictValue: "Y", DictType: "sys_yes_no", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "系统默认是"},
			{DictId: 9, Sort: 0, DictLabel: "否", DictValue: "N", DictType: "sys_yes_no", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "系统默认否"},
			{DictId: 10, Sort: 0, DictLabel: "正常", DictValue: "2", DictType: "sys_job_status", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "正常状态"},
			{DictId: 11, Sort: 0, DictLabel: "停用", DictValue: "1", DictType: "sys_job_status", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "停用状态"},
			{DictId: 12, Sort: 0, DictLabel: "默认", DictValue: "DEFAULT", DictType: "sys_job_group", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "默认分组"},
			{DictId: 13, Sort: 0, DictLabel: "系统", DictValue: "SYSTEM", DictType: "sys_job_group", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "系统分组"},
			{DictId: 14, Sort: 0, DictLabel: "通知", DictValue: "1", DictType: "sys_notice_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "通知"},
			{DictId: 15, Sort: 0, DictLabel: "公告", DictValue: "2", DictType: "sys_notice_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "公告"},
			{DictId: 16, Sort: 0, DictLabel: "正常", DictValue: "0", DictType: "sys_common_status", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "正常状态"},
			{DictId: 17, Sort: 0, DictLabel: "关闭", DictValue: "1", DictType: "sys_common_status", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "关闭状态"},
			{DictId: 18, Sort: 0, DictLabel: "新增", DictValue: "1", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "新增操作"},
			{DictId: 19, Sort: 0, DictLabel: "修改", DictValue: "2", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "修改操作"},
			{DictId: 20, Sort: 0, DictLabel: "删除", DictValue: "3", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "删除操作"},
			{DictId: 21, Sort: 0, DictLabel: "授权", DictValue: "4", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "授权操作"},
			{DictId: 22, Sort: 0, DictLabel: "导出", DictValue: "5", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "导出操作"},
			{DictId: 23, Sort: 0, DictLabel: "导入", DictValue: "6", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "导入操作"},
			{DictId: 24, Sort: 0, DictLabel: "强退", DictValue: "7", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "强退操作"},
			{DictId: 25, Sort: 0, DictLabel: "生成代码", DictValue: "8", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "生成操作"},
			{DictId: 26, Sort: 0, DictLabel: "清空数据", DictValue: "9", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "清空操作"},
			{DictId: 27, Sort: 0, DictLabel: "成功", DictValue: "0", DictType: "sys_notice_status", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "成功状态"},
			{DictId: 28, Sort: 0, DictLabel: "失败", DictValue: "1", DictType: "sys_notice_status", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "", Remark: "失败状态"},
			{DictId: 29, Sort: 0, DictLabel: "登录", DictValue: "10", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "1", Remark: "登录操作"},
			{DictId: 30, Sort: 0, DictLabel: "退出", DictValue: "11", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "1", Remark: ""},
			{DictId: 31, Sort: 0, DictLabel: "获取验证码", DictValue: "12", DictType: "sys_oper_type", CssClass: "", ListClass: "", IsDefault: "", Status: "0", Default: "", Creator: "1", Updator: "1", Remark: "获取验证码"},
		}

		list10 := []models.Setting{
			{1, "go-admin管理系统", "https://gitee.com/mydearzwj/image/raw/master/img/go-admin.png", models.Model{}},
		}

		list11 := []models.Job{
			{1, "接口测试", "DEFAULT", 1, "0/5 * * * * ", "http://localhost:8000", "", 1, 1, 1, 0, "", "", models.Model{}, ""},
			{2, "函数测试", "DEFAULT", 2, "0/5 * * * * ", "ExamplesOne", "参数", 1, 1, 1, 0, "", "", models.Model{}, ""},
		}

		err := tx.Create(list1).Error
		if err != nil {
			return err
		}
		err = tx.Create(list2).Error
		if err != nil {
			return err
		}

		err = tx.Create(list3).Error
		if err != nil {
			return err
		}

		err = tx.Create(list4).Error
		if err != nil {
			return err
		}

		err = tx.Create(list5).Error
		if err != nil {
			return err
		}

		err = tx.Create(list6).Error
		if err != nil {
			return err
		}

		err = tx.Create(list7).Error
		if err != nil {
			return err
		}

		err = tx.Create(list8).Error
		if err != nil {
			return err
		}

		err = tx.Create(list9).Error
		if err != nil {
			return err
		}

		err = tx.Create(list10).Error
		if err != nil {
			return err
		}

		err = tx.Create(list11).Error
		if err != nil {
			return err
		}

		if err := models.InitDb(tx); err != nil {
			return err
		}

		return tx.Create(&models.Migration{
			Version: version,
		}).Error
	})
}
