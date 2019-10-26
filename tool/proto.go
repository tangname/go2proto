package main

type MessageType int32

// 消息类型;
const (
	// 数据库的实体;
	POES MessageType = 1
	// 请求;
	REQ MessageType = 3
	// 响应;
	RESP MessageType = 5
	// 其他类型;
	ELSE MessageType = 7
)

// 云,包含多个服务聚合;
type Cloud struct {
	Name     string     // 名称
	Services []*Service //服务集合
	Poes     *Dto       // 数据库实体集合
	// 导入的dto,Key为dto的包路径(小写),该集合用于汇总所有服务导入的dto包
	ImportDtos map[string]*Dto
}

/*
 服务
*/
type Service struct {
	Name             string             // 服务名称
	ServiceFunctions []*ServiceFunction // 方法列表
	RootDto          *Dto               // 根节点的dto,为服务对应的dto文件;
}

/*
导入信息
*/
type Import struct {
	PackageName string // 包名,若是引用了别名,则使用别名
	PackagePath string // 包的路径
}

/*
 Service对应的Dto信息
*/
type Dto struct {
	IsRoot       bool                //是否根节点;
	Name         string              // dto名称
	Imports      map[string]*Import  // dto中引用的其他的Dto,其中key为包的名称(非完整路径)|别名,例如:import "t1/commm",则Key为common
	Messages     map[string]*Message // 当前Dto的消息集合,key为struct的名称
	MessageIndex int32
	Dtos         map[string]*Dto // 导入包中关联的Dto,key为包的路径(小写)
}

func (this *Dto) Create() {
	this.Imports = make(map[string]*Import)
	this.Messages = make(map[string]*Message)
	this.Dtos = make(map[string]*Dto)
}

/*
ServiceFunction 接口定义的方法
*/
type ServiceFunction struct {
	Name        string
	ParamTypes  []string
	ResultTypes []string
	Comment     string
	Stream      bool
	PingPong    bool
}

/*
Message 接口定义的结构体
*/
type Message struct {
	Name          string
	Type          MessageType
	MessageFields []MessageField
	Index         int32
}

/*
MessageField 结构体字段属性
*/
type MessageField struct {
	Index     int
	FieldName string
	FieldType string
	Comment   string
}
