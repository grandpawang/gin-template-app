package conf

import (
	"flag"
	"gin-template-app/pkg/cache/redis"
	"gin-template-app/pkg/database/minio"
	"gin-template-app/pkg/database/orm"
	xsql "gin-template-app/pkg/database/sql"
	"gin-template-app/pkg/log"
	"gin-template-app/pkg/net/http"
	"gin-template-app/pkg/net/mqtt"
	time "gin-template-app/pkg/time"

	"github.com/BurntSushi/toml"
)

// Redis config
type Redis struct {
	*redis.Config
	Expire time.Duration
}

// Config .
type Config struct {
	ORM        *orm.Config
	DB         *xsql.Config
	HTTP       *http.Config
	Log        *log.Config
	MQTT       *mqtt.Config
	Redis      *Redis
	Minio      *minio.Config
	HTTPClient *http.ClientConfig
}

var (
	confPath string
	// Conf config struct
	Conf = &Config{}
)

func init() {
	flag.StringVar(&confPath, "conf", "./cmd/config.toml", "default config path")
}

func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

// Init init conf
func Init() error {
	log.Info("🚀 make up config file")
	if confPath != "" {
		return local()
	}
	return nil
}
