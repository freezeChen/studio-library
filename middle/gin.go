/*
   @Time : 2019-05-16 10:38
   @Author : frozenchen
   @File : gin
   @Software: studio
*/
package middle

import (
	"strings"
	"time"

	"github.com/freezeChen/studio-library/metadata"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token, X-File-Name")
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}
	}
}

func GeneralMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(metadata.GinStartTime, time.Now())
		ctx.Set(metadata.GinTraceId, strings.ReplaceAll(uuid.New().String(), "-", ""))
		ctx.Next()
	}
}
