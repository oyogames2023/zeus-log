package zap

import (
	"fmt"
	xlog "github.com/oyogames2023/zeus-log"
	"github.com/oyogames2023/zeus-log/rollwriter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strconv"
	"time"
)

// Some ZapCore constants.
const (
	ConsoleZapCore = "console"
	FileZapCore    = "file"
)

var (
	defaultConfig = []xlog.OutputConfig{
		{
			Writer:    "console",
			Level:     "debug",
			Formatter: "console",
		},
	}
	Levels = map[string]zapcore.Level{
		"":      zapcore.DebugLevel,
		"trace": zapcore.DebugLevel,
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"fatal": zapcore.FatalLevel,
		"panic": zapcore.PanicLevel,
	}
	levelToZapLevel = map[xlog.Level]zapcore.Level{
		xlog.LevelTrace: zapcore.DebugLevel,
		xlog.LevelDebug: zapcore.DebugLevel,
		xlog.LevelInfo:  zapcore.InfoLevel,
		xlog.LevelWarn:  zapcore.WarnLevel,
		xlog.LevelError: zapcore.ErrorLevel,
		xlog.LevelFatal: zapcore.FatalLevel,
		xlog.LevelPanic: zapcore.PanicLevel,
	}
	zapLevelToLevel = map[zapcore.Level]xlog.Level{
		zapcore.DebugLevel: xlog.LevelDebug,
		zapcore.InfoLevel:  xlog.LevelInfo,
		zapcore.WarnLevel:  xlog.LevelWarn,
		zapcore.ErrorLevel: xlog.LevelError,
		zapcore.FatalLevel: xlog.LevelFatal,
		zapcore.PanicLevel: xlog.LevelPanic,
	}
)

// NewZapLog creates a zap Logger object whose caller skip is set to 2.
func NewZapLog(c xlog.Config) xlog.Logger {
	return NewZapLogWithCallerSkip(c, 2)
}

// NewZapLogWithCallerSkip creates a default Logger from zap.
func NewZapLogWithCallerSkip(cfg xlog.Config, callerSkip int) xlog.Logger {
	var (
		cores  []zapcore.Core
		levels []zap.AtomicLevel
	)
	for _, c := range cfg {
		writer := xlog.GetWriter(c.Writer)
		if writer == nil {
			panic("log: writer core: " + c.Writer + " no registered")
		}
		decoder := &Decoder{OutputConfig: &c}
		if err := writer.Setup(c.Writer, decoder); err != nil {
			panic("log: writer core: " + c.Writer + " setup fail: " + err.Error())
		}
		cores = append(cores, decoder.Core)
		levels = append(levels, decoder.ZapLevel)
	}
	return &zapLog{
		levels: levels,
		logger: zap.New(
			zapcore.NewTee(cores...),
			zap.AddCallerSkip(callerSkip),
			zap.AddCaller(),
		),
	}
}

func newEncoder(c *xlog.OutputConfig) zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        GetLogEncoderKey("T", c.FormatConfig.TimeKey),
		LevelKey:       GetLogEncoderKey("L", c.FormatConfig.LevelKey),
		NameKey:        GetLogEncoderKey("N", c.FormatConfig.NameKey),
		CallerKey:      GetLogEncoderKey("C", c.FormatConfig.CallerKey),
		FunctionKey:    GetLogEncoderKey(zapcore.OmitKey, c.FormatConfig.FunctionKey),
		MessageKey:     GetLogEncoderKey("M", c.FormatConfig.MessageKey),
		StacktraceKey:  GetLogEncoderKey("S", c.FormatConfig.StacktraceKey),
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     NewTimeEncoder(c.FormatConfig.TimeFormat),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if c.EnableColor {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	switch c.Formatter {
	case "console":
		return zapcore.NewConsoleEncoder(encoderCfg)
	case "json":
		return zapcore.NewJSONEncoder(encoderCfg)
	default:
		return zapcore.NewConsoleEncoder(encoderCfg)
	}
}

// GetLogEncoderKey gets user defined log output name, uses defKey if empty.
func GetLogEncoderKey(defKey, key string) string {
	if key == "" {
		return defKey
	}
	return key
}

