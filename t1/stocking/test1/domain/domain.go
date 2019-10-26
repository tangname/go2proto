package domain

import (
	"t1/stocking/poes"
	"t1/stocking/test1/dto"
)

type Test1Domain struct {
}

// 测试方法;
func (this *Test1Domain) Hello(request *dto.DemoRequest) (tt *dto.Demo1Response, err error) {
	return &dto.Demo1Response{}, nil
}

// 测试方法1;
func (this *Test1Domain) Hello1(request *dto.DemoRequest) (tt *poes.DeskBasic) {
	return &poes.DeskBasic{}
}

// 测试方法2;
func (this *Test1Domain) Hello2(request *dto.DemoRequest) (tt *dto.Demo1Response, err error) {
	return &dto.Demo1Response{}, nil
}
