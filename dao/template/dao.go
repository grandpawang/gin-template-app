package template

import (
	"context"
	"gin-template-app/conf"
	"gin-template-app/pkg/cache/redis"
	"gin-template-app/pkg/database/orm"
	"gin-template-app/pkg/database/sql"
	xsql "gin-template-app/pkg/database/sql"
	xhttp "gin-template-app/pkg/net/http"
	"time"

	jsoniterator "github.com/json-iterator/go"

	"github.com/pkg/errors"
)

var json = jsoniterator.ConfigCompatibleWithStandardLibrary

// Dao .
type Dao struct {
	c        *conf.Config
	orm      *orm.ORM
	db       *xsql.DB
	redis    *redis.Pool
	client   *xhttp.Client
	hwExpire int32
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:        c,
		orm:      orm.NewMysql(c.ORM),
		db:       xsql.NewMySQL(c.DB),
		redis:    redis.NewPool(c.Redis.Config),
		client:   xhttp.NewClient(c.HTTPClient),
		hwExpire: int32(time.Duration(c.Redis.Expire) / time.Second),
	}

	return
}

// Close close the resource.
func (d *Dao) Close() {
	if d.db != nil {
		d.orm.DB.Close()
	}
}

// BeginTran .
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	if tx, err = d.db.Begin(c); err != nil {
		err = errors.WithStack(err)
	}
	return
}
