package main

import (
	"flag"
	"gin-template-app/conf"
	"gin-template-app/pkg/database/orm"
	"gin-template-app/pkg/log"
)

func main() {
	flag.Parse()

	if err := conf.Init(); err != nil {
		log.Error("😟 conf.Init() error(%v)", err)
		panic(err)
	}
	db := orm.NewMysql(conf.Conf.ORM)
	db.DB.AutoMigrate()
	log.Info("😀 database migrate ok")

	// minio := minio.New(conf.Conf.Minio)
	// if err := minio.NewBucket(models.BoxImageBucket, "us-east-1"); err != nil {
	// 	log.Error("minio error(%v)", err)
	// }

	log.Info("😀 minio new bucket ok")
}
