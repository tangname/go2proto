package poes

// 柜台资料;
type DeskBasic struct {
	DeskId         int64  `remark:"柜台序号" gorm:"column:DeskId;type:bigint(20);primary_key"`
	CharacterId    int64  `remark:"角色序号" gorm:"column:CharacterId;type:bigint(20)"`
	DeskName       string `remark:"柜台名称" gorm:"column:DeskName;type:varchar(50)"`
	GroupTypeDk    int64  `remark:"柜台分组(SettingDictionary.DictId)" gorm:"column:GroupTypeDk;type:bigint(20)"`
	GroupTypeDv    string `remark:"柜台分组(SettingDictionary.DictName)" gorm:"column:GroupTypeDv;type:varchar(50)"`
	ChargeUserId   int64  `remark:"责任人员序号" gorm:"column:ChargeUserId;type:bigint(20)"`
	ChargeUser     string `remark:"负责人员" gorm:"column:ChargeUser;type:varchar(50)"`
	IsDefault      int32  `remark:"是否默认(枚举)" gorm:"column:IsDefault;type:int(11)"`
	LastPickupTime string `remark:"最后领货时间" gorm:"column:LastPickupTime;type:datetime(6)"`
	LastReturnTime string `remark:"最后退货时间" gorm:"column:LastReturnTime;type:datetime(6)"`
	State          int32  `remark:"柜台状态(枚举)" gorm:"column:State;type:int(11)"`
}

// 柜台领退单;
type DeskPickretOrderBasic struct {
	PickretId    int64  `remark:"领退单序号" gorm:"column:PickretId;type:bigint(20);primary_key"`
	CharacterId  int64  `remark:"角色序号" gorm:"column:CharacterId;type:bigint(20)"`
	DeskId       int64  `remark:"柜台序号" gorm:"column:DeskId;type:bigint(20)"`
	PickretCode  string `remark:"领退单号" gorm:"column:PickretCode;type:varchar(50)"`
	PickretType  int32  `remark:"领退类型(枚举)" gorm:"column:PickretType;type:int(11)"`
	ChargeUserId int64  `remark:"(账户验证)负责人员序号(领退人员)" gorm:"column:ChargeUserId;type:bigint(20)"`
	ChargeUser   string `remark:"(账户验证)负责人员(领退人员)" gorm:"column:ChargeUser;type:varchar(50)"`
	Quantity     int32  `remark:"货品总数" gorm:"column:Quantity;type:int(11)"`
	CostPrice    int64  `remark:"成本价(元)精度2位" gorm:"column:CostPrice;type:bigint(20)"`
	LabelPrice   int64  `remark:"标签价(元)精度2位" gorm:"column:LabelPrice;type:bigint(20)"`
	GoldWeight   int64  `remark:"净金重(g)精度4位" gorm:"column:GoldWeight;type:bigint(20)"`
	Weight       int64  `remark:"货重(g)精度4位" gorm:"column:Weight;type:bigint(20)"`
	CreateUserId int64  `remark:"创建人序号" gorm:"column:CreateUserId;type:bigint(20)"`
	CreateUser   string `remark:"创建人员" gorm:"column:CreateUser;type:varchar(50)"`
	CreateTime   string `remark:"创建时间" gorm:"column:CreateTime;type:datetime(6)"`
	Stamp        string `remark:"时间戳(分区)" gorm:"column:Stamp;type:datetime(6)"`
	State        int32  `remark:"领货退货单状态(枚举)" gorm:"column:State;type:int(11)"`
}
