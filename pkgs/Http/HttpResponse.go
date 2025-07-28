package Http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct{}

func (Response) Right(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, RightResponse{
		Result: data,
	})
}

func (Response) Error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, ErrorResponse{
		Err: Error{
			Code:   ParamsError3000.Code(),
			Msg:    err.Error(),
			Result: nil,
		},
		Result: nil,
	})
}

func (Response) ErrorWithCode(c *gin.Context, code int, err error) {
	c.JSON(http.StatusOK, ErrorResponse{
		Err: Error{
			Code:   code,
			Msg:    err.Error(),
			Result: "",
		},
		Result: nil,
	})
}

func (Response) Page(c *gin.Context, data []interface{}, count int64) {
	c.JSON(http.StatusOK, PageResponse{
		Items: data,
		Count: count,
	})
}
