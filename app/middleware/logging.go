package middleware

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] %s %s %d %s \n",
			param.ClientIP,
			param.TimeStamp.Format(time.ANSIC),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	})
}

func SaveLogs() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}
