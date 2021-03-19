package rule

import (
	"hongxia_api/core/clRouter"
	"hongxia_api/core/clUserPool"
)

func InitRule() {

	clRouter.AddRule(clRouter.RouterRule{
		Ac:       "userLogin",
		Param:    []clRouter.RouterParam{
			clRouter.RouterParam{
				Key:   "user",
				Def:   "",
				PType: 0,
				Static: false,
			},
			clRouter.RouterParam{
				Key:   "pass",
				Def:   "",
				PType: 0,
				Static: false,
			},
		},
		Static:   false,
		Callback: func(_uInfo *clUserPool.ClNetUserInfo, _params map[string]string) string {
			return "OK"
		},
		Login:    false,
	})

}