package paginator

type Page struct {
	List      interface{} `json:"list"`
	Total     int64       `json:"total"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
}