func newConsoleCore(c *xlog.OutputConfig) (zapcore.Core, zap.AtomicLevel) {
	lvl := zap.NewAtomicLevelAt(Levels[c.Level])
	return zapcore.NewCore(
		newEncoder(c),
		zapcore.Lock(os.Stdout),
		lvl), lvl
}

func newFileCore(c *xlog.OutputConfig) (zapcore.Core, zap.AtomicLevel, error) {
	fmt.Println(c.WriterConfig.MaxAge)
	opts := []rollwriter.Option{
		rollwriter.WithMaxAge(c.WriterConfig.MaxAge),
		rollwriter.WithMaxBackups(c.WriterConfig.MaxBackups),
		rollwriter.WithCompress(c.WriterConfig.Compress),
		rollwriter.WithMaxSize(c.WriterConfig.MaxSize),
	}
	// roll by time.
	if c.WriterConfig.RollType != xlog.RollingBySizeStr {
		opts = append(opts, rollwriter.WithRotationTime(c.WriterConfig.TimeUnit.Format()))
	}
	writer, err := rollwriter.NewRollWriter(c.WriterConfig.FileName, opts...)
	if err != nil {
		return nil, zap.AtomicLevel{}, err
	}

	// write mode.
	var ws zapcore.WriteSyncer
	switch m := xlog.GetWriteMode(c.WriterConfig.WriteMode); m {
	case 0, xlog.WriteFast:
		// Use WriteFast as default mode.
		// It has better performance, discards logs on full and avoid blocking service.
		ws = rollwriter.NewAsyncRollWriter(writer, rollwriter.WithDropLog(true))
	case xlog.WriteSync:
		ws = zapcore.AddSync(writer)
	case xlog.WriteAsync:
		ws = rollwriter.NewAsyncRollWriter(writer, rollwriter.WithDropLog(false))
	default:
		return nil, zap.AtomicLevel{}, fmt.Errorf("validating WriteMode parameter: got %d, "+
			"but expect one of WriteFast(%d), WriteAsync(%d), or WriteSync(%d)", m,
			xlog.WriteFast, xlog.WriteAsync, xlog.WriteSync)
	}

	// log level.
	lvl := zap.NewAtomicLevelAt(Levels[c.Level])
	return zapcore.NewCore(
		newEncoder(c),
		ws, lvl,
	), lvl, nil
}

// NewTimeEncoder creates a time format encoder.
func NewTimeEncoder(format string) zapcore.TimeEncoder {
	switch format {
	case "":
		return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendByteString(defaultTimeFormat(t))
		}
	case "seconds":
		return zapcore.EpochTimeEncoder
	case "milliseconds":
		return zapcore.EpochMillisTimeEncoder
	case "nanoseconds":
		return zapcore.EpochNanosTimeEncoder
	default:
		return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(format))
		}
	}
}

// defaultTimeFormat returns the default time format "2006-01-02 15:04:05.000",
// which performs better than https://pkg.go.dev/time#Time.AppendFormat.
func defaultTimeFormat(t time.Time) []byte {
	t = t.Local()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	micros := t.Nanosecond() / 1000

	buf := make([]byte, 23)
	buf[0] = byte((year/1000)%10) + '0'
	buf[1] = byte((year/100)%10) + '0'
	buf[2] = byte((year/10)%10) + '0'
	buf[3] = byte(year%10) + '0'
	buf[4] = '-'
	buf[5] = byte((month)/10) + '0'
	buf[6] = byte((month)%10) + '0'
	buf[7] = '-'
	buf[8] = byte((day)/10) + '0'
	buf[9] = byte((day)%10) + '0'
	buf[10] = ' '
	buf[11] = byte((hour)/10) + '0'
	buf[12] = byte((hour)%10) + '0'
	buf[13] = ':'
	buf[14] = byte((minute)/10) + '0'
	buf[15] = byte((minute)%10) + '0'
	buf[16] = ':'
	buf[17] = byte((second)/10) + '0'
	buf[18] = byte((second)%10) + '0'
	buf[19] = '.'
	buf[20] = byte((micros/100000)%10) + '0'
	buf[21] = byte((micros/10000)%10) + '0'
	buf[22] = byte((micros/1000)%10) + '0'
	return buf
}

