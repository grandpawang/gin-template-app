package template

import (
	"gin-template-app/conf"
	"gin-template-app/dao/template"
	"gin-template-app/pkg/database/minio"
)

// Service struct
type Service struct {
	c     *conf.Config
	dao   *template.Dao
	minio *minio.Client
	mqtt  *MQTTOrder
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:     c,
		dao:   template.New(c),
		minio: minio.New(c.Minio),
		mqtt:  NewMQTTOrder(c.MQTT),
	}
	return s
}
