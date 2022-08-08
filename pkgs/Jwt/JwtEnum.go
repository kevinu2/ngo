package Jwt

import (
	"github.com/gin-gonic/gin"
)

type ErrorType uint8

const (
	ErrorTokenEmpty ErrorType = iota + 1
	ErrorIsBlacklist
	ErrorTokenParse
	ErrorTokenExpire
	ErrorRequestTimeout
	ErrorTokenSet
)

func (et ErrorType) GetMsg(c *gin.Context, msg string) *gin.Context {
	switch et {
	case ErrorTokenEmpty:
		c.JSON(401, gin.H{"result": -1, "reason": "未登录", "data": "token error"})
		return c
	case ErrorIsBlacklist:
		c.JSON(401, gin.H{"result": -1, "reason": "非法访问", "data": "token in blacklist"})
		return c
	case ErrorTokenParse:
		c.JSON(401, gin.H{"result": -1, "reason": "解析错误", "data": "parse token err " + msg})
		return c
	case ErrorTokenExpire:
		c.JSON(401, gin.H{"result": -1, "reason": "授权已过期", "data": "token expired"})
		return c
	case ErrorRequestTimeout:
		c.JSON(401, gin.H{"result": -1, "reason": "请求超时，请重新登录", "data": "request timeout"})
		return c
	case ErrorTokenSet:
		c.JSON(401, gin.H{"result": -1, "reason": "授权失败", "data": "set token err " + msg})
		return c
	default:
		return c
	}
}
