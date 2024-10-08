package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Flyingmn/gzap"
	"go.uber.org/zap"
)

var (
	//log.Logger写文件支持并发
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	logLevel    = 0 //默认的LogLevel为0，即所有级别的日志都打印

	logOut        *os.File
	day           int
	dayChangeLock sync.RWMutex
	logFile       string
)

const (
	DebugLevel = iota //iota=0
	InfoLevel         //InfoLevel=iota, iota=1
	WarnLevel         //WarnLevel=iota, iota=2
	ErrorLevel        //ErrorLevel=iota, iota=3
)

func SetLogLevel(level int) {
	logLevel = level

	//zap设置
	switch level {
	case DebugLevel:
		gzap.SetZapCfg(gzap.ZapLevel("debug"))
	case InfoLevel:
		gzap.SetZapCfg(gzap.ZapLevel("info"))
	case WarnLevel:
		gzap.SetZapCfg(gzap.ZapLevel("warn"))
	case ErrorLevel:
		gzap.SetZapCfg(gzap.ZapLevel("error"))
	default:
		gzap.SetZapCfg(gzap.ZapLevel("info"))
	}
}
func SetLogServiceName(serviceName string) {
	gzap.SetZapCfg(gzap.SetPresetFields(map[string]any{"service": serviceName}))
}

func SetLogFile(file string) {
	logFile = file
	now := time.Now()
	var err error
	if logOut, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664); err != nil {
		panic(err)
	} else {
		debugLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags)
		infoLogger = log.New(logOut, "[INFO] ", log.LstdFlags)
		warnLogger = log.New(logOut, "[WARN] ", log.LstdFlags)
		errorLogger = log.New(logOut, "[ERROR] ", log.LstdFlags)
		day = now.YearDay()
		dayChangeLock = sync.RWMutex{}

		//zap设置
		gzap.SetZapCfg(
			gzap.ZapOutFile(
				logFile,                       //文件位置
				gzap.ZapOutFileMaxSize(2048),  // 日志文件的最大大小(MB为单位)
				gzap.ZapOutFileMaxAge(365),    //保留旧文件的最大天数量
				gzap.ZapOutFileMaxBackups(50), //保留旧文件的最大个数
			),
		)
	}
}

// 检查是否需要切换日志文件，如果需要则切换
func checkAndChangeLogfile() {
	dayChangeLock.Lock()
	defer dayChangeLock.Unlock()
	now := time.Now()
	if now.YearDay() == day {
		return
	}

	logOut.Close()
	postFix := now.Add(-24 * time.Hour).Format("20060102") //昨天的日期
	if err := os.Rename(logFile, logFile+"."+postFix); err != nil {
		fmt.Printf("append date postfix %s to log file %s failed: %v\n", postFix, logFile, err)
		return
	}
	var err error
	if logOut, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664); err != nil {
		fmt.Printf("create log file %s failed %v\n", logFile, err)
		return
	} else {
		debugLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags)
		infoLogger = log.New(logOut, "[INFO] ", log.LstdFlags)
		warnLogger = log.New(logOut, "[WARN] ", log.LstdFlags)
		errorLogger = log.New(logOut, "[ERROR] ", log.LstdFlags)
		day = now.YearDay()
	}
}

func addPrefix() string {
	file, _, line := getLineNo()
	arr := strings.Split(file, "/")
	if len(arr) > 3 {
		arr = arr[len(arr)-3:]
	}
	return strings.Join(arr, "/") + ":" + strconv.Itoa(line)
}

func Debug(format string, v ...interface{}) {
	if logLevel <= DebugLevel {
		checkAndChangeLogfile()
		debugLogger.Printf(addPrefix()+" "+format, v...)
	} else {
		fmt.Printf("[DEBUG] "+addPrefix()+" "+format+"\n", v...)
	}
}

func Info(format string, v ...interface{}) {
	if logLevel <= InfoLevel {
		checkAndChangeLogfile()
		infoLogger.Printf(addPrefix()+" "+format, v...) //format末尾如果没有换行符会自动加上
	} else {
		fmt.Printf("[INFO] "+addPrefix()+" "+format+"\n", v...)
	}
}

func Warn(format string, v ...interface{}) {
	if logLevel <= WarnLevel {
		checkAndChangeLogfile()
		warnLogger.Printf(addPrefix()+" "+format, v...)
	} else {
		fmt.Printf("[WARN] "+addPrefix()+" "+format+"\n", v...)
	}
}

func Error(format string, v ...interface{}) {
	if logLevel <= ErrorLevel {
		checkAndChangeLogfile()
		errorLogger.Printf(addPrefix()+" "+format, v...)
	} else {
		fmt.Printf("[ERROR] "+addPrefix()+" "+format+"\n", v...)
	}
}

func getLineNo() (string, string, int) {
	funcName, file, line, ok := runtime.Caller(3)
	if ok {
		return file, runtime.FuncForPC(funcName).Name(), line
	} else {
		return "", "", 0
	}
}

func Debugz(msg string, fields ...zap.Field) {
	gzap.Debug(msg, fields...)
}
func Infoz(msg string, fields ...zap.Field) {
	gzap.Info(msg, fields...)
}
func Warnz(msg string, fields ...zap.Field) {
	gzap.Warn(msg, fields...)
}
func Errorz(msg string, fields ...zap.Field) {
	gzap.Error(msg, fields...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	gzap.Debugw(msg, keysAndValues...)
}
func Infow(msg string, keysAndValues ...interface{}) {
	gzap.Infow(msg, keysAndValues...)
}
func Warnw(msg string, keysAndValues ...interface{}) {
	gzap.Warnw(msg, keysAndValues...)
}
func Errorw(msg string, keysAndValues ...interface{}) {
	gzap.Errorw(msg, keysAndValues...)
}
