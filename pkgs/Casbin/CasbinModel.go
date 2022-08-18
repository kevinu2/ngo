package Casbin

type CasbinAuthInfo struct {
	Id          string `form:"id" json:"id"`
	PType       string `form:"ptype" json:"ptype"`
	PPType      string `form:"p_type" json:"p_type"`
	AuthorityId string `form:"role_id" json:"role_id"`
	Url         string `form:"url" json:"url"`
	Method      string `form:"method" json:"method"`
}

type CasbinModel struct {
	Id          uint   `json:"id" gorm:"column:id"`
	PType       string `json:"ptype" gorm:"column:ptype"`
	PPType      string `json:"p_type"  gorm:"column:p_type"`
	AuthorityId string `json:"role_id" gorm:"column:v0"`
	Url         string `json:"url" gorm:"column:v1"`
	Method      string `json:"method" gorm:"column:v2"`
}

type CasbinInfo struct {
	Url    string `json:"url"`
	Method string `json:"method"`
}

type CasbinRangeInfo struct {
	//Tbg        string `form:"tbg" json:"tbg"`
	//Ted        string `form:"ted" json:"ted"`
	//OrderField string `form:"order_field" json:"order_field"`
	//Order      string `form:"order" json:"order"`
	Offset int `form:"offset" json:"offset"`
	Limit  int `form:"limit" json:"limit"`
}

// CasbinInReceive structure for input parameters
type CasbinInReceive struct {
	AuthorityId string       `form:"auth_id" json:"auth_id" binding:"required"`
	CasbinInfos []CasbinInfo `form:"casbin_infos" json:"casbin_infos"`
}

type CasbinPolicyPathResponse struct {
	CasbinInfo []CasbinInfo `json:"casbin_info"`
}

type CheckAuthParam struct {
	RuleId string `json:"rule_id"`
	Url    string `json:"url"`
	Method string `json:"method"`
}
