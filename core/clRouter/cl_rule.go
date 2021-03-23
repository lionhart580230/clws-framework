package clRouter

import (
	"net/url"
)

// 管理路由规则

func (this *RouterRule) CheckParam(_paramList string) (bool, map[string]string) {
	var vMap = make(map[string]string)
	vlist, err := url.ParseQuery(_paramList)
	if err != nil {
		return false, vMap
	}

	for _, pInfo := range this.Param {
		val := ""
		v, exist := vlist[pInfo.Key]

		if !exist {
			if pInfo.Static {
				return false, vMap
			}
			val = pInfo.Def
		} else {
			val = v[0]
		}

		val = checkParam(pInfo.PType, val)
		if val == "VALUE_IS_ERROR" {
			return false, vMap
		}
		vMap[pInfo.Key] = val
	}
	return true, vMap
}