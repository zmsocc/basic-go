package web

import (
	"github.com/gin-gonic/gin"
)

type jwtHandler struct {
}

func NewJwtHandler() jwtHandler {
	return jwtHandler{}
}

func ExtractToken(ctx *gin.Context) string {
	panic("implement me")
}
