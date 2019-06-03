package zlog

import "testing"

func init() {
	InitLogger(&Config{
		Debug: true,
		Name:  "names",
	})
}

func TestInfof(t *testing.T) {
	Infof("print info (%s)", "test info")
}

func TestErrorf(t *testing.T) {
	Errorf("print error (%s)", "test error")
}
