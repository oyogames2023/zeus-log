package zeus_log

import (
	yaml "gopkg.in/yaml.v3"
	"time"
)

const (
	OutputConsole = "console"
	OutputFile    = "file"
)

// Config is the log config. Each log may have multiple outputs.
type Config []OutputConfig

type OutputConfig struct {
	// Writer is the output of log, includes console, file and remote.
	Writer       string       `yaml:"writer"`
	WriterConfig WriterConfig `yaml:"writer_config"`

	// Formatter is the format of log, such as console or json.
	Formatter    string       `yaml:"formatter"`
	FormatConfig FormatConfig `yaml:"format_config"`

	// RemoteConfig is the remote config. It's defined by business and should be
	// registered by third-party modules.
	RemoteConfig yaml.Node `yaml:"remote_config"`

	// Level controls the log level, like debug, info or error, etc...
	Level string `yaml:"level"`

	// CallerSkip controls the nesting depth of log function.
	CallerSkip int `yaml:"caller_skip"`

	// EnableColor determines if the output is colored. The default value is false.
	EnableColor bool `yaml:"enable_color"`
}

// WriterConfig is the local file config.
type WriterConfig struct {
	// LogPath is the log path like "/usr/local/logs" or
	// "C:\Users<YourUsername>\AppData\Local\Temp".
	LogPath string `yaml:"log_path"`
	// FileName is the file name like "app.log".
	FileName string `yaml:"file_name"`
	// WriteMode is the log write mod. sync, async, fast.(default as fast)
	WriteMode string `yaml:"write_mode"`
	// RollType is the log rolling type. Split files by size/time.(default as time)
	RollType string `yaml:"roll_type"`
	// MaxAge is the max expire times(day).
	MaxAge int `yaml:"max_age"`
	// MaxBackups is the max backup files.
	MaxBackups int `yaml:"max_backups"`
	// Compress defines whether log should be compressed.
	Compress bool `yaml:"compress"`
	// MaxSize is the max size of log file(MB).
	MaxSize int `yaml:"max_size"`

	// TimeUnit splits files by time unit, like year/month/hour/minute, default
	// as day. It takes effect only when split by time.
	TimeUnit TimeUnit `yaml:"time_unit"`
}

// FormatConfig is the log format config.
type FormatConfig struct {
	// TimeFormat specifies the time format for log output, with a default
	// value of "2006-01-02 15:04:05.000" when left empty.
	TimeFormat string `yaml:"time_format"`

	// TimeKey is the time key of log output, default as "time".
	TimeKey string `yaml:"time_key"`
	// LevelKey is the level key of log output, default as "level".
	LevelKey string `yaml:"level_key"`
	// LoggerKey is the logger key of log output, default as "logger".
	NameKey string `yaml:"name_key"`
	// CallerKey is the caller key of log output, default as "caller".
	CallerKey string `yaml:"caller_key"`
	// FunctionKey is the function key of log output, default as "", which means
	// not to print function name.
	FunctionKey string `yaml:"function_key"`
	// MessageKey is the message key of log output, default as "msg".
	MessageKey string `yaml:"message_key"`
	// StacktraceKey is the stack trace key of log output, default as "stacktrace".
	StacktraceKey string `yaml:"stacktrace_key"`
}

// WriteMode is the log write mode, one of 1, 2, 3.
type WriteMode int

const (
	// WriteSync writes synchronously.
	WriteSync WriteMode = iota + 1
	// WriteAsync writes asynchronously.
	WriteAsync
	// WriteFast writes fast(may drop logs asynchronously).
	WriteFast
)

var (
	// WriteModeStrings is the map from write mod to its string representation.
	WriteModeStrings = map[WriteMode]string{
		WriteSync:  "sync",
		WriteAsync: "async",
		WriteFast:  "fast",
	}
	// WriteModeNames is the map from string to write mod.
	WriteModeNames = map[string]WriteMode{
		"sync":  WriteSync,
		"async": WriteAsync,
		"fast":  WriteFast,
	}
)

func GetWriteMode(mode string) WriteMode {
	if v, ok := WriteModeNames[mode]; ok {
		return v
	}
	return WriteAsync
}

func (wm *WriteMode) String() string {
	return WriteModeStrings[*wm]
}

// RollType is the log rolling type, one of 1, 2.
type RollType int

const (
	// RollingBySize rolls logs by file size.
	RollingBySize RollType = iota + 1
	// RollingByTime rolls logs by time.
	RollingByTime

	RollingBySizeStr = "size"
	RollingByTimeStr = "time"
)

var (
	// RollTypeStrings is the map from rolling type to its string representation.
	RollTypeStrings = map[RollType]string{
		RollingBySize: RollingBySizeStr,
		RollingByTime: RollingByTimeStr,
	}
	// RollTypeNames is the map from string to rolling type.
	RollTypeNames = map[string]RollType{
		RollingBySizeStr: RollingBySize,
		RollingByTimeStr: RollingByTime,
	}
)

func GetRollingType(rollingType RollType) string {
	if v, ok := RollTypeStrings[rollingType]; ok {
		return v
	}
	return RollingByTimeStr
}

func (rt *RollType) String() string {
	return RollTypeStrings[*rt]
}

// Some common used time formats.
const (
	// TimeFormatMinute is accurate to the minute.
	TimeFormatMinute = "%Y%m%d%H%M"
	// TimeFormatHour is accurate to the hour.
	TimeFormatHour = "%Y%m%d%H"
	// TimeFormatDay is accurate to the day.
	TimeFormatDay = "%Y%m%d"
	// TimeFormatMonth is accurate to the month.
	TimeFormatMonth = "%Y%m"
	// TimeFormatYear is accurate to the year.
	TimeFormatYear = "%Y"
)

const (
	// Minute splits by the minute.
	Minute = "minute"
	// Hour splits by the hour.
	Hour = "hour"
	// Day splits by the day.
	Day = "day"
	// Month splits by the month.
	Month = "month"
	// Year splits by the year.
	Year = "year"
)

// TimeUnit is the time unit by which files are split, one of minute/hour/day/month/year.
type TimeUnit string

// Format returns a string preceding with `.`. Use TimeFormatDay as default.
func (t TimeUnit) Format() string {
	var timeFmt string
	switch t {
	case Minute:
		timeFmt = TimeFormatMinute
	case Hour:
		timeFmt = TimeFormatHour
	case Day:
		timeFmt = TimeFormatDay
	case Month:
		timeFmt = TimeFormatMonth
	case Year:
		timeFmt = TimeFormatYear
	default:
		timeFmt = TimeFormatDay
	}
	return "." + timeFmt
}

// RotationGap returns the time.Duration for time unit. Use one day as the default.
func (t TimeUnit) RotationGap() time.Duration {
	switch t {
	case Minute:
		return time.Minute
	case Hour:
		return time.Hour
	case Day:
		return time.Hour * 24
	case Month:
		return time.Hour * 24 * 30
	case Year:
		return time.Hour * 24 * 365
	default:
		return time.Hour * 24
	}
}
