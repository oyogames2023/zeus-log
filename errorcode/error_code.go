package errorcode

import "errors"

var (
	ErrInvalidFilePath            = errors.New("invalid file path")
	ErrInvalidTimePattern         = errors.New("invalid time pattern")
	ErrIgnoreCurrentLogFile       = errors.New("ignore current logfile")
	ErrOpenFileFailed             = errors.New("open file failed")
	ErrInvalidWriterDecoderObject = errors.New("invalid writer decoder object")
	ErrInvalidWriterDecoderType   = errors.New("invalid writer decoder type")
)
