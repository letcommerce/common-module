package logs

import (
	"bytes"
	"encoding/json"
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

const defaultLogPath = "gin.log"

func SetRequestId(ginCtx *gin.Context) {
	ctx = ginCtx
}

func InitLogger(path string, env string, serviceName string) {
	// Setting Gin Logger
	if path == "" {
		path = defaultLogPath
	}
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
	if env == "prod" {
		log.SetFormatter(&JsonFormatter{TimestampFormat: "2006-01-02 15:04:05", LevelDesc: []string{"PANIC", "FATAL", "ERROR", "WARN", "INFO", "DEBUG"}})
	} else {
		log.SetFormatter(&PlainFormatter{TimestampFormat: "2006-01-02 15:04:05", LevelDesc: []string{"PANIC", "FATAL", "ERROR", "WARN", "INFO", "DEBUG"}})
	}
}

type PlainFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := fmt.Sprintf(entry.Time.Format(f.TimestampFormat))
	var requestId string
	if ctx != nil {
		requestId = requestid.GetRequestIDFromContext(ctx)
	}
	return []byte(fmt.Sprintf("[%s] [%s] - %s [%v:%v:%v - %v]\n", f.LevelDesc[entry.Level], timestamp, entry.Message, ServiceName, Env, requestId, Caller(entry.Caller))), nil
}

func Caller(f *runtime.Frame) string {
	p, _ := os.Getwd()
	fileName := strings.TrimPrefix(f.File, p)
	fileName = strings.ReplaceAll(fileName, "/go/pkg/mod/github.com/letcommerce/", "")
	return fmt.Sprintf("%s:%d", fileName, f.Line)
}

type JsonFormatter struct {
	TimestampFormat   string
	LevelDesc         []string
	PrettyPrint       bool
	DisableHTMLEscape bool
}

func (f *JsonFormatter) Format(entry *log.Entry) ([]byte, error) {
	result := map[string]interface{}{}

	result["time"] = fmt.Sprintf(entry.Time.Format(f.TimestampFormat))
	result["msg"] = entry.Message
	if entry.HasCaller() {
		function, file := CallerWithFunc(entry.Caller)
		if function != "" {
			result["func"] = function
		}
		result["file"] = file
	}
	result["level"] = f.LevelDesc[entry.Level]
	result["service"] = ServiceName
	result["env"] = Env
	if ctx != nil {
		requestId := requestid.GetRequestIDFromContext(ctx)
		result["request_id"] = requestId
	}
	b := &bytes.Buffer{}
	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(!f.DisableHTMLEscape)
	if f.PrettyPrint {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(result); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}

	return b.Bytes(), nil
}

func CallerWithFunc(f *runtime.Frame) (string, string) {
	p, _ := os.Getwd()
	fileName := strings.TrimPrefix(f.File, p)
	fileName = strings.ReplaceAll(fileName, "/go/pkg/mod/github.com/letcommerce/", "")
	return f.Func.Name(), fmt.Sprintf("%s:%d", fileName, f.Line)
}
