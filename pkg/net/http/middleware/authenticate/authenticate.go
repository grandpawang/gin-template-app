package authenticate

import (
	"gin-template-app/pkg/ecode"
	"gin-template-app/pkg/net/http/ctx"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT Authentication MW
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code ecode.Codes
		var data interface{}

		code = ecode.OK
		t := c.GetHeader("Authorization")
		if t == "" {
			code = ecode.ParamsErr
		} else {
			userInfo, err := ParseToken(t)
			ctx.WithCurrentUser(c, userInfo)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = ecode.AccessTokenExpires
				default:
					code = ecode.SignCheckErr
				}
			}
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
