/*
   @Time : 2019-05-10 09:55
   @Author : frozenchen
   @File : config
   @Software: studio-library
*/
package redis

import (
	"time"
)


type Config struct {
	Addr         string
	Auth         string
	Idle         int
	Active       int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}
