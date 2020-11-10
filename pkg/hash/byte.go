package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

/* 对字节数组取hash值 */

// Md5Byte byte md5
func Md5Byte(s []byte) string {
	h := md5.New()
	h.Write(s)
	return hex.EncodeToString(h.Sum(nil))
}

// Sha1Byte byte sha1
func Sha1Byte(s []byte) string {
	h := sha1.New()
	h.Write(s)
	return hex.EncodeToString(h.Sum(nil))
}

// Sha256Byte byte sha256
func Sha256Byte(s []byte) string {
	h := sha256.New()
	h.Write(s)
	return hex.EncodeToString(h.Sum(nil))
}

// Sha512Byte byte sha512
func Sha512Byte(s []byte) string {
	h := sha512.New()
	h.Write(s)
	return hex.EncodeToString(h.Sum(nil))
}
