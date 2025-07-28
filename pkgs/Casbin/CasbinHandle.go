package Casbin

import (
	"github.com/gin-gonic/gin"
	"ngo2/pkgs/Default"
	"ngo2/pkgs/Jwt"
	"strings"
)

func Handler() gin.HandlerFunc {
	return c.Handler()
}
func (c *Casbin) Handler() gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, _ := context.Get(DefaultClaims)
		waitUse := claims.(*Jwt.Claims)
		obj := context.Request.URL.RequestURI()
		objs := strings.Split(obj, ObjOffset)
		realObj := strings.TrimPrefix(objs[0], DefaultPrefix)
		realObj = strings.Split(realObj, c.Prefix)[1]
		act := context.Request.Method
		sub := waitUse.AuthorityId
		if waitUse.AppRights == Admin {
			sub = AdminSub
		}
		_ = c.GetCasbin()
		enforceRs, _ := c.Enforcer.Enforce(sub, realObj, act)
		if enforceRs {
			context.Next()
		} else {
			Jwt.ErrorNoPrivileges.GetMsg(context, Default.StringEmpty).Abort()
			return
		}

	}
}
