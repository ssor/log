package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// PrintList define a list to print
type PrintList interface {
	ListName() string
	InfoList() []string
}

const (
	line = "----------------------------------------------------------------------------"
)

var (
	// DebugLevel shows what level is now
	DebugLevel = LevelInfo
	logWriter  io.WriteCloser
	logFile    *os.File
)

func init() {
}

func createLogFile() (err error) {
	fileName := "error_" + time.Now().Format("2006_01_02_15_04_05") + ".log"
	logFile, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	return err
}

// exists returns whether the given file or directory exists or not
func isFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// ToFileF write log to file
func ToFileF(format string, args ...interface{}) {
	log := fmt.Sprintf(format, args...)
	ToFile(log)
}

// ToFile write log to file
func ToFile(log string) {
	if logFile == nil {
		if err := createLogFile(); err != nil {
			return
		}
	}

	_, err := logFile.WriteString(time.Now().Format(time.RFC3339) + " : " + log + getFileLocation() + " \r\n")
	if err != nil {
		logFile.Close()
		logFile = nil
	}
}

// SetLogWriter set a write instead of default
func SetLogWriter(writer io.WriteCloser) {
	logWriter = writer
}

// PrintListInfo print a list at LevelInfo level
func PrintListInfo(list PrintList) {
	printList(list, Info)
	InfoF("%s%s", line, getFileLocation())
}

// PrintListTrace print a list at LevelTrace level
func PrintListTrace(list PrintList) {
	printList(list, Trace)
}

// PrintList print list
func printList(list PrintList, printFunc func(log string)) {
	log := fmt.Sprintf(line+"%s", getFileLocation())
	printFunc(log)
	printFunc(list.ListName() + " 列表：")
	strs := list.InfoList()
	for _, str := range strs {
		printFunc(str)
	}
	printFunc(log)
}

//Must 能够造成系统不正常运行的问题
func Must(log string) {
	debugOutput(log, LevelCritical)
}

//MustF 能够造成系统不正常运行的问题
func MustF(format string, args ...interface{}) {
	log := fmt.Sprintf(format+"%s", append(args, getFileLocation())...)
	Must(log)
}

// SysF system not well
func SysF(format string, args ...interface{}) {
	log := fmt.Sprintf(format+"%s", append(args, getFileLocation())...)
	Sys(log)
}

// Sys 出现异常信息，系统能够正常运行，但是可能和使用者想象的不同
func Sys(log string) {
	debugOutput(log, 2)
}

// InfoF 关键步骤或者信息的提醒
func InfoF(format string, args ...interface{}) {
	log := fmt.Sprintf(format+"%s", append(args, getFileLocation())...)
	Info(log)
}

// Info 关键步骤或者信息的提醒
func Info(log string) {
	debugOutput(log, 3)
}

// TraceF 运行数据的打印
func TraceF(format string, args ...interface{}) {
	log := fmt.Sprintf(format+"%s", append(args, getFileLocation())...)
	Trace(log)
}

// Trace 运行数据的打印
func Trace(log string) {
	debugOutput(log, 4)
}

func debugOutput(log string, level int) {
	if level <= DebugLevel {
		prefix := ""
		switch level {
		case 1:
			prefix = "[ERROR]: "
		case 2:
			prefix = "[WARN]:  "
		case 3:
			prefix = "[INFO]:  "
		case 4:
			prefix = "[TRACE]: "

		}
		output := prefix + log
		fmt.Println(output)

		if logWriter != nil {
			_, err := logWriter.Write([]byte(output))
			if err != nil {
				logWriter.Close()
				logWriter = nil
			}
		}
	}
}

func getFileLocation() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		array := strings.Split(file, "/")
		return fmt.Sprintf(" (%s %d)", array[len(array)-1], line)
	}
	return "  ???"

}
