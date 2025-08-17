package usecase_test

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func makeGinCtx() *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c
}
