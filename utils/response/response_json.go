package response

import (
	"online-food/dto"

	"github.com/gin-gonic/gin"
)

func ToResponseJson(ctx *gin.Context, code int, status string, message string, data interface{}) {
	ctx.JSON(code, dto.WebResponse{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	},
	)
}
