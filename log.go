package log

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	InfoLevel   Level = zap.InfoLevel   // 0, default level
	WarnLevel   Level = zap.WarnLevel   // 1
	ErrorLevel  Level = zap.ErrorLevel  // 2
	DPanicLevel Level = zap.DPanicLevel // 3, used in development log
	// PanicLevel logs a message, then panics
	PanicLevel Level = zap.PanicLevel // 4
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel Level = zap.FatalLevel // 5
	DebugLevel Level = zap.DebugLevel // -1
)

type Field = zap.Field

type Logger struct {
	l     *zap.Logger // zap ensure that zap.Logger is safe for concurrent use
	level zap.AtomicLevel
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Debugf(msg string, args ...any) {
	l.l.Sugar().Debugf(msg, args)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Infof(msg string, args ...any) {
	l.l.Sugar().Infof(msg, args...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Warnf(msg string, args ...any) {
	l.l.Sugar().Warnf(msg, args...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}

func (l *Logger) Errorf(msg string, args ...any) {
	l.l.Sugar().Errorf(msg, args...)
}

func (l *Logger) DPanic(msg string, fields ...Field) {
	l.l.DPanic(msg, fields...)
}

func (l *Logger) DPanicf(msg string, args ...any) {
	l.l.Sugar().DPanicf(msg, args...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}

func (l *Logger) Panicf(msg string, args ...any) {
	l.l.Sugar().Panicf(msg, args...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Fatalf(msg string, args ...any) {
	l.l.Sugar().Fatalf(msg, args...)
}

var (
	std = New(os.Stderr, InfoLevel, WithCaller(true), AddCallerSkip(1))

	Info   = std.Info
	Warn   = std.Warn
	Error  = std.Error
	DPanic = std.DPanic
	Panic  = std.Panic
	Fatal  = std.Fatal
	Debug  = std.Debug
)

// not safe for concurrent use, replace default std
func ResetDefault(l *Logger) {
	std = l
	Info = std.Info
	Warn = std.Warn
	Error = std.Error
	DPanic = std.DPanic
	Panic = std.Panic
	Fatal = std.Fatal
	Debug = std.Debug
}

func Default() *Logger { return std }

type Option = zap.Option

var (
	WithCaller    = zap.WithCaller
	AddCallerSkip = zap.AddCallerSkip
	AddStacktrace = zap.AddStacktrace
)

// New create a new logger
func New(w io.Writer, level Level, opts ...Option) *Logger {
	if w == nil {
		panic("the writer is nil")
	}

	cfg := zap.NewProductionConfig()
	// set time format
	cfg.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format(time.RFC3339Nano))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(w),
		zapcore.Level(level),
	)

	return &Logger{
		l:     zap.New(core, opts...),
		level: zap.NewAtomicLevelAt(level),
	}
}

// SetLevel alters the logging level on runtime
// it is concurrent-safe
func (l *Logger) SetLevel(level Level) {
	l.level.SetLevel(zapcore.Level(level))
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func Sync() error {
	if std != nil {
		return std.Sync()
	}

	return nil
}
