package http

import (
	"bytes"
	"context"
	"fmt"
	"gin-template-app/pkg/net/breaker"
	"gin-template-app/pkg/stat"
	xtime "gin-template-app/pkg/time"
	"io"
	"net"
	xhttp "net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	pkgerr "github.com/pkg/errors"
)

const (
	_minRead = 16 * 1024 // 16kb

)

var (
	clientStats = stat.HTTPClient
)

// ClientConfig is http client conf.
type ClientConfig struct {
	Dial      xtime.Duration
	Timeout   xtime.Duration
	KeepAlive xtime.Duration
	Breaker   *breaker.Config
	URL       map[string]*ClientConfig
	Host      map[string]*ClientConfig
}

// Client is http client.
type Client struct {
	conf   *ClientConfig
	client *xhttp.Client
	dialer *net.Dialer

	urlConf  map[string]*ClientConfig
	hostConf map[string]*ClientConfig
	mutex    sync.RWMutex
	breaker  *breaker.Group
}

// NewClient new a http client.
func NewClient(c *ClientConfig) *Client {
	client := new(Client)
	client.conf = c
	client.dialer = &net.Dialer{
		Timeout:   time.Duration(c.Dial),
		KeepAlive: time.Duration(c.KeepAlive),
	}

	// wraps RoundTripper for tracer
	client.client = &xhttp.Client{}
	client.urlConf = make(map[string]*ClientConfig)
	client.hostConf = make(map[string]*ClientConfig)
	client.breaker = breaker.NewGroup(c.Breaker)
	if c.Timeout <= 0 {
		panic("must config http timeout!!!")
	}
	for uri, cfg := range c.URL {
		client.urlConf[uri] = cfg
	}
	for host, cfg := range c.Host {
		client.hostConf[host] = cfg
	}
	return client
}

// SetConfig set client config.
func (client *Client) SetConfig(c *ClientConfig) {
	client.mutex.Lock()
	if c.Timeout > 0 {
		client.conf.Timeout = c.Timeout
	}
	if c.KeepAlive > 0 {
		client.dialer.KeepAlive = time.Duration(c.KeepAlive)
		client.conf.KeepAlive = c.KeepAlive
	}
	if c.Dial > 0 {
		client.dialer.Timeout = time.Duration(c.Dial)
		client.conf.Timeout = c.Dial
	}
	if c.Breaker != nil {
		client.conf.Breaker = c.Breaker
		client.breaker.Reload(c.Breaker)
	}
	for uri, cfg := range c.URL {
		client.urlConf[uri] = cfg
	}
	for host, cfg := range c.Host {
		client.hostConf[host] = cfg
	}
	client.mutex.Unlock()
}

// NewRequest new http request with method, uri, ip, values and headers.
func (client *Client) NewRequest(method, uri string, params url.Values) (req *xhttp.Request, err error) {
	enc, err := client.sign(params)
	if err != nil {
		err = pkgerr.Wrapf(err, "uri:%s,params:%v", uri, params)
		return
	}
	ru := uri
	if enc != "" {
		ru = uri + "?" + enc
	}
	if method == xhttp.MethodGet {
		req, err = xhttp.NewRequest(xhttp.MethodGet, ru, nil)
	} else {
		req, err = xhttp.NewRequest(xhttp.MethodPost, uri, strings.NewReader(enc))
	}
	if err != nil {
		err = pkgerr.Wrapf(err, "method:%s,uri:%s", method, ru)
		return
	}
	const (
		_contentType = "Content-Type"
		_urlencoded  = "application/x-www-form-urlencoded"
		_userAgent   = "User-Agent"
	)
	if method == xhttp.MethodPost {
		req.Header.Set(_contentType, _urlencoded)
	}
	req.Header.Set(_userAgent, "coint.server")
	return
}

// Get issues a GET to the specified URL.
func (client *Client) Get(c context.Context, uri string, params url.Values, res interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodGet, uri, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res)
}

// Post issues a Post to the specified URL.
func (client *Client) Post(c context.Context, uri string, params url.Values, res interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodPost, uri, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res)
}

// RESTfulGet issues a RESTful GET to the specified URL.
func (client *Client) RESTfulGet(c context.Context, uri string, params url.Values, res interface{}, v ...interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodGet, fmt.Sprintf(uri, v...), params)
	if err != nil {
		return
	}
	return client.Do(c, req, res, uri)
}

