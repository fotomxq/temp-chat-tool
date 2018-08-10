package ModForm

import "encoding/json"

//处理form提交的数据模块

//解析POST的JSON结构体
// 仅限于处理[]string{}类型
//param value string 数据
//return []string 数据组
func ModFormGetArray(value string) []string {
	res := []string{}
	if value == "" {
		return res
	}
	err := json.Unmarshal([]byte(value), &res)
	if err != nil {
		res = []string{
			value,
		}
	}
	return res
}
