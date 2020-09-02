/*
   @Time : 2019-05-16 10:45
   @Author : frozenchen
   @File : mysql
   @Software: studio
*/
package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

func New(c *Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Source), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Use(prometheus.New(prometheus.Config{
		DBName:           "db1",
		RefreshInterval:  0,
		PushAddr:         "",
		StartServer:      false,
		HTTPServerPort:   0,
		MetricsCollector: nil,
	}))

	return db
}
