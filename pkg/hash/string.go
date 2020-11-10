package hash

/* 对字符串取hash值 */

// Md5String string md5
func Md5String(s string) string {
	return Md5Byte([]byte(s))
}

// Sha1String string sha1
func Sha1String(s string) string {
	return Sha1Byte([]byte(s))
}

// Sha256String string sha256
func Sha256String(s string) string {
	return Sha256Byte([]byte(s))
}

// Sha512String string sha512
func Sha512String(s string) string {
	return Sha512Byte([]byte(s))
}
