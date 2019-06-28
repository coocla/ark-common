package logs

import (
	"ark-common/constants"
	"ark-common/utils/filesystem"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// SetupLog 初始化日志配置
// normalName: 正常输出的文件名
// errName: 错误输出的文件名
func SetupLog(normalName, errName, logLevel string) {
	normalFp, err := filesystem.OpenFile(normalName, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Printf("open log file %s failed: %v\n", normalName, err)
		os.Exit(1)
	}

	errFp, err := filesystem.OpenFile(errName, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Printf("open log file %s failed: %v\n", errName, err)
		os.Exit(1)
	}

	loglevel, err := log.ParseLevel(logLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.SetLevel(loglevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   true,
		TimestampFormat: constants.ISO8601,
	})
	// if loglevel == log.DebugLevel {
	// 	log.SetReportCaller(true)
	// }
	log.AddHook(NewFilterHook(normalFp, []log.Level{
		log.InfoLevel,
		log.DebugLevel,
	}))
	log.AddHook(NewFilterHook(errFp, []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
	}))
}

// FilterHook 自定义过滤器
type FilterHook struct {
	Writer       io.Writer
	LogLevels    []log.Level
	filterLevels map[log.Level]bool
}

// NewFilterHook 返回特定的过滤器
func NewFilterHook(writer io.Writer, levels []log.Level) *FilterHook {
	f := &FilterHook{
		Writer:       writer,
		LogLevels:    levels,
		filterLevels: make(map[log.Level]bool),
	}
	f.setup()
	return f
}

func (hook *FilterHook) setup() {
	for idx := range hook.LogLevels {
		hook.filterLevels[hook.LogLevels[idx]] = true
	}
}

// Levels 实现hook的接口
func (hook *FilterHook) Levels() []log.Level {
	return log.AllLevels
}

// Fire 实现hook的接口
func (hook *FilterHook) Fire(entry *log.Entry) (err error) {
	line, err := entry.String()
	if err != nil {
		return
	}
	if ok := hook.filterLevels[entry.Level]; ok {
		_, err = hook.Writer.Write([]byte(line))
	}
	return
}
