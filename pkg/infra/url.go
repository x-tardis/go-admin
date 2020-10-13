package infra

import (
	"strconv"
	"strings"
)

// ParseIdsGroup 解析URL中以','分隔的批量id,
func ParseIdsGroup(keys string) []int {
	ids := strings.Split(keys, ",")
	IDS := make([]int, 0, len(ids))
	for i := 0; i < len(ids); i++ {
		id, err := strconv.Atoi(strings.Trim(ids[i], " "))
		if err != nil {
			continue
		}
		IDS = append(IDS, id)
	}
	return IDS
}
