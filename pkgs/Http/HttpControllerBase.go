package Http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"ngo/pkgs/Db/RedisDB"
	log "ngo/pkgs/Log"
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
			log.Logger().Errorf("err:%v", err.Error())
			c.JSON(http.StatusOK, Response{Err: Error{Code: CodeUnLogin, Msg: err.Error()}})
			c.Abort()
			return
		}
		log.Logger().Infof("session_id:%+v", sessionId)

		if sessionId == "" {
			log.Logger().Error("sessionId 为空")
			//c.AbortWithError(200, errors.New("获取token为空"))
			c.JSON(http.StatusOK, Response{Err: Error{Code: CodeUnLogin, Msg: "获取session_id为空"}})
			c.Abort()
			return
		}
		//TODO getUserByKey from redis
		redisClient := RedisDB.GetDB()
		r, err := redisClient.GetString(sessionId)
		if err != nil {
			log.Logger().Errorf("get redis err:%v", err.Error())
			//c.AbortWithError(200, errors.New("获取token失败"))
			c.JSON(http.StatusOK, Response{Err: Error{Code: CodeUnLogin, Msg: "redis获取session_id为空"}})
			c.Abort()
			return
		}
		log.Logger().Infof("get redis result:%+v", r)
		var user UserInfo
		err = json.Unmarshal([]byte(r), &user)
		if err != nil {
			c.JSON(http.StatusOK, Response{Err: Error{Code: CodeUnLogin, Msg: "反序列化user失败"}})
			c.Abort()
			log.Logger().Errorf("Marshal error:%v", err.Error())
			return
		}
		log.Logger().Infof("get user from cookie:%+v", user)
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
