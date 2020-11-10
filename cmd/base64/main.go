package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"os"
)

func main(){
	bs ,_ := ioutil.ReadFile("./gbbmn-box.bin")
	bss := base64.StdEncoding.EncodeToString(bs)
	ioutil.WriteFile("out.txt", bytes.NewBufferString(bss).Bytes(), os.ModePerm)
	return
}