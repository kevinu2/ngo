package Jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinu2/ngo2/pkgs/Http"
)

type ErrorType uint8

const (
	ErrorTokenEmpty ErrorType = iota + 1
	ErrorIsBlacklist
	ErrorTokenParse
	ErrorTokenExpire
	ErrorRequestExpire
	ErrorTokenSet
	ErrorNoPrivileges
)

func (et ErrorType) GetMsg(c *gin.Context, msg string) *gin.Context {
	switch et {
	case ErrorTokenEmpty:
		c.JSON(401, Http.Response{Err: Http.Error{Code: 401, Msg: "未登录或非法访问", Result: msg}, Result: Http.NoResult})
		return c
	case ErrorIsBlacklist:
		c.JSON(401, Http.Response{Err: Http.Error{Code: 401, Msg: "您的帐户异地登陆或令牌失效", Result: msg}, Result: Http.NoResult})
		return c
	case ErrorTokenParse:
		c.JSON(401, Http.Response{Err: Http.Error{Code: 401, Msg: "Token解析失败", Result: msg}, Result: Http.NoResult})
		return c
	case ErrorTokenExpire:
		c.JSON(401, Http.Response{Err: Http.Error{Code: 401, Msg: "授权已过期", Result: msg}, Result: Http.NoResult})
		return c
	case ErrorRequestExpire:
		c.JSON(401, Http.Response{Err: Http.Error{Code: 401, Msg: "请求超时，请重新登录", Result: msg}, Result: Http.NoResult})
		return c
	case ErrorTokenSet:
		c.JSON(401, Http.Response{Err: Http.Error{Code: 401, Msg: "授权失败", Result: msg}, Result: Http.NoResult})
		return c
	case ErrorNoPrivileges:
		c.JSON(403, Http.Response{Err: Http.Error{Code: 403, Msg: "权限不足，请联系管理员", Result: msg}, Result: Http.NoResult})
		return c
	default:
		return c
	}
}
