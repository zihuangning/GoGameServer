package logger

import (
	"os"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	fileLogger    *zap.Logger
	consoleLogger *zap.Logger
	conf          option
	logLimiter    *rateLimiter
)

type option struct {
	debug bool
	both  bool
	file  bool
	name  string
}

// 默认参数
func defaultOption() option {
	return option{
		debug: false,
		both:  true,
		file:  true,
		name:  "log",
	}
}

// Option 参数
type Option func(*option)

// Init Init
func Init(opts ...Option) {
	conf := defaultOption()
	for _, opt := range opts {
		opt(&conf)
	}

	var logLevel = zapcore.InfoLevel
	if conf.debug {
		logLevel = zapcore.DebugLevel
	}

	// 初始化日志限流器
	logLimiter = newRateLimiter(100, time.Second) // 每秒最多100条日志

	if conf.both {
		createConsoleLogger(conf.name, logLevel)
	}

	if conf.file {
		createFileLogger(conf.name, logLevel)
	}
}

func createFileLogger(name string, logLevel zapcore.Level) {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/" + name + ".log",
		MaxSize:    5,    // 每个日志文件最大5MB
		MaxBackups: 3,    // 最多保留3个备份
		MaxAge:     7,    // 日志文件最多保留7天
		Compress:   true, // 压缩旧日志
	}
	writeSyncer := zapcore.AddSync(lumberJackLogger)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("15:04:05.000"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "msg"
	encoderConfig.NameKey = "name"
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	fileLogger = zap.New(zapcore.NewCore(encoder, writeSyncer, logLevel), zap.AddCaller(), zap.AddCallerSkip(1))
	fileLogger = fileLogger.Named(name)
}

func createConsoleLogger(fileName string, logLevel zapcore.Level) {
	writeSyncer := os.Stderr

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("15:04:05.000"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "msg"
	encoderConfig.NameKey = "name"
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	consoleLogger = zap.New(zapcore.NewCore(encoder, writeSyncer, logLevel), zap.AddCaller(), zap.AddCallerSkip(1))
	consoleLogger = consoleLogger.Named(fileName)
}

// WithDebug 是否开启Debug参数
func WithDebug(debug bool) Option {
	return func(o *option) {
		o.debug = debug
	}
}

// WithBoth 是否开启Both参数
func WithBoth(both bool) Option {
	return func(o *option) {
		o.both = both
	}
}

// WithFile 是否开启File参数
func WithFile(file bool) Option {
	return func(o *option) {
		o.file = file
	}
}

// WithName 设置name参数
func WithName(name string) Option {
	return func(o *option) {
		o.name = name
	}
}

// rateLimiter 日志限流器
type rateLimiter struct {
	rate      int
	interval  time.Duration
	lastReset time.Time
	count     int
	mutex     sync.Mutex
}

func newRateLimiter(rate int, interval time.Duration) *rateLimiter {
	return &rateLimiter{
		rate:      rate,
		interval:  interval,
		lastReset: time.Now(),
	}
}

func (r *rateLimiter) allow() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	if now.Sub(r.lastReset) >= r.interval {
		r.count = 0
		r.lastReset = now
	}

	if r.count >= r.rate {
		return false
	}

	r.count++
	return true
}

// Error Error
func Error(msg string, fields ...zap.Field) {
	if !logLimiter.allow() {
		return
	}

	if consoleLogger != nil {
		consoleLogger.Error(msg, fields...)
	}

	if fileLogger != nil {
		fileLogger.Error(msg, fields...)
	}
}

// Warn Warn
func Warn(msg string, fields ...zap.Field) {
	if !logLimiter.allow() {
		return
	}

	if consoleLogger != nil {
		consoleLogger.Warn(msg, fields...)
	}

	if fileLogger != nil {
		fileLogger.Warn(msg, fields...)
	}
}

// Info Info
func Info(msg string, fields ...zap.Field) {
	if !logLimiter.allow() {
		return
	}

	if consoleLogger != nil {
		consoleLogger.Info(msg, fields...)
	}

	if fileLogger != nil {
		fileLogger.Info(msg, fields...)
	}
}

// Debug Debug
func Debug(msg string, fields ...zap.Field) {
	if !logLimiter.allow() {
		return
	}

	if consoleLogger != nil {
		consoleLogger.Debug(msg, fields...)
	}

	if fileLogger != nil {
		fileLogger.Debug(msg, fields...)
	}
}
