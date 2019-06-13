/*
   @Time : 2019-05-15 17:34
   @Author : frozenchen
   @File : conf.go
   @Software: studio
*/
package conf

import (
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-micro/config/source/file"
	"os"
)

func LoadFileSource(path string) source.Source {
	_, err := os.Stat(path)
	if err != nil {
		return nil
	}
	fileSource := file.NewSource(file.WithPath(path))
	return fileSource
}
