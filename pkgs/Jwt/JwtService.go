package Jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kevinu2/ngo/v2/pkgs/RedisGo"
	"strconv"
	"time"
)

var j *JWT

type JWT struct {
	SigningKey []byte
	Target     string
	SessionId  string
	Redis      *RedisGo.Cacher
}

func init() {
	j = New()
}

func New() *JWT {
	return new(JWT)
}

func AddConfig(signed, platform, SessionId string, redis *RedisGo.Cacher) {
	j.AddConfig(signed, platform, SessionId, redis)
}

func (j *JWT) AddConfig(signKey, platform, SessionId string, redis *RedisGo.Cacher) {
	j.SigningKey = []byte(signKey)
	j.Target = platform
	j.SessionId = SessionId
	j.Redis = redis
}

func Init() *JWT {
	return j.Init()
}

func (j *JWT) Init() *JWT {
	return j
}

func Auth() gin.HandlerFunc {
	return j.Auth()
}

func (j *JWT) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			ErrorTokenEmpty.GetMsg(c, "").Abort()
			c.Abort()
			return
		}
		claims, err := j.Parse(token)
		if err != nil {
			if err == errors.New(TokenExpired) {
				ErrorTokenParse.GetMsg(c, err.Error()).Abort()
				return
			}
			ErrorTokenParse.GetMsg(c, err.Error()).Abort()
			return
		}

		sessionId := claims.SessionId
		jwtKey := fmt.Sprintf("%s_%s_%s", j.Target, j.SessionId, sessionId)
		userSessionKey, err := j.Redis.Exists(jwtKey)
		if err != nil || !userSessionKey {
			ErrorTokenExpire.GetMsg(c, err.Error()).Abort()
			return
		}

		pageExpiredTime, err := j.Redis.GetString("page_timeout")
		expireTime := time.Duration(-1) * time.Minute
		if err != nil {
			err = j.Redis.Set(jwtKey, token, int64(expireTime))
			if err != nil {
				ErrorTokenSet.GetMsg(c, err.Error()).Abort()
				return
			}
		} else {
			pageExpireTimeInt, _ := strconv.Atoi(pageExpiredTime)
			keyExists, err := j.Redis.Exists(jwtKey)
			if err != nil {
				ErrorTokenSet.GetMsg(c, err.Error()).Abort()
				return
			}
			refreshTime := time.Duration(pageExpireTimeInt) * time.Minute
			if keyExists {
				err = j.Redis.Expire(jwtKey, int64(refreshTime))
				if err != nil {
					ErrorTokenSet.GetMsg(c, err.Error()).Abort()
					return
				}
			}

		}
		c.Set("claims", claims)
		c.Next()
	}
}

func Create(claims Claims) (string, error) {
	return j.Create(claims)
}
func (j *JWT) Create(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func Parse(tokenString string) (*Claims, error) {
	return j.Parse(tokenString)
}

func (j *JWT) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New(TokenMalformed)
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errors.New(TokenExpired)
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New(TokenNotValidYet)
			} else {
				return nil, errors.New(TokenInvalid)
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.New(TokenInvalid)

	} else {
		return nil, errors.New(TokenInvalid)
	}
}

func Refresh(tokenString string) (string, error) {
	return j.Refresh(tokenString)
}

func (j *JWT) Refresh(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
		return j.Create(*claims)
	}
	return "", errors.New(TokenInvalid)
}
func Get(sessionId string) (string, error) {
	return j.Get(sessionId)
}
func (j *JWT) Get(sessionId string) (string, error) {
	userKey := fmt.Sprintf("%s_%s_%s", j.SigningKey, "session_id", sessionId)
	redisJWT, err := j.Redis.GetString(userKey)
	if err != nil {
		return "", err
	}
	var jwtList Blacklist
	err = json.Unmarshal([]byte(redisJWT), &jwtList)
	if err != nil {
		return "", err
	}

	return jwtList.Jwt, err
}

func Set(jwtList Blacklist, userId int, sessionId string, expiredTime int) (err error) {
	return j.Set(jwtList, userId, sessionId, expiredTime)
}

func (j *JWT) Set(jwtList Blacklist, userId int, sessionId string, expiredTime int) (err error) {
	userKey := fmt.Sprintf("%s_%s_%s", j.Target, "session_id", sessionId)
	userKeys := fmt.Sprintf("%s_%d", j.Target, userId)
	err = j.Redis.Set(userKey, jwtList.Jwt, int64(expiredTime))
	if err != nil {
		return err
	}

	_, err = j.Redis.HSet(userKeys, sessionId, userKey)
	if err != nil {
		return err
	}

	err = j.Redis.Expire(userKeys, int64(expiredTime))
	if err != nil {
		return err
	}
	return nil
}
