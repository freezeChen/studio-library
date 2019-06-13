/*
   @Time : 2019-05-28 15:33
   @Author : frozenchen
   @File : error
   @Software: studio
*/
package redis

import "errors"

var (
	ErrKey    = errors.New("redis: key is empty")
	ErrNotFound = errors.New("not found")
)
