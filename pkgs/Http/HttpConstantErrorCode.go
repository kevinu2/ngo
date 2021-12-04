package Http

//服务错误码
var (
	NeedInvitationCode1005       = 1005  //需要邀请码
	DbError1999                  = 1999  //数据库错误
	EmailFormatErr2000           = 2000  //邮箱格式错误
	PwdLengthErr2006             = 2006  //密码长度不够
	AccountLock2011              = 2011  //账户被锁定
	NotLogin2013                 = 2013  //未登录
	PwdErr2017                   = 2017  //密码错误
	PhoneCodeError2018           = 2018  //验证码错误
	AccountAlreadyRegistered2074 = 2074  //账户早已经注册
	CodeExpire2076               = 2076  //验证码过期
	GtCodeError2082              = 2082  //人机验证错误
	NewPwdAndOldSame2094         = 2094  //新旧密码相同
	UserNotExist2106             = 2106  //用户不存在
	InvitationCodeError2124      = 2124  //无效验证码
	ParamsError3000              = 3000  //参数错误
	ParamsNull3003               = 3002  //参数为空
	ServiceNetworkBad4000        = 4000  //当前服务网络不稳定
	ServiceBusy4003              = 4003  //服务器繁忙
	PermissionDecline10000       = 10000 //用户没有权限
	InvalidInviter15002          = 15002 //邀请码无效
)
