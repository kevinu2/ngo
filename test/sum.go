package main

import (
	"crypto/md5"
	"fmt"
	"hash/adler32"
	"hash/crc32"
)

func main() {

	data := []byte("敏感信息/重要文件泄漏[")
	h1 := adler32.Checksum(data)
	h2 := crc32.ChecksumIEEE(data)
	h3 := md5.Sum(data)
	fmt.Println(h1, h2, h3)
	fmt.Printf("%x\n%x\n%x\n", h1, h2, h3)
}
