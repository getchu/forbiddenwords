package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogLevel int

const (
	LOG_DEBUG   = LogLevel(1)
	LOG_INFO    = LogLevel(2)
	LOG_WARNING = LogLevel(3)
	LOG_ERROR   = LogLevel(4)
	LOG_FATAL   = LogLevel(5)
)

func (this LogLevel) String() string {
	switch this {
	case LOG_DEBUG:
		return "debug"
	case LOG_INFO:
		return "info"
	case LOG_WARNING:
		return "warning"
	case LOG_ERROR:
		return "error"
	case LOG_FATAL:
		return "fatal"
	default:
		return "info"
	}
}

type Logger struct {
	*log.Logger
	level   LogLevel
	logFile *os.File
}

func getLevelByString(levelString string) LogLevel {
	switch strings.ToLower(levelString) {
	case "debug":
		return LOG_DEBUG
	case "info":
		return LOG_INFO
	case "warning":
		return LOG_WARNING
	case "error":
		return LOG_ERROR
	case "fatal":
		return LOG_FATAL
	default:
		return LOG_INFO
	}
}

func NewLogger(fileNamePrefix string, prefix string, level string) *Logger {
	//目录是否存在
	pathSliceList := strings.Split(fileNamePrefix, "/")
	//不是当前目录
	if len(pathSliceList) > 1 {
		pathSliceList = pathSliceList[:len(pathSliceList)-1]
		path := strings.Join(pathSliceList, "/")
		_, err := os.Stat(path)
		//创建目录
		if err != nil && os.IsNotExist(err) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				log.Fatalln("create dir failed: " + err.Error())
			}
		}
	}
	//加日期后缀
	fileName := fileNamePrefix + "." + time.Now().Format("2006-01-02")
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("open log file " + fileName + " error: " + err.Error())
	}
	logger := log.New(logFile, prefix, log.LstdFlags)
	logger.SetFlags(log.Ldate | log.Ltime)
	//删除7天前的日志
	//获取当前日志文件所在路径
	parentPath := filepath.Dir(fileName)
	//七天前时间戳
	oldestFileNameShuffix := time.Now().Unix() - int64(86400*7)
	filepath.Walk(parentPath, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		shuffix := strings.TrimPrefix(path, fileNamePrefix)
		if len(shuffix) < 11 {
			return nil
		}
		shuffix = strings.Replace(shuffix, ".", "", 1)
		t, err := time.Parse("2006-01-02", shuffix)
		if err != nil {
			return nil
		}
		if t.Unix() < oldestFileNameShuffix {
			os.Remove(path)
		}
		return nil
	})
	return &Logger{
		Logger:  logger,
		level:   getLevelByString(level),
		logFile: logFile,
	}
}

func (this *Logger) Output(msgLevel LogLevel, content string, args ...interface{}) {
	if msgLevel < this.level {
		return
	}
	if len(args) > 0 {
		content = fmt.Sprintf("["+msgLevel.String()+"]"+content, args...)
	} else {
		content = "[" + msgLevel.String() + "]" + content
	}
	this.Logger.Println(content)
}

func (this *Logger) Close() bool {
	this.logFile.Close()
	return true
}
