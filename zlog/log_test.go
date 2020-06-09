package zlog

import (
	"testing"
)

func init() {
	InitLogger(&Config{
		Debug:      false,
		Name:       "names",
		WriteKafka: false,
	})
}

func TestInfof(t *testing.T) {

	Debug("debug消息")
	Info("info消息")
	Error("error消息")
	//dd()

}

func TestErrorf(t *testing.T) {

	for i := 0; i < 1000; i++ {
		Info("sdf")
	}

}

func dd() {

}