// zapLog is a Logger implementation based on zaplogger.
type zapLog struct {
	levels []zap.AtomicLevel
	logger *zap.Logger
}

func (l *zapLog) WithOptions(opts ...xlog.Option) xlog.Logger {
	o := &xlog.Options{}
	for _, opt := range opts {
		opt(o)
	}
	return &zapLog{
		levels: l.levels,
		logger: l.logger.WithOptions(zap.AddCallerSkip(o.Skip)),
	}
}

// With returns a new logger with key/value paris.
func (l *zapLog) With(args ...any) xlog.Logger {
	return nil
}

// WithFields returns a new logger with key/value paris.
func (l *zapLog) WithFields(fields ...xlog.Field) xlog.Logger {
	zapFields := make([]zap.Field, len(fields))
	for i := range fields {
		zapFields[i] = zap.Any(fields[i].Key, fields[i].Value)
	}

	return &zapLog{
		levels: l.levels,
		logger: l.logger.With(zapFields...)}
}

func getLogMsg(args ...interface{}) string {
	msg := fmt.Sprint(args...)
	return msg
}

func getLogMsgf(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return msg
}

// Trace logs to TRACE log. Arguments are handled in the manner of fmt.Print.
func (l *zapLog) Trace(args ...any) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(getLogMsg(args...))
	}
}

// Tracef logs to TRACE log. Arguments are handled in the manner of fmt.Printf.
func (l *zapLog) Tracef(format string, args ...any) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(getLogMsgf(format, args...))
	}
}

// Traceln logs to TRACE log. Arguments are handled in the manner of fmt.Println.
func (l *zapLog) Traceln(args ...any) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(fmt.Sprintf("%s\n", getLogMsg(args...)))
	}
}

// Debug logs to DEBUG log. Arguments are handled in the manner of fmt.Print.
func (l *zapLog) Debug(args ...any) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(getLogMsg(args...))
	}
}

// Debugf logs to DEBUG log. Arguments are handled in the manner of fmt.Printf.
func (l *zapLog) Debugf(format string, args ...any) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(getLogMsgf(format, args...))
	}
}

// Debugln logs to DEBUG log. Arguments are handled in the manner of fmt.Println.
func (l *zapLog) Debugln(args ...any) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(fmt.Sprintf("%s\n", getLogMsg(args...)))
	}
}

// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
func (l *zapLog) Info(args ...any) {
	if l.logger.Core().Enabled(zapcore.InfoLevel) {
		l.logger.Info(getLogMsg(args...))
	}
}

// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func (l *zapLog) Infof(format string, args ...any) {
	if l.logger.Core().Enabled(zapcore.InfoLevel) {
		l.logger.Info(getLogMsgf(format, args...))
	}
}

// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (l *zapLog) Infoln(args ...any) {
	if l.logger.Core().Enabled(zapcore.InfoLevel) {
		l.logger.Info(fmt.Sprintf("%s\n", getLogMsg(args...)))
	}
}

// Warn logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func (l *zapLog) Warn(args ...any) {
	if l.logger.Core().Enabled(zapcore.WarnLevel) {
		l.logger.Warn(getLogMsg(args...))
	}
}

// Warnf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func (l *zapLog) Warnf(format string, args ...any) {
	if l.logger.Core().Enabled(zapcore.WarnLevel) {
		l.logger.Warn(getLogMsgf(format, args...))
	}
}

// Warnln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (l *zapLog) Warnln(args ...any) {
	if l.logger.Core().Enabled(zapcore.WarnLevel) {
		l.logger.Warn(fmt.Sprintf("%s\n", getLogMsg(args...)))
	}
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func (l *zapLog) Error(args ...any) {
	if l.logger.Core().Enabled(zapcore.ErrorLevel) {
		l.logger.Error(getLogMsg(args...))
	}
}

// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func (l *zapLog) Errorf(format string, args ...any) {
	if l.logger.Core().Enabled(zapcore.ErrorLevel) {
		l.logger.Error(getLogMsgf(format, args...))
	}
}

