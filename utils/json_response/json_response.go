package json_response

import "github.com/gin-gonic/gin"

type APIResponse struct {
	CodeSystem   string      `json:"codeSystem"`
	Code         string      `json:"code"`
	Message      string      `json:"message,omitempty"`
	MessageError string      `json:"messageError,omitempty"`
	Result       interface{} `json:"result,omitempty"`
}

func Success(ctx *gin.Context, codeSystem string, result interface{}) {
	ctx.JSON(200, APIResponse{
		CodeSystem: codeSystem,
		Code:       "00",
		Message:    "success",
		Result:     result,
	})
	panic(nil)
}

func Error(ctx *gin.Context, codeSystem string, httpCode int, code string, message string) {
	ctx.JSON(httpCode, APIResponse{
		CodeSystem:   codeSystem,
		Code:         code,
		MessageError: message,
		Result:       "",
	})
	panic(nil)
}
