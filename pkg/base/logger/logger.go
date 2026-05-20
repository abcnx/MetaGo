package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Logger 统一日志接口
type Logger struct {
	functionName string
	output       io.Writer // 可自定义输出目标
}

// New 创建新的日志实例
func New(functionName string) *Logger {
	return &Logger{
		functionName: functionName,
		output:       os.Stdout, // 默认输出到标准输出
	}
}

// NewWithOutput 创建带自定义输出的日志实例
func NewWithOutput(functionName string, output io.Writer) *Logger {
	return &Logger{
		functionName: functionName,
		output:       output,
	}
}

// Info 输出信息日志
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.log("INFO", msg, fields...)
}

// Error 输出错误日志
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.log("ERROR", msg, fields...)
}

// Warn 输出警告日志
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.log("WARN", msg, fields...)
}

// Debug 输出调试日志
func (l *Logger) Debug(msg string, fields ...interface{}) {
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		l.log("DEBUG", msg, fields...)
	}
}

// SetOutput 设置输出目标
func (l *Logger) SetOutput(output io.Writer) {
	l.output = output
}

// log 内部日志方法
func (l *Logger) log(level, msg string, fields ...interface{}) {
	timestamp := time.Now().Format("2006-01-02T15:04:05.000Z07:00")
	output := fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, level, l.functionName, msg)
	if len(fields) > 0 {
		output += fmt.Sprintf(" %v", fields)
	}
	
	// 输出到指定目标
	if l.output != nil {
		fmt.Fprintln(l.output, output)
	} else {
		fmt.Println(output)
	}
}

// 全局日志级别控制
var globalLogLevel string

// SetGlobalLogLevel 设置全局日志级别
// 可选值: DEBUG, INFO, WARN, ERROR
// 低于设定级别的日志不会输出
func SetGlobalLogLevel(level string) {
	globalLogLevel = level
}

// GetGlobalLogLevel 获取全局日志级别
func GetGlobalLogLevel() string {
	if globalLogLevel != "" {
		return globalLogLevel
	}
	return os.Getenv("LOG_LEVEL")
}

// shouldLog 判断是否应该输出日志
func shouldLog(level string) bool {
	currentLevel := GetGlobalLogLevel()
	if currentLevel == "" {
		return true // 未设置级别时输出所有日志
	}
	
	// 日志级别优先级
 levels := map[string]int{
		"DEBUG": 0,
		"INFO":  1,
		"WARN":  2,
		"ERROR": 3,
	}
	
	currentPriority, ok := levels[currentLevel]
	if !ok {
		return true
	}
	
	msgPriority, ok := levels[level]
	if !ok {
		return true
	}
	
	return msgPriority >= currentPriority
}