/*
   @Time : 2019-05-16 10:45
   @Author : frozenchen
   @File : mysql
   @Software: studio
*/
package mysql

import (
	"time"

	"xorm.io/core"
	"xorm.io/xorm"
)

var (
	engine *xorm.Engine
	err    error
)

func New(c *Config) xorm.EngineInterface {
	engine, err = xorm.NewEngine("mysql", c.Source)
	if err != nil {
		panic(err)
	}
	engine.TZLocation = time.Local
	engine.SetMaxOpenConns(c.Active)
	engine.SetMaxIdleConns(c.Idle)
	engine.SetMapper(core.SameMapper{})
	engine.ShowSQL(c.Show)
	return engine
}
