package Http

type Code int

const (
	Success2000                  Code = 2000
	NeedInvitationCode1005       Code = 1005  //需要邀请码
	DbError1999                  Code = 1999  //数据库错误
	EmailFormatErr2000           Code = 2001  //邮箱格式错误
	PwdLengthErr2006             Code = 2006  //密码长度不够
	AccountLock2011              Code = 2011  //账户被锁定
	NotLogin2013                 Code = 2013  //未登录
	PwdErr2017                   Code = 2017  //密码错误
	PhoneCodeError2018           Code = 2018  //验证码错误
	AccountAlreadyRegistered2074 Code = 2074  //账户早已经注册
	CodeExpire2076               Code = 2076  //验证码过期
	GtCodeError2082              Code = 2082  //人机验证错误
	NewPwdAndOldSame2094         Code = 2094  //新旧密码相同
	UserNotExist2106             Code = 2106  //用户不存在
	InvitationCodeError2124      Code = 2124  //无效验证码
	ParamsError3000              Code = 3000  //参数错误
	ParamsNull3003               Code = 3002  //参数为空
	ServiceNetworkBad4000        Code = 4000  //当前服务网络不稳定
	ServiceBusy4003              Code = 4003  //服务器繁忙
	PermissionDecline10000       Code = 10000 //用户没有权限
	InvalidInviter15002          Code = 15002 //邀请码无效
)

func (c Code) Code() int {
	switch c {
	case Success2000:
		return 2000
	case NeedInvitationCode1005:
		return 1005
	case DbError1999:
		return 1999
	case EmailFormatErr2000:
		return 2001
	case PwdLengthErr2006:
		return 2006
	case AccountLock2011:
		return 2011
	case NotLogin2013:
		return 2013
	case PwdErr2017:
		return 2017
	case PhoneCodeError2018:
		return 2018
	case AccountAlreadyRegistered2074:
		return 2074
	case CodeExpire2076:
		return 2076
	case GtCodeError2082:
		return 2082
	case NewPwdAndOldSame2094:
		return 2094
	case UserNotExist2106:
		return 2106
	case InvitationCodeError2124:
		return 2124
	case ParamsError3000:
		return 3000
	case ParamsNull3003:
		return 3002
	case ServiceNetworkBad4000:
		return 4000
	case ServiceBusy4003:
		return 4003
	case PermissionDecline10000:
		return 10000
	case InvalidInviter15002:
		return 15002
	default:
		return 0

	}
}
