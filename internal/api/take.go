package api

import (
	"codesignal/datatypes"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckStatus(ctx *gin.Context, limiter *datatypes.RateManager) {
	userMethod := ctx.Query("method")
	userRoute := ctx.Query("route")
	data, found := limiter.CheckStatus(ctx, userMethod, userRoute)
	if !found {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid input data or endpoint not initialized"))
		return
	}
	tokens := data.AvailableTokens()
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"available_tokens": tokens,
		"allow":            tokens > 0,
	})
}

func TakeApiGetUserId(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Hello World with user id")
}

func PatchUserId(ctx *gin.Context) {
	ctx.String(http.StatusOK, "this is patch user id")
}

func PostUserInfo(ctx *gin.Context) {
	ctx.String(http.StatusOK, "this is post user info")
}
