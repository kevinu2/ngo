package Casbin

import (
	"bytes"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"github.com/kevinu2/ngo2/pkgs/Error"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

var c *Casbin

const path = "rbac_model.conf"

type Casbin struct {
	Enforcer *casbin.Enforcer
	GormDB   *gorm.DB
	Table    string
	Prefix   string
}

func init() {
	c = New()
}

func New() *Casbin {
	return c.New()
}

func (c *Casbin) New() *Casbin {
	if c != nil {
		return c
	}
	v := new(Casbin)
	return v
}
func AddConfig(gorm *gorm.DB, table, prefix string) {
	c.AddConfig(gorm, table, prefix)
}
func (c *Casbin) AddConfig(gorm *gorm.DB, table, prefix string) {
	c.GormDB = gorm
	c.Table = table
	c.Prefix = prefix
	c.Enforcer, _ = c.GetCasbin()
}

func GetCasbin() (*casbin.Enforcer, error) { return c.GetCasbin() }
func (c *Casbin) GetCasbin() (*casbin.Enforcer, error) {
	if c.Enforcer != nil {
		return c.Enforcer, nil
	}
	a, err := gormAdapter.NewAdapterByDB(c.GormDB)
	if err != nil {
		return c.Enforcer, Error.ErrorNotFound.GetMsg(err.Error())
	}
	e, err := casbin.NewEnforcer("rbac_model.conf", a)
	if err != nil {
		return c.Enforcer, Error.ErrorNotMatch.GetMsg(err.Error())
	}
	_ = e.LoadPolicy()
	return c.Enforcer, nil
}

func UpdateCasbinAuths(authInfos []CasbinAuthInfo, roleId int, tx *gorm.DB) error {
	return c.UpdateCasbinAuths(authInfos, roleId, tx)
}
func (c *Casbin) UpdateCasbinAuths(authInfos []CasbinAuthInfo, roleId int, tx *gorm.DB) error {
	conn := c.GormDB.Table(c.Table)
	err := c.DeleteCasbinRuleByRoleId(strconv.Itoa(roleId))
	if err != nil {
		tx.Rollback()
		return err
	}
	err = c.BatchSaveCasbinAuth(conn, authInfos)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func DeleteCasbinRuleByRoleId(id string) error { return c.DeleteCasbinRuleByRoleId(id) }
func (c *Casbin) DeleteCasbinRuleByRoleId(id string) error {
	conn := c.GormDB.Table(c.Table)
	err := conn.Where("v0 = ?", id).Delete(&CasbinModel{}).Error
	return err
}

func DeleteCasbinRuleById(id string) error { return c.DeleteCasbinRuleById(id) }
func (c *Casbin) DeleteCasbinRuleById(id string) error {
	conn := c.GormDB.Table(c.Table)
	err := conn.Where("id = ?", id).Delete(&CasbinModel{}).Error
	return err
}

func BatchSaveCasbinAuth(conn *gorm.DB, authInfos []CasbinAuthInfo) error {
	return c.BatchSaveCasbinAuth(conn, authInfos)
}
func (c *Casbin) BatchSaveCasbinAuth(conn *gorm.DB, authInfos []CasbinAuthInfo) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("insert into %s(p_type,ptype,v0,v1,v2) values", c.Table)
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, authInfo := range authInfos {
		if i == len(authInfos)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s');", authInfo.PType, authInfo.PPType,
				authInfo.AuthorityId, authInfo.Url, authInfo.Method))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s'),", authInfo.PType, authInfo.PPType,
				authInfo.AuthorityId, authInfo.Url, authInfo.Method))
		}
	}
	return conn.Exec(buffer.String()).Error
}

func ClearCasbinAuth(v int, p ...string) bool { return c.ClearCasbinAuth(v, p...) }
func (c *Casbin) ClearCasbinAuth(v int, p ...string) bool {
	rs, _ := c.Enforcer.RemoveFilteredPolicy(v, p...)
	return rs
}

