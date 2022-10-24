package nsq

import (
	"strings"

	"github.com/nsqio/go-nsq"
	log "github.com/sirupsen/logrus"
)

type nsqLogrusLogger struct{}

func NewNSQLogrusLogger(l log.Level) (nsqLogrusLogger, nsq.LogLevel) {
	level := nsq.LogLevelWarning
	switch l {
	case log.DebugLevel:
		level = nsq.LogLevelDebug
	case log.InfoLevel:
		level = nsq.LogLevelInfo
	case log.WarnLevel:
		level = nsq.LogLevelWarning
	case log.ErrorLevel:
		level = nsq.LogLevelError
	}
	return nsqLogrusLogger{}, level
}

/*
go-nsq consumer.go line:1191
logger.Output(2, fmt.Sprintf("%-4s %3d [%s/%s] %s",
	lvl, r.id, r.topic, r.channel,
	fmt.Sprintf(line, args...)))
*/
func (n nsqLogrusLogger) Output(_ int, s string) error {
	if len(s) > 3 {
		msg := strings.TrimSpace(s[3:])
		switch s[:3] {
		case nsq.LogLevelDebug.String():
			log.Debugln(msg)
		case nsq.LogLevelInfo.String():
			log.Infoln(msg)
		case nsq.LogLevelWarning.String():
			log.Warnln(msg)
		case nsq.LogLevelError.String():
			log.Errorln(msg)
		default:
			log.Infoln(msg)
		}
	}
	return nil
}
