package file

import (
	"bytes"
	"compress/zlib"
	"io"
)

// ZlibEncode zlib 压缩
func ZlibEncode(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

// ZlibDecode zlib解压
func ZlibDecode(src []byte) []byte {
	b := bytes.NewReader(src)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}
