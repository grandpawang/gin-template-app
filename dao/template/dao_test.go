package template

import (
	"flag"
	"gin-template-app/conf"
	"os"
	"testing"
)

var (
	d   *Dao
	ids []uint
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../../cmd/config.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
