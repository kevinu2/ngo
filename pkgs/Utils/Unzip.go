package Utils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func UnzipByte(data []byte) (unzipData []byte, err error) {
	if len(data) == 0 {
		return
	}
	closer, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return
	}
	defer closer.Close()

	return ioutil.ReadAll(closer)
}
