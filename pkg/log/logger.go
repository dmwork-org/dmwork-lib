package log

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger
var errorLogger *zap.Logger
var warnLogger *zap.Logger
var testLogger *zap.Logger
var atom = zap.NewAtomicLevel()
var configOnce sync.Once

// ensureConfigured ensures the logger is configured exactly once.
// It uses sync.Once to prevent race conditions during concurrent access.
func ensureConfigured() {
	configOnce.Do(func() {
		doConfigure(NewOptions())
	})
}

// Configure configures the logger with the provided options.
// This function is safe for concurrent use - only the first call takes effect.
func Configure(opts *Options) {
	configOnce.Do(func() {
		doConfigure(opts)
	})
}

// doConfigure performs the actual logger configuration.
// This should only be called via sync.Once to ensure thread safety.
func doConfigure(opts *Options) {
	atom.SetLevel(opts.Level)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(newEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		atom,
	)
	if opts.LineNum {
		testLogger = zap.New(core, zap.AddCaller())
	} else {
		testLogger = zap.New(core)
	}

	infoWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(opts.LogDir, "info.log"),
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(newEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(infoWriter)),
		atom,
	)
	if opts.LineNum {
		logger = zap.New(core, zap.AddCaller())
	} else {
		logger = zap.New(core)
	}

	errorWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(opts.LogDir, "error.log"),
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(newEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(errorWriter)),
		zap.ErrorLevel,
	)
	if opts.LineNum {
		errorLogger = zap.New(core, zap.AddCaller())
	} else {
		errorLogger = zap.New(core)
	}

	warnWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(opts.LogDir, "warn.log"),
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(newEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(warnWriter)),
		zap.WarnLevel,
	)
	if opts.LineNum {
		warnLogger = zap.New(core, zap.AddCaller())
	} else {
		warnLogger = zap.New(core)
	}
}

// resetForTesting resets the logger state for testing purposes.
// This allows tests to reconfigure the logger with different options.
func resetForTesting() {
	configOnce = sync.Once{}
	logger = nil
	errorLogger = nil
	warnLogger = nil
	testLogger = nil
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "linenum",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeCaller:  zapcore.FullCallerEncoder,     // 全路径编码器
		EncodeName:    zapcore.FullNameEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}
}
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// Info Info
func Info(msg string, fields ...zap.Field) {
	ensureConfigured()
	logger.Info(msg, fields...)
}

// Debug Debug
func Debug(msg string, fields ...zap.Field) {
	ensureConfigured()
	logger.Debug(msg, fields...)
}

// Error Error
func Error(msg string, fields ...zap.Field) {
	ensureConfigured()
	errorLogger.Error(msg, fields...)
}

// Warn Warn
func Warn(msg string, fields ...zap.Field) {
	ensureConfigured()
	warnLogger.Warn(msg, fields...)
}

// Log Log
type Log interface {
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
}

// LIMLog TLog
type TLog struct {
	prefix string // 日志前缀
}

// NewLIMLog NewLIMLog
func NewTLog(prefix string) *TLog {

	return &TLog{prefix: prefix}
}

// Info Info
func (t *TLog) Info(msg string, fields ...zap.Field) {
	Info(fmt.Sprintf("【%s】%s", t.prefix, msg), fields...)
}

// Debug Debug
func (t *TLog) Debug(msg string, fields ...zap.Field) {
	Debug(fmt.Sprintf("【%s】%s", t.prefix, msg), fields...)
}

// Error Error
func (t *TLog) Error(msg string, fields ...zap.Field) {
	Error(fmt.Sprintf("【%s】%s", t.prefix, msg), fields...)
}

// Warn Warn
func (t *TLog) Warn(msg string, fields ...zap.Field) {
	Warn(fmt.Sprintf("【%s】%s", t.prefix, msg), fields...)
}
