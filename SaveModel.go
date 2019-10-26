package main

type SaveModel struct {
	ServiceName      string             // 服务名称
	ServiceFunctions []*ServiceFunction // 方法列表
	Reqs             []*Message         // 消息集合
	Resps            []*Message         // 消息集合
	Else             []*Message         // 消息集合
	Poes             []*Message         // 消息集合
}
