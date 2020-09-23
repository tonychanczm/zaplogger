package zaplogger

import (
	"errors"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"sync"
)

var (
	resetLock sync.RWMutex
	defaultZapLogger *ZapLogger
)

type ZapLogger struct {
	sugarLogger *zap.SugaredLogger
	zapLogger *zap.Logger
	logWriter io.Writer
	config *LogConfig
}

func (z *ZapLogger) GetSugarLogger() *zap.SugaredLogger {
	return z.sugarLogger
}

func (z *ZapLogger) GetZapLogger() *zap.Logger {
	return z.zapLogger
}

type LogConfig struct {
	FilenameOrIoWriter interface{}
	FileMaxSize int
	FileMaxBackup int
	FileMaxAge int
	FileCompress bool
}

func NewZapLogger(config *LogConfig, logLevel zapcore.LevelEnabler) (zl *ZapLogger) {
	if config.FileMaxAge <= 0 {
		config.FileMaxAge = 30
	}

	if config.FileMaxBackup <= 0 {
		config.FileMaxBackup = 3
	}

	if config.FileMaxSize <= 0 {
		config.FileMaxSize = 100
	}
	zl = new(ZapLogger)
	zl.config = config
	var err error
	zl.logWriter, err = getWriter(config)

	if err != nil {
		config.FilenameOrIoWriter = os.Stdout
		zl.logWriter, err = getWriter(config)
		if err != nil {
			panic("Fail to enable default logWriter: " + err.Error())
		}
	}

	encoder := getEncoder()
	core := zapcore.NewCore(encoder, zapcore.AddSync(zl.logWriter), logLevel)

	zl.zapLogger = zap.New(core, zap.AddCaller(),zap.AddCallerSkip(1))
	zl.sugarLogger = zl.zapLogger.Sugar()
	return zl
}

func InitDefaultLogger(config *LogConfig, logLevel zapcore.LevelEnabler) {
	resetLock.Lock()
	defaultZapLogger = NewZapLogger(config, logLevel)
	resetLock.Unlock()
}

func init() {
	InitDefaultLogger(&LogConfig{}, zap.DebugLevel)
}

func getWriter(config *LogConfig) (io.Writer, error) {
	switch config.FilenameOrIoWriter.(type) {
	case string:
		return &lumberjack.Logger{
			Filename:   config.FilenameOrIoWriter.(string),
			MaxSize:    config.FileMaxSize,
			MaxBackups: config.FileMaxBackup,
			MaxAge:     config.FileMaxAge,
			Compress:   config.FileCompress,
		}, nil
	case io.Writer:
		return config.FilenameOrIoWriter.(io.Writer), nil
	}
	return nil, errors.New("string / io.Writer only")
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func (z *ZapLogger) Debug(args ...interface{}) {
	z.sugarLogger.Debug(args...)
}

func (z *ZapLogger) Debugf(template string, args ...interface{}) {
	z.sugarLogger.Debugf(template, args...)
}

func (z *ZapLogger) Info(args ...interface{}) {
	z.sugarLogger.Info(args...)
}

func (z *ZapLogger) Infof(template string, args ...interface{}) {
	z.sugarLogger.Infof(template, args...)
}

func (z *ZapLogger) Warn(args ...interface{}) {
	z.sugarLogger.Warn(args...)
}

func (z *ZapLogger) Warnf(template string, args ...interface{}) {
	z.sugarLogger.Warnf(template, args...)
}

func (z *ZapLogger) Fatal(args ...interface{}) {
	z.sugarLogger.Fatal(args...)
}

func (z *ZapLogger) Fatalf(template string, args ...interface{}) {
	z.sugarLogger.Fatalf(template, args...)
}

func (z *ZapLogger) Error(args ...interface{}) {
	z.sugarLogger.Error(args...)
}

func (z *ZapLogger) Errorf(template string, args ...interface{}) {
	z.sugarLogger.Errorf(template, args...)
}

func (z *ZapLogger) DPanic(args ...interface{}) {
	z.sugarLogger.DPanic(args...)
}

func (z *ZapLogger) DPanicf(template string, args ...interface{}) {
	z.sugarLogger.DPanicf(template, args...)
}

func Debug(args ...interface{}) {
	resetLock.RLock()
	defaultZapLogger.sugarLogger.Debug(args...)
	resetLock.RUnlock()
}

func Debugf(template string, args ...interface{}) {
	resetLock.RLock()
	defaultZapLogger.sugarLogger.Debugf(template, args...)
	resetLock.RUnlock()
}

func Info(args ...interface{}) {
	resetLock.RLock()
	defaultZapLogger.sugarLogger.Info(args...)
	resetLock.RUnlock()
}

func Infof(template string, args ...interface{}) {
	resetLock.RLock()
	defaultZapLogger.sugarLogger.Infof(template, args...)
	resetLock.RUnlock()
}

func Warn(args ...interface{}) {
	resetLock.RLock()
	defaultZapLogger.sugarLogger.Warn(args...)
	resetLock.RUnlock()
}

func Warnf(template string, args ...interface{}) {
	resetLock.RLock()
	defaultZapLogger.sugarLogger.Warnf(template, args...)
	resetLock.RUnlock()
}

func Fatal(args ...interface{}) {
	resetLock.RLock()
	defaultZapLogger.sugarLogger.Fatal(args...)
	resetLock.RUnlock()
}

func Fatalf(template string, args ...interface{}) {
	resetLock.RLock()
	defaultZapLogger.sugarLogger.Fatalf(template, args...)
	resetLock.RUnlock()
}

func Error(args ...interface{}) {
	resetLock.RLock()
	defaultZapLogger.sugarLogger.Error(args...)
	resetLock.RUnlock()
}

func Errorf(template string, args ...interface{}) {
	resetLock.RLock()
	defaultZapLogger.sugarLogger.Errorf(template, args...)
	resetLock.RUnlock()
}

func DPanic(args ...interface{}) {
	resetLock.RLock()
	defaultZapLogger.sugarLogger.DPanic(args...)
	resetLock.RUnlock()
}

func DPanicf(template string, args ...interface{}) {
	resetLock.RLock()
	defaultZapLogger.sugarLogger.DPanicf(template, args...)
	resetLock.RUnlock()
}

func GetDefaultZapLogger() *zap.Logger {
	return defaultZapLogger.zapLogger
}

func GetDefaultSugarLogger() *zap.SugaredLogger {
	return defaultZapLogger.sugarLogger
}

type ZapLogWriter struct {
	level string
}

func NewZapLogWriter(level string) *ZapLogWriter {
	return &ZapLogWriter{level:level}
}

func (z *ZapLogWriter) Write(b []byte) (int, error) {
	resetLock.RLock()
	switch string(z.level) {
	case "debug", "DEBUG":
		defaultZapLogger.sugarLogger.Debug(string(b))
	case "info", "INFO", "": // make the zero value useful
		defaultZapLogger.sugarLogger.Info(string(b))
	case "warn", "WARN":
		defaultZapLogger.sugarLogger.Warn(string(b))
	case "error", "ERROR":
		defaultZapLogger.sugarLogger.Error(string(b))
	case "panic", "PANIC":
		defaultZapLogger.sugarLogger.Panic(string(b))
	case "fatal", "FATAL":
		defaultZapLogger.sugarLogger.Fatal(string(b))
	default:
		defaultZapLogger.sugarLogger.Info(string(b))
	}
	resetLock.RUnlock()
	return len(b), nil
}

type NullWriter struct {

}

func (n *NullWriter) Write(b []byte) (int, error) {
	return len(b), nil
}