/*
   @Time : 2019-05-15 17:34
   @Author : frozenchen
   @File : conf.go
   @Software: studio
*/
package conf

import (
	"os"

	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/config/source/file"
)

func LoadFileSource(path string) source.Source {
	_, err := os.Stat(path)
	if err != nil {
		return nil
	}
	fileSource := file.NewSource(file.WithPath(path))
	return fileSource
}
