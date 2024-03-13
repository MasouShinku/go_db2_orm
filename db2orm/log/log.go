package log

import (
	"io"
	"log"
	"os"
	"sync"
)

const (
	preError = "\033[31m[Error]\033[0m"
	preInfo  = "\033[34m[Info]\033[0m"
)

// 设置日志的层级
const (
	InfoLevel  = 0
	ErrorLevel = 1
	Disabled   = 2
)

var (
	errorLog = log.New(os.Stdout, preError, log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, preInfo, log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

// 暴露log方法
var (
	Errorln = errorLog.Println
	Errorf  = errorLog.Printf
	Infoln  = infoLog.Println
	Infof   = infoLog.Printf
)

func setLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		errorLog.SetOutput(io.Discard)
	}

	if InfoLevel < level {
		infoLog.SetOutput(io.Discard)
	}
}
