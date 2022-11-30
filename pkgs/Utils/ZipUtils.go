package Utils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

type zip struct {
}

func Zip() zip {
	return zip{}
}

func (zip) UnzipByte(data []byte) (unzipData []byte, err error) {
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
