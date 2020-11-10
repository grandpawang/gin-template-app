package http

import (
	"gin-template-app/pkg/ecode"
	"gin-template-app/pkg/log"
	"gin-template-app/pkg/valid"

	"github.com/gin-gonic/gin"
)

// BindAndValid Bind request form data and validate them
func BindAndValid(c *gin.Context, data interface{}) error {
	// bind data
	err := c.Bind(data)
	if err != nil {
		log.Error("Bind() err(%v)", err)
		err = ecode.RequestErr
		return err
	}
	err = valid.Valid(data)
	if err != nil {
		log.Error("Valid() err(%v)", err)
		err = ecode.ParamsErr
		return err
	}
	return err
}
