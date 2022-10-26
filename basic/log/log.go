package log

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	//初始化日志
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		QuoteEmptyFields: true, //empty field will set in ""
		ForceColors:      true,
		FullTimestamp:    true,
		DisableQuote:     true,
		TimestampFormat:  "2006-01-02 15:04:05 ",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理函数名
			fs := strings.Split(frame.Function, ".")
			fun := ""
			if len(fs) > 0 {
				fun = fs[len(fs)-1]
			}
			fileName := path.Base(frame.File)
			return fmt.Sprintf("[\033[1;34m%s\033[0m]", fun), fmt.Sprintf("[%s:%d]", fileName, frame.Line)
		},
	})
}

func getLogLevel(logLevel string) log.Level {
	level := log.InfoLevel
	logLevel = strings.ToUpper(logLevel)
	switch logLevel {
	case "DEBUG":
		level = log.DebugLevel
	case "INFO":
		level = log.InfoLevel
	case "ERROR":
		level = log.ErrorLevel
	case "FATAL":
		level = log.FatalLevel
	case "TRACE":
		level = log.TraceLevel
	case "WARN":
		level = log.WarnLevel
	}
	return level
}

//Support modify log level dynamicly
func InitLogServer(addr string) {
	router := gin.Default()

	router.PUT("/config", func(c *gin.Context) {
		oldlevel := getLogLevel(viper.GetString("LOG_LEVEL"))
		log.Info("OldLevel", oldlevel)

		expire := c.Query("expire")
		newlevel := c.Query("logLevel")

		level := getLogLevel(newlevel)
		log.SetLevel(level)
		log.Info("log level", level) //打印日志级别

		//显示log和viper中的日志级别
		c.JSON(200, gin.H{
			"req time":          expire,
			"req level":         level,
			"successed:":        true,
			"log level now at ": log.GetLevel(),
		})

		if expire != "" {
			tm, _ := strconv.Atoi(expire)
			timer := time.NewTimer(time.Duration(tm) * time.Second)
			go func() {
				<-timer.C
				timer.Stop()
				//定时任务结束，还原日志级别
				log.SetLevel(oldlevel)
			}()
		}
	})

	go router.Run(addr)
}
