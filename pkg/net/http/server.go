package http

import (
	"gin-template-app/pkg/log"
	"gin-template-app/pkg/net/http/middleware/core"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	jsoniterator "github.com/json-iterator/go"
)

var json = jsoniterator.ConfigCompatibleWithStandardLibrary

// Config http config
type Config struct {
	Host           string
	Port           int64
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

// New http server
func New() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	// add MW
	engine.Use(core.Core())
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// add prometheus
	engine.GET("/metrics", monitor())
	// engine.addRoute("GET", "/metadata", engine.metadata())
	// startPerf(engine)
	return engine
}

// Init Http Server
func Init(c *Config, engine *gin.Engine) {
	readTimeout := c.ReadTimeout * time.Second
	writeTimeout := c.WriteTimeout * time.Second
	endPoint := c.Host + ":" + strconv.FormatInt(c.Port, 10)
	maxHeaderBytes := c.MaxHeaderBytes
	log.Infoln("🚀 setup server: ")
	log.Infoln("	🎉 read timeout: ", readTimeout)
	log.Infoln("	🎉 write timeout: ", writeTimeout)
	log.Infoln("	🎉 end point: ", endPoint)
	log.Infoln("	🎉 max header bytes: ", maxHeaderBytes)

	s := &http.Server{
		Addr:           endPoint,
		Handler:        engine,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Panicf("gin(http) Start error(%v)", err)
			panic(err)
		}
	}()
}
