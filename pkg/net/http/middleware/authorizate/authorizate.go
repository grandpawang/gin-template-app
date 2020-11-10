package authorizate

import (
	"gin-template-app/pkg/ecode"
	"gin-template-app/pkg/log"
	"gin-template-app/pkg/net/http/ctx"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CasbinMW casbin MW
func CasbinMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := ecode.OK
		var data interface{} = nil

		user, ok := ctx.CurrentUserFromCtx(c)
		if !ok {
			code = ecode.NoLogin
			return
		}
		if user.Username == "root" {
			c.Next()
			return
		}
		uid := user.ID

		p := c.Request.URL.Path
		m := c.Request.Method
		b, err := CsbinCheckPermission(strconv.FormatUint(uint64(uid), 10), p, m)
		log.Infoln("🐕 CasbinMW", p, m, b)
		if err != nil {
			code = ecode.Unauthorized
			data = err.Error()
		} else if !b {
			code = ecode.AccessDenied
		}

		if code.Equal(ecode.OK) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  code.Message(),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
