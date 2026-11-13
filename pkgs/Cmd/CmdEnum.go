package Cmd

// CmdType identifies Cmd controller/executor endpoints. Identities are assigned sequentially from 1001101.
type CmdType int

const (
	CmdPipelineCtl CmdType = iota + 1001101
	CmdPipelineExec
	CmdBehaviorCtl
	CmdBehaviorExec
	CmdSuppressCtl
	CmdSuppressExec
	CmdModelCtl
	CmdModelExec
	CmdFeatureCtl
	CmdFeatureExec
	CmdRuleToggleCtl
	CmdRuleToggleExec
	CmdRuleApplyCtl
	CmdRuleApplyExec
	CmdRuleReplayCtl
	CmdRuleReplayExec
	CmdIPCtl
	CmdIPExec
	CmdPortCtl
	CmdPortExec
	CmdTimeCtl
	CmdTimeExec
	CmdNtpCtl
	CmdNtpExec
	CmdRebootCtl
	CmdRebootExec
	CmdShutdownCtl
	CmdShutdownExec
	CmdHostNicCtl
	CmdHostNicExec
	CmdTrainCtl
	CmdTrainExec
	CmdEvalCtl
	CmdEvalExec
	CmdFeedbackCtl
	CmdFeedbackExec
)

func (st CmdType) GetCode() int {
	switch st {
	case CmdPipelineCtl:
		return 1001101
	case CmdPipelineExec:
		return 1001102
	case CmdBehaviorCtl:
		return 1001103
	case CmdBehaviorExec:
		return 1001104
	case CmdSuppressCtl:
		return 1001105
	case CmdSuppressExec:
		return 1001106
	case CmdModelCtl:
		return 1001107
	case CmdModelExec:
		return 1001108
	case CmdFeatureCtl:
		return 1001109
	case CmdFeatureExec:
		return 1001110
	case CmdRuleToggleCtl:
		return 1001111
	case CmdRuleToggleExec:
		return 1001112
	case CmdRuleApplyCtl:
		return 1001113
	case CmdRuleApplyExec:
		return 1001114
	case CmdRuleReplayCtl:
		return 1001115
	case CmdRuleReplayExec:
		return 1001116
	case CmdIPCtl:
		return 1001117
	case CmdIPExec:
		return 1001118
	case CmdPortCtl:
		return 1001119
	case CmdPortExec:
		return 1001120
	case CmdTimeCtl:
		return 1001121
	case CmdTimeExec:
		return 1001122
	case CmdNtpCtl:
		return 1001123
	case CmdNtpExec:
		return 1001124
	case CmdRebootCtl:
		return 1001125
	case CmdRebootExec:
		return 1001126
	case CmdShutdownCtl:
		return 1001127
	case CmdShutdownExec:
		return 1001128
	case CmdHostNicCtl:
		return 1001129
	case CmdHostNicExec:
		return 1001130
	case CmdTrainCtl:
		return 1001131
	case CmdTrainExec:
		return 1001132
	case CmdEvalCtl:
		return 1001133
	case CmdEvalExec:
		return 1001134
	case CmdFeedbackCtl:
		return 1001135
	case CmdFeedbackExec:
		return 1001136
	default:
		return 0
	}
}

type ActionType int

const (
	ActionCreate ActionType = iota + 2008701
	ActionUpdate
	ActionDisable
	ActionDelete
	ActionDebug
	ActionReload
)

func (at ActionType) GetCode() int {
	switch at {
	case ActionCreate:
		return 2008701
	case ActionUpdate:
		return 2008702
	case ActionDisable:
		return 2008703
	case ActionDelete:
		return 2008704
	case ActionDebug:
		return 2008705
	case ActionReload:
		return 2008706
	default:
		return 0
	}
}
