package utils

import (
	bytes2 "bytes"

	"github.com/gin-gonic/gin"
)

func Failed(c *gin.Context, applies ...func(b *Base)) {
	baseResponse := Failure()
	for _, apply := range applies {
		apply(baseResponse)
	}

	respond(c, baseResponse)
}

// ErrorMessage ...
func respond(c *gin.Context, baseResponse *Base) {
	c.JSON(baseResponse.StatusCode, baseResponse)
}

type Response[T any] struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Data            T      `json:"data,omitempty"`
}

func (r *Response[T]) SetToSuccess() {
	r.ResponseCode = "2000000"
	r.ResponseMessage = "Success"
}

func (r *Response[T]) SetToSuccessCreated() {
	r.ResponseCode = "2010000"
	r.ResponseMessage = "Success"
}

func PadLeft(str string, pad string, length int) string {
	var buffer bytes2.Buffer
	for i := 0; i < (length - len(str)); i = i + len(pad) {
		buffer.WriteString(pad)
	}
	result := buffer.String() + str
	return result[0:length]
}
