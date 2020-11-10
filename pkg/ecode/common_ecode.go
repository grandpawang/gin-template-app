package ecode

// All common ecode
var (
	// System ecode
	OK                 = add(0)    // 正确
	SignCheckErr       = add(-3)   // API校验密匙错误
	ParamsErr          = add(-5)   // 参数错误
	NoLogin            = add(-101) // 账号未登录
	RequestErr         = add(-400) // 请求参数错误
	Unauthorized       = add(-401) // 未认证
	AccessDenied       = add(-403) // 访问权限不足
	NothingFound       = add(-404) // 不存在该参数实体
	ServerErr          = add(-500) // 服务器错误
	ServiceUnavailable = add(-503) // 过载保护服务暂不可用
	Deadline           = add(-504) // 服务调用超时
	LimitExceed        = add(-509) // 超出限制
	AccessTokenExpires = add(-658) // Token过期
	NotDelAsscoiation  = add(-659) // 存在关联关系不能删除

	// Box ecode
	ExistBox        = add(-10000) // 已存在盒子
	NotExistBox     = add(-10001) // 不存在盒子
	BoxNotOnline    = add(-10002) // 盒子未上线
	ExistOrder      = add(-10003) // 已存在订单号
	BoxExistMission = add(-10004) // 盒子存在未处理任务
	BoxOrderError   = add(-10005) // 订单信息下载异常
	BoxPrintError   = add(-10006) // 屏幕异常
	BoxPowerError   = add(-10007) // 电量异常

	// MW ecode
	ExistMW          = add(-11000) // 已存在基站
	NotExistMW       = add(-11001) // 不存在基站
	ExistMWCfgName   = add(-11002) // 已存在基站配置名称
	MWStatusStop     = add(-11003) // 基站已经停止传达盒子工作任务
	MWStatusOutline  = add(-11004) // 基站网络异常
	MWStatusConfing  = add(-11005) // 基站正在配置或者已经配置失败，无法继续执行该指令
	MWStatusUpdating = add(-11006) // 基站正在更新或者已经升级失败，无法继续执行该指令

	// version ecode
	NotExistVersion    = add(-12001) // 不存在该硬件本身版本
	NotExistBoxVersion = add(-12002) // 不存在存放在该硬件的盒子版本
	NotValidVersion    = add(-12003) // 上传的固件非官方固件

	// process ecode
	NotExistProcess  = add(-13001) // 不存在流程
	ExistProcessName = add(-13002) // 已存在流程名称
)
