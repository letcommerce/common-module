package logs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"strings"
)

var (
	LogWriter io.Writer
)

func InitLogger(path string, env string) {
	// Setting Gin Logger
	f, err := os.Create(path)
	if err != nil {
		log.Panicf("can't open log file: gin.log, error: %s", err)
	}
	LogWriter = io.MultiWriter(os.Stdout, f)
	gin.DefaultWriter = LogWriter
	if env != "local" {
		gin.SetMode(gin.ReleaseMode)
	}
	log.SetReportCaller(true)
	log.SetOutput(LogWriter)
	log.SetFormatter(&PlainFormatter{TimestampFormat: "2006-01-02 15:04:05", LevelDesc: []string{"PANIC", "FATAL", "ERROR", "WARN", "INFO", "DEBUG"}})
}

type PlainFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := fmt.Sprintf(entry.Time.Format(f.TimestampFormat))
	return []byte(fmt.Sprintf("[%s] [%s] - %s [%v]\n", f.LevelDesc[entry.Level], timestamp, entry.Message, Caller(entry.Caller))), nil
}

func Caller(f *runtime.Frame) string {
	p, _ := os.Getwd()
	return fmt.Sprintf("%s:%d", strings.TrimPrefix(f.File, p), f.Line)
}
