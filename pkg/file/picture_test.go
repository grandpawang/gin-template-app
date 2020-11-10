package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceAddBoxConfig(t *testing.T) {
	convey.Convey("service - AddBoxConfig", t, func(ctx convey.C) {

		ctx.Convey("add box config", func(ctx convey.C) {
			file, _ := os.Open("../../test/1.png")
			byts, _ := GetPNGImagePic(file)
			// ioutil.WriteFile("test", byts, os.ModePerm)
			fmt.Println(len(byts))
			file.Close()

			file, _ = os.Open("../../test/2.png")
			byts, _ = GetPNGImagePic(file)
			fmt.Println(len(byts))
			file.Close()

			file, _ = os.Open("../../test/3.png")
			byts, _ = GetPNGImagePic(file)
			fmt.Println(len(byts))
			file.Close()
		})
	})
}
