package logger

import "github.com/tryfix/log"

var (
	NoopLogger log.Logger
)

func init() {
	NoopLogger = log.NewLog(log.WithColors(true), log.WithFilePath(true), log.WithLevel("DEBUG")).
		Log(log.Prefixed("saban"))
}
