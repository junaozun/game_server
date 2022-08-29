package middleware

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							// Here is the "broken pipe" error so just ignore it.
							c.Abort()
							return
						}
					}
				}

				trace := string(debug.Stack())
				fmt.Println("gin handleErrors panic:", err, trace)
				log.Printf("[gin] handleErrors panic=%v stack=%s", err, trace)
				var (
					errMsg string
					ok     bool
				)
				if errMsg, ok = err.(string); ok {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code":    http.StatusInternalServerError,
						"message": "system error, " + errMsg,
					})
					return
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code":    500,
						"message": "system error",
					})
					return
				}
			}
		}()
		c.Next()
	}
}
