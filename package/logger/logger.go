package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] %s %s %d %s", params.ClientIP, params.TimeStamp.Format(time.RFC822),
				params.Path, params.Method, params.StatusCode, params.Latency)
		})
}
