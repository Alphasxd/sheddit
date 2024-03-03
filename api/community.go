package api

import (
	"github.com/gin-gonic/gin"

	"sheddit/common"
	"sheddit/logic"
)

func GetCommunityCategory(ctx *gin.Context) {
	category := logic.GetCommunityCategory()
	common.Success(ctx, category)
}