func AddCasbinAuth(authInfo CasbinAuthInfo) error { return c.AddCasbinAuth(authInfo) }
func (c *Casbin) AddCasbinAuth(authInfo CasbinAuthInfo) error {
	var (
		cs CasbinModel
		n  int64
	)
	err := c.GormDB.Table(c.Table).Where("v0 = ? AND v1 = ? AND v3 = ?", authInfo.AuthorityId, authInfo.Url, authInfo.Method).Find(&cs).Count(&n).Error
	if err != nil {
		return Error.ErrorNotFound.GetMsg(err.Error())
	}
	if n > 0 {
		return Error.ErrorNotRequired.GetMsg("")

	}
	if strings.TrimSpace(authInfo.PType) == `` {
		cs.PType = `p`
	} else {
		cs.PType = authInfo.PType
	}
	cs.AuthorityId = authInfo.AuthorityId
	cs.Url = authInfo.Url
	cs.Method = authInfo.Method

	err = c.GormDB.Table(c.Table).Create(&cs).Error
	return err
}

func UpdateUserCasbin(authorityId string, casbinInfos []CasbinInfo) error {
	return c.UpdateUserCasbin(authorityId, casbinInfos)
}
func (c *Casbin) UpdateUserCasbin(authorityId string, casbinInfos []CasbinInfo) error {
	c.ClearCasbinAuth(0, authorityId)
	for _, v := range casbinInfos {
		cm := CasbinModel{
			Id:          0,
			PType:       "p",
			AuthorityId: authorityId,
			Url:         v.Url,
			Method:      v.Method,
		}
		cms := make([]CasbinModel, 0)
		cms = append(cms, cm)
		addFlag := c.AddUserCasbin(cm)
		if addFlag == false {
			return Error.ErrorNotRequired.GetMsg("")

		}
	}
	return nil
}

func AddUserCasbin(cm CasbinModel) bool { return c.AddUserCasbin(cm) }
func (c *Casbin) AddUserCasbin(cm CasbinModel) bool {
	rs, _ := c.Enforcer.AddPolicy(cm.AuthorityId, cm.Url, cm.Method)
	return rs
}

func GetAllCasbin(limit, offset int) ([]CasbinModel, error) {
	return c.GetAllCasbin(limit, offset)
}
func (c *Casbin) GetAllCasbin(limit, offset int) ([]CasbinModel, error) {
	var cs []CasbinModel

	if err := c.GormDB.Table(c.Table).Limit(limit).Offset(offset).Find(&cs).Error; err != nil {
		return nil, err
	}
	return cs, nil
}

func UpdateCasbinAuth(param CasbinAuthInfo) error { return c.UpdateCasbinAuth(param) }
func (c *Casbin) UpdateCasbinAuth(param CasbinAuthInfo) error {
	var cs CasbinModel
	rows := c.GormDB.Table(c.Table).Where("v0 = ? AND v1 = ? AND v2 = ?", param.AuthorityId, param.Url, param.Method).Find(&cs).RowsAffected
	if rows > 0 {
		return Error.ErrorNotRequired.GetMsg("")
	}

	err := c.GormDB.Table(c.Table).Where("id = ?", param.Id).Updates(map[string]interface{}{
		"ptype": param.PType,
		"v0":    param.AuthorityId,
		"v1":    param.Url,
		"v2":    param.Method,
	}).Error
	return err
}

func (c *Casbin) clearCasbin(v int, p ...string) bool {
	rs, _ := c.Enforcer.RemoveFilteredPolicy(v, p...)
	return rs
}

func (c *Casbin) DeleteCasbinAuthById(id string) error {
	var cs CasbinModel
	err := c.GormDB.Table(c.Table).Where("id = ?", id).First(&cs).Error
	if err != nil {
		return err
	}
	err = c.GormDB.Table(c.Table).Delete(cs).Error
	c.clearCasbin(1, cs.AuthorityId, cs.Url, cs.Method)
	return err
}

func (c *Casbin) MatchRoleCasbinRule(param CheckAuthParam) error {
	conn := c.GormDB.Table(c.Table)
	rows := conn.Where("vo = ? and v1 = ? and v2 = ?", param.RuleId, param.Url, param.Method).RowsAffected
	if rows <= 0 {
		return Error.ErrorNotFound.GetMsg("")
	}
	return nil
}
