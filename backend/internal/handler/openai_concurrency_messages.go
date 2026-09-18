package handler

import (
	"fmt"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/gin-gonic/gin"
)

func requestConcurrencyLimitMessage(_ *gin.Context, limit int) string {
	if limit > 0 {
		return fmt.Sprintf("request concurrency limit reached: max %d", limit)
	}
	return "request concurrency limit reached"
}

func authSubjectConcurrencyLimitMessage(subject middleware2.AuthSubject) string {
	return requestConcurrencyLimitMessage(nil, subject.Concurrency)
}
