package http

import (
	"bytes"
	"gin-template-app/pkg/ecode"
	"gin-template-app/pkg/log"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type jsonResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// JSON http response
func JSON(c *gin.Context, data interface{}, err error) {
	code := http.StatusOK
	bcode := ecode.Cause(err)
	if err != nil {
		code = http.StatusInternalServerError
	}
	c.JSON(code, jsonResp{
		Code:    bcode.Code(),
		Message: bcode.Message(),
		Data:    data,
	})
}

// File http response
func File(c *gin.Context, data io.Reader, err error) {
	code := http.StatusOK
	if err != nil {
		log.Error("response File err(%v)", err)
		JSON(c, nil, err)
		return
	}
	bs, err := ioutil.ReadAll(data)
	if err != nil {
		JSON(c, nil, err)
		return
	}
	file := bytes.NewReader(bs)
	c.DataFromReader(code, file.Size(), "applicate/file", file, nil)
}

// Redirect http redirect
func Redirect(c *gin.Context, data string, err error) {
	code := http.StatusFound
	if err != nil {
		log.Error("response File err(%v)", err)
		JSON(c, nil, err)
		return
	}
	c.Redirect(code, data)
	c.Abort()
}
