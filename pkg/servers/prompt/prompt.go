package prompt

type Language int

const (
	EN Language = iota
	CN
)

type Prompt int

const (
	Unavailable Prompt = iota
	QuerySuccess
	QueryFailed
	CreateSuccess
	CreateFailed
	UpdatedSuccess
	UpdateFailed
	DeleteSuccess
	DeleteFailed
	NotFound
	OperationTimeout
	IncorrectData
)

var language = CN
var message = map[Prompt]map[Language]string{
	Unavailable:      {EN: "unavailable", CN: "未知错误"},
	QuerySuccess:     {EN: "query success", CN: "查询成功"},
	QueryFailed:      {EN: "query failed", CN: "查询失败"},
	CreateSuccess:    {EN: "create success", CN: "创建成功"},
	CreateFailed:     {EN: "create failed", CN: "创建失败"},
	UpdatedSuccess:   {EN: "update success", CN: "更新成功"},
	UpdateFailed:     {EN: "update failed", CN: "更新失败"},
	DeleteSuccess:    {EN: "delete success", CN: "删除成功"},
	DeleteFailed:     {EN: "delete failed", CN: "删除失败"},
	NotFound:         {EN: "not found", CN: "未找到相关内容"},
	OperationTimeout: {EN: "operation timeout", CN: "操作超时"},
	IncorrectData:    {EN: "incorrect data", CN: "数据不正确"},
}

func (p Prompt) String() string {
	prom := message[p]
	if len(prom) == 0 {
		return message[Unavailable][language]
	}
	return prom[language]
}
