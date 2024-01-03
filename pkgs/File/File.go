package File

import "os"

func Exist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func Write(path string, data []byte, mode uint32) error {
	f, err := os.OpenFile(path, os.O_WRONLY, os.FileMode(mode))
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
