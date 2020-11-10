package ctx

import (
	"github.com/gin-gonic/gin"
)

// UserInfo save in jwt
type UserInfo struct {
	ID       uint
	Username string
	Name     string
}

// currentUser context key
const currentUser = "current_user"

// WithCurrentUser gin with context
func WithCurrentUser(c *gin.Context, value interface{}) {
	c.Set(currentUser, value)
}

// CurrentUserFromCtx CurrentUser From gin Ctx
func CurrentUserFromCtx(c *gin.Context) (UserInfo, bool) {
	val, ok := c.Value(currentUser).(UserInfo)
	return val, ok
}