// Errorln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (l *zapLog) Errorln(args ...any) {
	if l.logger.Core().Enabled(zapcore.ErrorLevel) {
		l.logger.Error(fmt.Sprintf("%s\n", getLogMsg(args...)))
	}
}

// Fatal logs to FATAL log. Arguments are handled in the manner of fmt.Print.
func (l *zapLog) Fatal(args ...any) {
	if l.logger.Core().Enabled(zapcore.FatalLevel) {
		l.logger.Fatal(getLogMsg(args...))
	}
}

// Fatalf logs to FATAL log. Arguments are handled in the manner of fmt.Printf.
func (l *zapLog) Fatalf(format string, args ...any) {
	if l.logger.Core().Enabled(zapcore.FatalLevel) {
		l.logger.Fatal(getLogMsgf(format, args...))
	}
}

// Fatalln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (l *zapLog) Fatalln(args ...any) {
	if l.logger.Core().Enabled(zapcore.FatalLevel) {
		l.logger.Fatal(fmt.Sprintf("%s\n", getLogMsg(args...)))
	}
}

// Panic logs to FATAL log. Arguments are handled in the manner of fmt.Print.
func (l *zapLog) Panic(args ...any) {
	if l.logger.Core().Enabled(zapcore.PanicLevel) {
		l.logger.Panic(getLogMsg(args...))
	}
}

// Panicf logs to FATAL log. Arguments are handled in the manner of fmt.Printf.
func (l *zapLog) Panicf(format string, args ...any) {
	if l.logger.Core().Enabled(zapcore.PanicLevel) {
		l.logger.Panic(getLogMsgf(format, args...))
	}
}

// Panicln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (l *zapLog) Panicln(args ...any) {
	if l.logger.Core().Enabled(zapcore.PanicLevel) {
		l.logger.Panic(fmt.Sprintf("%s\n", getLogMsg(args...)))
	}
}

// Sync calls the zap logger's Sync method, and flushes any buffered log entries.
// Applications should take care to call Sync before exiting.
func (l *zapLog) Sync() error {
	return l.logger.Sync()
}

// SetLevel sets output log level.
func (l *zapLog) SetLevel(output string, level xlog.Level) {
	i, e := strconv.Atoi(output)
	if e != nil {
		return
	}
	if i < 0 || i >= len(l.levels) {
		return
	}
	l.levels[i].SetLevel(levelToZapLevel[level])
}

// GetLevel gets output log level.
func (l *zapLog) GetLevel(output string) xlog.Level {
	i, e := strconv.Atoi(output)
	if e != nil {
		return xlog.LevelDebug
	}
	if i < 0 || i >= len(l.levels) {
		return xlog.LevelDebug
	}
	return zapLevelToLevel[l.levels[i].Level()]
}

// CustomTimeFormat customize time format.
// Deprecated: Use https://pkg.go.dev/time#Time.Format instead.
func CustomTimeFormat(t time.Time, format string) string {
	return t.Format(format)
}

// DefaultTimeFormat returns the default time format "2006-01-02 15:04:05.000".
// Deprecated: Use https://pkg.go.dev/time#Time.AppendFormat instead.
func DefaultTimeFormat(t time.Time) []byte {
	return defaultTimeFormat(t)
}

// RedirectStdLog redirects std log to trpc logger as log level INFO.
// After redirection, log flag is zero, the prefix is empty.
// The returned function may be used to recover log flag and prefix, and redirect output to
// os.Stderr.
func RedirectStdLog(logger xlog.Logger) (func(), error) {
	return RedirectStdLogAt(logger, zap.InfoLevel)
}

// RedirectStdLogAt redirects std log to trpc logger with a specific level.
// After redirection, log flag is zero, the prefix is empty.
// The returned function may be used to recover log flag and prefix, and redirect output to
// os.Stderr.
func RedirectStdLogAt(logger xlog.Logger, level zapcore.Level) (func(), error) {
	if l, ok := logger.(*zapLog); ok {
		return zap.RedirectStdLogAt(l.logger, level)
	}

	return nil, fmt.Errorf("log: only supports redirecting std logs to trpc zap logger")
}