// RESTfulPost issues a RESTful Post to the specified URL.
func (client *Client) RESTfulPost(c context.Context, uri string, params url.Values, res interface{}, v ...interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodPost, fmt.Sprintf(uri, v...), params)
	if err != nil {
		return
	}
	return client.Do(c, req, res, uri)
}

// Raw sends an HTTP request and returns bytes response
func (client *Client) Raw(c context.Context, req *xhttp.Request, v ...string) (bs []byte, err error) {
	var (
		ok      bool
		code    string
		cancel  func()
		resp    *xhttp.Response
		config  *ClientConfig
		timeout time.Duration
		uri     = fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.Host, req.URL.Path)
	)
	// NOTE fix prom & config uri key.
	if len(v) == 1 {
		uri = v[0]
	}
	// breaker
	brk := client.breaker.Get(uri)
	if err = brk.Allow(); err != nil {
		code = "breaker"
		clientStats.Incr(uri, code)
		return
	}
	defer client.onBreaker(brk, &err)
	// stat
	now := time.Now()
	defer func() {
		clientStats.Timing(uri, int64(time.Since(now)/time.Millisecond))
		if code != "" {
			clientStats.Incr(uri, code)
		}
	}()
	// get config
	// 1.url config 2.host config 3.default
	client.mutex.RLock()
	if config, ok = client.urlConf[uri]; !ok {
		if config, ok = client.hostConf[req.Host]; !ok {
			config = client.conf
		}
	}
	client.mutex.RUnlock()
	// timeout
	deliver := true
	timeout = time.Duration(config.Timeout)
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := time.Until(deadline); ctimeout < timeout {
			// deliver small timeout
			timeout = ctimeout
			deliver = false
		}
	}
	if deliver {
		c, cancel = context.WithTimeout(c, timeout)
		defer cancel()
	}
	req = req.WithContext(c)
	if resp, err = client.client.Do(req); err != nil {
		err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		code = "failed"
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= xhttp.StatusBadRequest {
		err = pkgerr.Errorf("incorrect http status:%d host:%s, url:%s", resp.StatusCode, req.URL.Host, realURL(req))
		code = strconv.Itoa(resp.StatusCode)
		return
	}
	if bs, err = readAll(resp.Body, _minRead); err != nil {
		err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		return
	}
	return
}

// Do sends an HTTP request and returns an HTTP json response.
func (client *Client) Do(c context.Context, req *xhttp.Request, res interface{}, v ...string) (err error) {
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil {
		if err = json.Unmarshal(bs, res); err != nil {
			err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

// JSON sends an HTTP request and returns an HTTP json response.
func (client *Client) JSON(c context.Context, req *xhttp.Request, res interface{}, v ...string) (err error) {
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil {
		if err = json.Unmarshal(bs, res); err != nil {
			err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

func (client *Client) onBreaker(breaker breaker.Breaker, err *error) {
	if err != nil && *err != nil {
		breaker.MarkFailed()
	} else {
		breaker.MarkSuccess()
	}
}

// sign calc appkey and appsecret sign.
func (client *Client) sign(params url.Values) (query string, err error) {
	client.mutex.RLock()
	client.mutex.RUnlock()
	if params == nil {
		params = url.Values{}
	}
	tmp := params.Encode()
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	// var b bytes.Buffer
	// b.WriteString(tmp)
	// mh := md5.Sum(b.Bytes())
	// // query
	var qb bytes.Buffer
	qb.WriteString(tmp)
	// qb.WriteString("&sign=")
	// qb.WriteString(hex.EncodeToString(mh[:]))
	query = qb.String()
	return
}

// realUrl return url with http://host/params.
func realURL(req *xhttp.Request) string {
	if req.Method == xhttp.MethodGet {
		return req.URL.String()
	} else if req.Method == xhttp.MethodPost {
		ru := req.URL.Path
		if req.Body != nil {
			rd, ok := req.Body.(io.Reader)
			if ok {
				buf := bytes.NewBuffer([]byte{})
				buf.ReadFrom(rd)
				ru = ru + "?" + buf.String()
			}
		}
		return ru
	}
	return req.URL.Path
}

// readAll reads from r until an error or EOF and returns the data it read
// from the internal buffer allocated with a specified capacity.
func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}
