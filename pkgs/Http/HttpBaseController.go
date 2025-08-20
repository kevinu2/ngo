package Http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinu2/ngo/v2/pkgs/Db/RedisDB"
	"github.com/kevinu2/ngo/v2/pkgs/Log"
)

func LoginFilter() gin.HandlerFunc {

	return func(c *gin.Context) {
		//h := request.Header{}
		userId := c.GetUint64("userId")
		if userId > 0 {
			c.Next()
			return
		}
		sessionId, err := c.Cookie(SessionId)
		if err != nil {
			//c.AbortWithError(200,err)
			Log.Logger().Errorf("err:%v", err.Error())
			c.JSON(http.StatusOK, ErrorResponse{Err: Error{Code: CodeUnLogin, Msg: err.Error()}})
			c.Abort()
			return
		}
		Log.Logger().Infof("session_id:%+v", sessionId)

		if sessionId == "" {
			Log.Logger().Error("sessionId 为空")
			//c.AbortWithError(200, errors.New("获取token为空"))
			c.JSON(http.StatusOK, ErrorResponse{Err: Error{Code: CodeUnLogin, Msg: "获取session_id为空"}})
			c.Abort()
			return
		}
		//TODO getUserByKey from redis
		redisClient := RedisDB.GetDB()
		r, err := redisClient.GetString(sessionId)
		if err != nil {
			Log.Logger().Errorf("get redis err:%v", err.Error())
			//c.AbortWithError(200, errors.New("获取token失败"))
			c.JSON(http.StatusOK, ErrorResponse{Err: Error{Code: CodeUnLogin, Msg: "redis获取session_id为空"}})
			c.Abort()
			return
		}
		Log.Logger().Infof("get redis result:%+v", r)
		var user UserInfo
		err = json.Unmarshal([]byte(r), &user)
		if err != nil {
			c.JSON(http.StatusOK, ErrorResponse{Err: Error{Code: CodeUnLogin, Msg: "反序列化user失败"}})
			c.Abort()
			Log.Logger().Errorf("Marshal error:%v", err.Error())
			return
		}
		Log.Logger().Infof("get user from cookie:%+v", user)
		c.Set("user", r)
		c.Set("userId", user.UserId)
	}

}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求
		c.Next()
	}
}
