package test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"

	jsoniterator "github.com/json-iterator/go"
)

var json = jsoniterator.ConfigCompatibleWithStandardLibrary

// Header like string <--> string
type Header map[string]string

// JSONData like string <-->
type JSONData map[string]interface{}

func request(w *httptest.ResponseRecorder, req *http.Request, router *gin.Engine) []byte {
	// handler
	// router.ServeHTTP(w, req)

	// get response
	result := w.Result()
	defer result.Body.Close()

	// get body
	body, _ := ioutil.ReadAll(result.Body)

	return body
}

func setHeaders(obj interface {
	Set(key string, value string)
}, param Header) {
	for k, v := range param {
		obj.Set(k, v)
	}
}

func setupRequest(uri, method string, body io.Reader) (w *httptest.ResponseRecorder, req *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest(method, uri, body)
}

// get get uri base
func get(uri string, param url.Values, headers Header, router *gin.Engine) []byte {
	// build request params
	w, req := setupRequest(uri, "GET", nil)

	// get request add param
	req.URL.RawQuery = param.Encode()

	// set headers
	setHeaders(&req.Header, headers)

	// request uri
	return request(w, req, router)
}

// post post uri base
func post(uri string, param JSONData, headers Header, router *gin.Engine) []byte {
	// build request params
	b, _ := json.Marshal(param)
	w, req := setupRequest(uri, "POST", bytes.NewBuffer(b))

	// set headers
	req.Header.Set("Content-Type", "application/json")
	setHeaders(&req.Header, headers)

	// request uri
	return request(w, req, router)
}

// Get method
func Get(uri string, router *gin.Engine) []byte {
	return get(uri, url.Values{}, Header{}, router)
}

// GetWithHeader by Headers
func GetWithHeader(uri string, headers Header, router *gin.Engine) []byte {
	return get(uri, url.Values{}, headers, router)
}

// GetForm by form
func GetForm(uri string, param url.Values, router *gin.Engine) []byte {
	return get(uri, param, Header{}, router)
}

// GetFormWithHeader by form, header
func GetFormWithHeader(uri string, param url.Values, headers Header, router *gin.Engine) []byte {
	return get(uri, param, headers, router)
}

// Post method
func Post(uri string, param JSONData, router *gin.Engine) []byte {
	return post(uri, param, Header{}, router)
}

// PostWithHeader by Header
func PostWithHeader(uri string, param JSONData, headers Header, router *gin.Engine) []byte {
	return post(uri, param, headers, router)
}
