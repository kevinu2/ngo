package Http

type UserInfo struct {
	UserId    uint64 `json:"user_id" xorm:"pk"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	SessionId string `json:"sessionId"`
}
