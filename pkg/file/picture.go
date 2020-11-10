package file

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"image/png"
	"io"
	"log"
)

const (
	dx = 296
	dy = 128
)

func lzwencode(src []byte) []byte {
	var in bytes.Buffer
	w := lzw.NewWriter(&in, lzw.LSB, 8)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func gzipencode(src []byte) []byte {
	var in bytes.Buffer
	w := gzip.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func flateencode(src []byte) []byte {
	var in bytes.Buffer
	w, _ := flate.NewWriter(&in, flate.DefaultCompression)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

// GetPNGImagePic get png images
func GetPNGImagePic(r io.Reader) (res []byte, err error) {
	// decode image
	img, err := png.Decode(r)
	if err != nil {
		log.Fatal(err)
		return
	}
	// get image bytes
	res = make([]byte, dx*dy/8)
	imgX := img.Bounds().Max.X
	imgY := img.Bounds().Max.Y
	for x := 0; x < dx; x++ {
		var byt byte = 0
		for y := 0; y < dy; y++ {
			var r uint32
			if x < imgX && (127-y) < imgY {
				r, _, _, _ = img.At(x, 127-y).RGBA()
			} else {
				r = 0xff
			}
			tmp := uint8(0)
			tmp |= uint8(r & 0xff)

			if tmp >= 0xbb {
				byt |= 1 << ((127 - y) % 8)
			}
			if (127-y)%8 == 0 {
				res[x*dy/8+y/8] = byt
				// fmt.Printf("%x ", byt)
				byt = 0
			}
		}
		// fmt.Println()
	}
	// fmt.Println(imgX, imgY)
	res = ZlibEncode(res)
	// fmt.Println(len(res))
	return res, nil
}
