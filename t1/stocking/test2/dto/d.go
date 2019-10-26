package dto

import (
	"t1/stocking/common/usercommon"
)

//type User struct {
//	UserName string
//}

type Demo1Request struct {
	Users  []usercommon.UserCommon
	UserId int64
}

type Demo1Response struct {
	usercommon.UserCommon
	test int `reamark:"名称"` // 测试返回数据
	//User
}
