package orm

import (
	"gin-template-app/pkg/ecode"
	"gin-template-app/pkg/log"
	xtime "gin-template-app/pkg/time"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql device
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ORM gorm database database mapping struct
type ORM struct {
	DB *gorm.DB
	c  *Config
}

// Config mysql config.
type Config struct {
	DSN         string         // orm DSN
	Debug       bool           // orm is debug
	Active      int            // pool
	Idle        int            // pool
	IdleTimeout xtime.Duration // connect max life time.
}

func connect(c *Config) (*gorm.DB, error) {
	// connect db
	log.WithFields(logrus.Fields{"connect": c.DSN, "debug": c.Debug}).Infoln("🚀 setup orm database")
	db, err := gorm.Open("mysql", c.DSN)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}

	// db setting
	db.SingularTable(true) // Table name is not plural
	db.DB().SetMaxIdleConns(c.Idle)
	db.DB().SetMaxOpenConns(c.Active)
	db.DB().SetConnMaxLifetime(time.Duration(c.IdleTimeout) / time.Second)
	db.LogMode(c.Debug)
	return db, nil
}

func init() {
	gorm.ErrRecordNotFound = ecode.NothingFound
}

// NewMysql initialization the gorm database mapping struct
func NewMysql(c *Config) *ORM {
	d, err := connect(c)
	if err != nil {
		panic(err)
	}
	return &ORM{DB: d, c: c}
}
