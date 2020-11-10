package minio

import (
	"gin-template-app/pkg/log"

	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
)

// Config .
type Config struct {
	Endpoint        string // end point
	AccessKeyID     string // access key id
	SecretAccessKey string // secret access key
	UseSSL          bool   // use ssl
}

// Client minio client
type Client struct {
	c *minio.Client
}

// New minio client
func New(c *Config) *Client {
	// Initialize minio client object.
	log.WithFields(logrus.Fields{"connect": c.Endpoint, "ssl": c.UseSSL}).Infoln("🚀 setup minio")
	minioClient, err := minio.New(c.Endpoint, c.AccessKeyID, c.SecretAccessKey, c.UseSSL)
	if err != nil {
		log.Fatalln(err)
	}

	cli := Client{c: minioClient}
	return &cli
}
