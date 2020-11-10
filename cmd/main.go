package main

import (
	"flag"
	"gin-template-app/conf"
	"gin-template-app/pkg/ecode"
	"gin-template-app/pkg/log"
	"gin-template-app/route"
	"gin-template-app/service/template"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("😟 conf.Init() error(%v)", err)
		panic(err)
	}
	// log init
	log.NewLog(conf.Conf.Log)
	// ecode init
	ecode.Init()
	// service init
	templateSvr := template.New(conf.Conf)
	// http init
	route.Init(templateSvr, *conf.Conf)
	log.Infoln("😀 up server ok")
	// signal init
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("server exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
