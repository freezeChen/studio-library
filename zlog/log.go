package zlog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
	"strings"
	"time"
)

const TIMEFORMAT = "2006-01-02 15:04:05"

type Config struct {
	Name       string
	Debug      bool
	KafkaAddr  string
	WriteKafka bool
}

var gLogger *zap.Logger

func InitLogger(c *Config) {
	var (
		writers []zapcore.WriteSyncer
		//path    = getCurrentDirectory()
	)

	//writeSync := zapcore.AddSync(&Logger{
	//	Filename:   path + "/log/log.txt",
	//	MaxSize:    10, // 单文件容量上限(MB)
	//	MaxBackups: 30, //
	//	MaxAge:     30, // 文件保存天数
	//	LocalTime:  true,
	//})

	if c.Debug {
		writers = append(writers, os.Stdout)
	}

	if c.WriteKafka {
		writeSync, err := initKafkaWriter(c)
		if err != nil {
			panic("Failed to connect kafka:"+err.Error())
		} else {
			writers = append(writers, writeSync)
		}
	}

	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "app",
		CallerKey:      "line",
		MessageKey:     "message",
		StacktraceKey:  "error",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	atomicLevel := zap.NewAtomicLevel()
	syncer := zapcore.NewMultiWriteSyncer(writers...)
	core := zapcore.NewCore(encoder, syncer, atomicLevel)

	if c.Debug {
		atomicLevel.SetLevel(zap.DebugLevel)
	} else {
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	gLogger = zap.New(core)

}

func Infof(format string, args ...interface{}) {
	defer gLogger.Sync()
	pc, _, _, _ := runtime.Caller(1)
	forPC := runtime.FuncForPC(pc)
	split := strings.Split(forPC.Name(), ".")

	gLogger.WithOptions(zap.AddCallerSkip(1)).
		Info(fmt.Sprintf(format, args...),
			zap.String("class", split[len(split)-2]),
			zap.String("method", split[len(split)-1]),
		)
}

func Errorf(format string, args ...interface{}) {
	defer gLogger.Sync()
	pc, _, _, _ := runtime.Caller(1)
	forPC := runtime.FuncForPC(pc)
	split := strings.Split(forPC.Name(), ".")

	gLogger.WithOptions(zap.AddCallerSkip(1)).
		Error(fmt.Sprintf(format, args...),
			zap.String("class", split[len(split)-2]),
			zap.String("method", split[len(split)-1]),
		)
}

//日志时间格式化
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(TIMEFORMAT))
}
