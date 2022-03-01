package logs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	requestid "github.com/sumit-tembe/gin-requestid"
	"io"
	"os"
	"runtime"
	"strings"
)

var (
	LogWriter   io.Writer
	ServiceName string
	Env         string
	ctx         *gin.Context
)

func SetRequestId(ginCtx *gin.Context) {
	ctx = ginCtx
}

func InitLogger(path string, env string, serviceName string) {
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
	Env = env
	ServiceName = serviceName
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
	return []byte(fmt.Sprintf("[%s] [%s] - %s [%v] [%v:%v - %v]\n", f.LevelDesc[entry.Level], timestamp, entry.Message, requestid.GetRequestIDFromContext(ctx), ServiceName, Env, Caller(entry.Caller))), nil
}

func Caller(f *runtime.Frame) string {
	p, _ := os.Getwd()
	return fmt.Sprintf("%s:%d", strings.TrimPrefix(f.File, p), f.Line)
}
