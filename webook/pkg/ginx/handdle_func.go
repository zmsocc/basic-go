package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WrapReq[T any](fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		res, err := fn(ctx, req)
		if err != nil {

		}
		ctx.JSON(http.StatusOK, res)
	}
}

//type Result Result
