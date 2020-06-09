/*
   @Time : 2019-05-10 09:55
   @Author : frozenchen
   @File : config
   @Software: studio-library
*/
package redis

import (
	"github.com/freezeChen/studio-library/xtime"
)

type Config struct {
	Addr         string
	Auth         string
	Idle         int
	Active       int
	DialTimeout  xtime.Duration
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
	IdleTimeout  xtime.Duration
}
