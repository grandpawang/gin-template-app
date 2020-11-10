package cron

import (
	"fmt"

	"github.com/robfig/cron"
)

// @yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 1 1 *
// @monthly               | Run once a month, midnight, first of month | 0 0 1 * *
// @weekly                | Run once a week, midnight between Sat/Sun  | 0 0 * * 0
// @daily (or @midnight)  | Run once a day, midnight                   | 0 0 * * *
// @hourly                | Run once an hour, beginning of hour        | 0 * * * *
type cr struct {
	c *cron.Cron
	t string
}

var (
	cronList map[string]cr
)

func init() {
	cronList = make(map[string]cr)

}

// Add a cron
// example Add("print", "* * * * * ?", func() {
//		fmt.Println("1s Timeout!")
//	})")
func Add(name string, time string, f func()) error {
	if _, exist := cronList[name]; exist {
		return fmt.Errorf("cron %v is exist", name)
	}
	c := cron.New()
	ok := c.AddFunc(time, f)
	if ok != nil {
		return fmt.Errorf("cron start fail, time is %v", time)
	}
	cronList[name] = cr{
		c: c,
		t: time,
	}
	c.Start()
	return nil
}

// Del a cron
// example Add("print")
func Del(name string) error {
	cr, exist := cronList[name]
	if !exist {
		return fmt.Errorf("cron %v not found", name)
	}

	cr.c.Stop()
	delete(cronList, name)

	return nil
}

// Edit a cron
// Edit("print", "* * * * * /2", func() {
//		fmt.Println("2s Timeout!")
//	})
func Edit(name string, time string, f func()) error {
	if e := Del(name); e != nil {
		return e
	}
	if e := Add(name, time, f); e != nil {
		return e
	}
	return nil
}

// Once Add a Once Cron
func Once(name string, time string, f func()) error {
	if e := Add(name, time, func() {
		Del(name)
		f()
	}); e != nil {
		return e
	}
	return nil
}
