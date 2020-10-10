package tools

import (
	"strconv"
	"strings"
)

type Ids struct {
	Ids []int
}

// 解析URL中批量id
func IdsStrToIdsIntGroup(keys string) []int {
	IDS := make([]int, 0)
	ids := strings.Split(keys, ",")
	for i := 0; i < len(ids); i++ {
		ID, _ := strconv.Atoi(ids[i])
		IDS = append(IDS, ID)
	}
	return IDS
}
