package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrResp struct {
	Error ErrorInfo `json:"error"`
}

func EncodeError(gCtx *gin.Context, err error) {
	var (
		code    int
		errInfo ErrorInfo
	)

	if ei, ok := err.(ErrorInfo); ok {
		code = ei.StatusCode()
		errInfo = ei
	} else {
		errInfo = ErrInternal(err.Error())
		code = errInfo.StatusCode()
	}

	gCtx.AbortWithStatusJSON(code, ErrResp{
		Error: errInfo,
	})
}

func EncodeResp(gCtx *gin.Context, data interface{}) {
	if data != nil {
		gCtx.JSON(http.StatusOK, data)
	} else {
		gCtx.JSON(http.StatusOK, struct{}{})
	}
}
