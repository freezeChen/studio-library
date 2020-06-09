package zlog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/freezeChen/studio-library/metadata"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const TIMEFORMAT = "2006-01-02 15:04:05"

type Config struct {
	Name       string
	Debug      bool
	KafkaAddr  string
	WriteKafka bool

	LogFileDir    string `json:"logFileDir"`
	ErrorFileName string `json:"errorFileName"`
	WarnFileName  string `json:"warnFileName"`
	InfoFileName  string `json:"infoFileName"`
	DebugFileName string `json:"debugFileName"`
	MaxSize       int    `json:"maxSize"` // megabytes
	MaxBackups    int    `json:"maxBackups"`
	MaxAge        int    `json:"maxAge"` // days
}

type Logger struct {
	*zap.Logger
	cfg       *Config
	zapConfig zap.Config
}

var (
	logger                         *Logger
	sp                             = string(filepath.Separator)
	errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer       // IO输出
	debugConsoleWS                 = zapcore.Lock(os.Stdout) // 控制台标准输出
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
	kafkaWS                        zapcore.WriteSyncer
)

func InitLogger(c *Config) {
	logger = &Logger{
		cfg: c,
	}

	logger.loadCfg()
	logger.setSyncers()

	logger.zapConfig.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "xtime",
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
	}
	var err error
	logger.Logger, err = logger.zapConfig.Build(logger.cores(), zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		panic(err)
	}
}

func (l *Logger) loadCfg() {
	if l.cfg.Debug {
		l.zapConfig = zap.NewDevelopmentConfig()
	} else {
		l.zapConfig = zap.NewProductionConfig()
	}

	// 默认输出到程序运行目录的logs子目录
	if l.cfg.LogFileDir == "" {
		l.cfg.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
		l.cfg.LogFileDir += sp + "logs" + sp
	}

	if l.cfg.Name == "" {
		l.cfg.Name = "app"
	}

	if l.cfg.ErrorFileName == "" {
		l.cfg.ErrorFileName = "error.log"
	}

	if l.cfg.WarnFileName == "" {
		l.cfg.WarnFileName = "warn.log"
	}

	if l.cfg.InfoFileName == "" {
		l.cfg.InfoFileName = "info.log"
	}

	if l.cfg.DebugFileName == "" {
		l.cfg.DebugFileName = "debug.log"
	}

	if l.cfg.MaxSize == 0 {
		l.cfg.MaxSize = 50
	}
	if l.cfg.MaxBackups == 0 {
		l.cfg.MaxBackups = 3
	}
	if l.cfg.MaxAge == 0 {
		l.cfg.MaxAge = 30
	}
}

func (l *Logger) setSyncers() {
	var err error
	if l.cfg.WriteKafka {
		kafkaWS, err = initKafkaWriter(l.cfg)
		if err != nil {
			panic("Failed to connect kafka:" + err.Error())
		}
		errWS = kafkaWS
		warnWS = kafkaWS
		infoWS = kafkaWS
		debugWS = kafkaWS
		return
	}

	f := func(fN string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   logger.cfg.LogFileDir + sp + logger.cfg.Name + "-" + fN,
			MaxSize:    logger.cfg.MaxSize,
			MaxBackups: logger.cfg.MaxBackups,
			MaxAge:     logger.cfg.MaxAge,
			Compress:   true,
			LocalTime:  true,
		})
	}
	errWS = f(l.cfg.ErrorFileName)
	warnWS = f(l.cfg.WarnFileName)
	infoWS = f(l.cfg.InfoFileName)
	debugWS = f(l.cfg.DebugFileName)
}

func (l *Logger) cores() zap.Option {

	fileEncoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(l.zapConfig.EncoderConfig)

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.zapConfig.Level.Level() > -1
	})

	cores := []zapcore.Core{
		// region 日志文件

		// error 及以上
		zapcore.NewCore(fileEncoder, errWS, errPriority),

		// warn
		zapcore.NewCore(fileEncoder, warnWS, warnPriority),

		// info
		zapcore.NewCore(fileEncoder, infoWS, infoPriority),

		// debug
		zapcore.NewCore(fileEncoder, debugWS, debugPriority),

		// endregion

		// region 控制台

		// 错误及以上
		zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),

		// 警告
		zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),

		// info
		zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),

		// debug
		zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority),

		// endregion
	}

	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

func Debug(format string, args ...interface{}) {
	class, method := funName(2)
	logger.WithOptions(zap.AddCallerSkip(1)).
		Debug(fmt.Sprintf(format, args...),
			zap.String("class", class),
			zap.String("method", method),
		)
}
func Info(format string, args ...interface{}) {
	class, method := funName(2)
	logger.WithOptions(zap.AddCallerSkip(1)).
		Info(
			fmt.Sprintf(format, args...),
			zap.String("class", class),
			zap.String("method", method),
		)
}
func Error(format string, args ...interface{}) {
	class, method := funName(2)
	logger.WithOptions(zap.AddCallerSkip(1)).
		Error(fmt.Sprintf(format, args...),
			zap.String("class", class),
			zap.String("method", method),
		)
}
func Warn(format string, args ...interface{}) {
	class, method := funName(2)
	logger.WithOptions(zap.AddCallerSkip(1)).
		Warn(fmt.Sprintf(format, args...),
			zap.String("class", class),
			zap.String("method", method),
		)
}

func DebugV(ctx context.Context, args ...KV) {
	if traceId, ok := ctx.Value(metadata.GinTraceId).(string); ok {
		args = append(args, KVString(metadata.GinTraceId, traceId))
	}
	name, method := funName(2)
	args = append(args, KVString("class", name), KVString("method", method))

	fields := kv2Field(args)
	logger.WithOptions(zap.AddCallerSkip(1)).
		Debug("log",
			fields...)

}
func InfoV(ctx context.Context, args ...KV) {
	if traceId, ok := ctx.Value(metadata.GinTraceId).(string); ok {
		args = append(args, KVString(metadata.GinTraceId, traceId))
	}
	name, method := funName(2)
	args = append(args, KVString("class", name), KVString("method", method))

	fields := kv2Field(args)
	logger.WithOptions(zap.AddCallerSkip(1)).
		Info("log",
			fields...)

}
func WarnV(ctx context.Context, args ...KV) {
	if traceId, ok := ctx.Value(metadata.GinTraceId).(string); ok {
		args = append(args, KVString(metadata.GinTraceId, traceId))
	}
	name, method := funName(2)
	args = append(args, KVString("class", name), KVString("method", method))

	fields := kv2Field(args)
	logger.WithOptions(zap.AddCallerSkip(1)).
		Warn("log",
			fields...)

}
func ErrorV(ctx context.Context, args ...KV) {
	if traceId, ok := ctx.Value(metadata.GinTraceId).(string); ok {
		args = append(args, KVString(metadata.GinTraceId, traceId))
	}
	name, method := funName(2)
	args = append(args, KVString("class", name), KVString("method", method))
	fields := kv2Field(args)
	logger.WithOptions(zap.AddCallerSkip(1)).
		Error("log",
			fields...)

}

func Sync() error {
	return logger.Sync()
}

//日志时间格式化
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(TIMEFORMAT))
}

func funName(skip int) (class, method string) {
	caller, _, _, _ := runtime.Caller(skip)
	forPC := runtime.FuncForPC(caller)
	split := strings.Split(forPC.Name(), ".")
	class = split[len(split)-2]
	method = split[len(split)-1]
	return
}
