package cron

import (
	"fmt"
	"gopkg.in/go-playground/assert.v1"
	"testing"
	"time"
)

func TestCorn(t *testing.T) {
	e := Add("print", "* * * * * ?", func() {
		fmt.Println("1s Timeout!")
	})
	assert.Equal(t, e, nil)
	time.Sleep(5 * time.Second)

	e = Edit("print", "0/2 * * * * ?", func() {
		fmt.Println("2s Timeout!")
	})
	assert.Equal(t, e, nil)
	time.Sleep(10 * time.Second)

	e = Del("print")
	assert.Equal(t, e, nil)
	time.Sleep(5 * time.Second)
}

func TestOnceCron(t *testing.T) {
	e := Once("@monthly", "@monthly", func() {
		fmt.Println("Once!!!")
	})
	assert.Equal(t, e, nil)
	time.Sleep(5 * time.Second)
}
