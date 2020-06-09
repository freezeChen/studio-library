package zlog

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type KV zapcore.Field

func KVString(k, v string) KV {
	return KV(zap.String(k, v))
}

func KVInt(k string, v int) KV {
	return KV(zap.Int(k, v))
}

func KVInt64(k string, v int64) KV {
	return KV(zap.Int64(k, v))
}

func KVFloat32(k string, v float32) KV {
	return KV(zap.Float32(k, v))
}

func KVFloat64(k string, v float64) KV {
	return KV(zap.Float64(k, v))
}

func KVBool(k string, v bool) KV {
	return KV(zap.Bool(k, v))
}

func KVTime(k string, v time.Time) KV {
	return KV(zap.Time(k, v))
}

func Kv(k string, v interface{}) KV {
	return KV(zap.Any(k, v))
}

func kv2Field(kv []KV) []zapcore.Field {
	var fields []zapcore.Field
	for _, k := range kv {
		fields = append(fields, zapcore.Field(k))
	}
	return fields
}
