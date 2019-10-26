package dto

import "t1/stocking/common/createbasic"

type DemoRequest struct {
	UserId int64
}

type Demo1Response struct {
	DemoRequest
	createbasic.CreateBasic
	User  *DemoRequest
	Users []*DemoRequest
	test  int `reamark:"名称"` // 测试返回数据
	//User
}
