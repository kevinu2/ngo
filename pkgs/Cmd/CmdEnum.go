package Cmd

// CmdType identifies a command family. Ctl/Exec are registration roles, not command identities.
type CmdType int

const (
	CmdPipeline CmdType = iota + 1001101
	CmdBehavior
	CmdSuppress
	CmdModel
	CmdFeature
	CmdRuleToggle
	CmdRuleApply
	CmdRuleReplay
	CmdIP
	CmdPort
	CmdTime
	CmdNtp
	CmdReboot
	CmdShutdown
	CmdHostNic
	CmdTrain
	CmdEval
	CmdFeedback
)

type ActionType int

const (
	ActionCreate ActionType = iota + 2008701
	ActionUpdate
	ActionDisable
	ActionDelete
	ActionDebug
	ActionReload
)

type RoleType string

const (
	RoleCtl  RoleType = "Ctl"
	RoleExec RoleType = "Exec"
)

type getter struct{ value any }

func (kind CmdType) Get() getter    { return getter{value: kind} }
func (kind ActionType) Get() getter { return getter{value: kind} }

func (g getter) Code() int {
	if g.Name() == "" {
		return 0
	}
	switch kind := g.value.(type) {
	case CmdType:
		return int(kind)
	case ActionType:
		return int(kind)
	default:
		return 0
	}
}
func (g getter) Name() string {
	switch kind := g.value.(type) {
	case CmdType:
		switch kind {
		case CmdPipeline:
			return "CmdPipeline"
		case CmdBehavior:
			return "CmdBehavior"
		case CmdSuppress:
			return "CmdSuppress"
		case CmdModel:
			return "CmdModel"
		case CmdFeature:
			return "CmdFeature"
		case CmdRuleToggle:
			return "CmdRuleToggle"
		case CmdRuleApply:
			return "CmdRuleApply"
		case CmdRuleReplay:
			return "CmdRuleReplay"
		case CmdIP:
			return "CmdIP"
		case CmdPort:
			return "CmdPort"
		case CmdTime:
			return "CmdTime"
		case CmdNtp:
			return "CmdNtp"
		case CmdReboot:
			return "CmdReboot"
		case CmdShutdown:
			return "CmdShutdown"
		case CmdHostNic:
			return "CmdHostNic"
		case CmdTrain:
			return "CmdTrain"
		case CmdEval:
			return "CmdEval"
		case CmdFeedback:
			return "CmdFeedback"
		}
	case ActionType:
		switch kind {
		case ActionCreate:
			return "CmdActionCreate"
		case ActionUpdate:
			return "CmdActionUpdate"
		case ActionDisable:
			return "CmdActionDisable"
		case ActionDelete:
			return "CmdActionDelete"
		case ActionDebug:
			return "CmdActionDebug"
		case ActionReload:
			return "CmdActionReload"
		}
	}
	return ""
}

type roleGetter string

func (g getter) Role() roleGetter {
	if _, ok := g.value.(CmdType); !ok {
		return ""
	}
	return roleGetter(g.Name())
}
func (r roleGetter) Exec() string {
	if r == "" {
		return ""
	}
	return string(r) + string(RoleExec)
}
func (r roleGetter) Ctl() string {
	if r == "" {
		return ""
	}
	return string(r) + string(RoleCtl)
}
func (r roleGetter) name(role RoleType) string {
	switch role {
	case RoleCtl:
		return r.Ctl()
	case RoleExec:
		return r.Exec()
	default:
		return ""
	}
}
