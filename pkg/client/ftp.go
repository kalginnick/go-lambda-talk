package client

import (
	"io"
	"path"
	"time"

	"github.com/jlaffaye/ftp"
)

func WriteFTP(user, pswd, url, file string, data io.Reader) error {
	cn, err := ftp.DialTimeout(url, 10*time.Second)
	if err != nil {
		return err
	}
	defer cn.Quit()

	err = cn.Login(user, pswd)
	if err != nil {
		return err
	}

	return cn.Stor(path.Join(url, file), data)
}
