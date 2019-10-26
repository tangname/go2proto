package domain

import "t1/stocking/test2/dto"

type Test2Domain struct {
}

// 测试方法1111;
func (this *Test2Domain) Hello(request *dto.DemoRequest) (tt *dto.DemoResponse, err error) {
	return &dto.DemoResponse{}, nil
}

// 测试方法;
func (this *Test2Domain) Hello111(request *dto.DemoRequest) *dto.DemoResponse {
	return &dto.DemoResponse{}
}
