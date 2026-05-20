package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	logger := New("test-function")
	if logger == nil {
		t.Fatal("New() 返回 nil")
	}
	if logger.functionName != "test-function" {
		t.Errorf("functionName = %s, 期望 test-function", logger.functionName)
	}
	if logger.output != os.Stdout {
		t.Error("默认输出应该是 os.Stdout")
	}
}

func TestNewWithOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithOutput("test", &buf)
	
	if logger == nil {
		t.Fatal("NewWithOutput() 返回 nil")
	}
	if logger.output != &buf {
		t.Error("输出应该是提供的 buffer")
	}
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithOutput("test-info", &buf)
	
	logger.Info("test message")
	
	output := buf.String()
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Info 日志应包含 [INFO], 实际: %s", output)
	}
	if !strings.Contains(output, "test-info") {
		t.Errorf("Info 日志应包含函数名, 实际: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Info 日志应包含消息, 实际: %s", output)
	}
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithOutput("test-error", &buf)
	
	logger.Error("error occurred")
	
	output := buf.String()
	if !strings.Contains(output, "[ERROR]") {
		t.Errorf("Error 日志应包含 [ERROR], 实际: %s", output)
	}
	if !strings.Contains(output, "error occurred") {
		t.Errorf("Error 日志应包含消息, 实际: %s", output)
	}
}

func TestWarn(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithOutput("test-warn", &buf)
	
	logger.Warn("warning message")
	
	output := buf.String()
	if !strings.Contains(output, "[WARN]") {
		t.Errorf("Warn 日志应包含 [WARN], 实际: %s", output)
	}
	if !strings.Contains(output, "warning message") {
		t.Errorf("Warn 日志应包含消息, 实际: %s", output)
	}
}

func TestDebug(t *testing.T) {
	t.Run("设置 DEBUG 级别", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "DEBUG")
		defer os.Unsetenv("LOG_LEVEL")
		
		var buf bytes.Buffer
		logger := NewWithOutput("test-debug", &buf)
		
		logger.Debug("debug message")
		
		output := buf.String()
		if !strings.Contains(output, "[DEBUG]") {
			t.Errorf("LOG_LEVEL=DEBUG 时 Debug 日志应包含 [DEBUG], 实际: %s", output)
		}
	})
	
	t.Run("未设置 DEBUG 级别", func(t *testing.T) {
		os.Unsetenv("LOG_LEVEL")
		
		var buf bytes.Buffer
		logger := NewWithOutput("test-debug", &buf)
		
		logger.Debug("debug message")
		
		output := buf.String()
		if output != "" {
			t.Errorf("未设置 LOG_LEVEL 时 Debug 日志应为空, 实际: %s", output)
		}
	})
}

func TestSetOutput(t *testing.T) {
	logger := New("test-setoutput")
	
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	
	if logger.output != &buf {
		t.Error("SetOutput 未更新输出目标")
	}
	
	logger.Info("test")
	
	if buf.String() == "" {
		t.Error("日志应写入新的输出目标")
	}
}

func TestWithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithOutput("test-fields", &buf)
	
	logger.Info("user action", "userId", 123, "action", "login")
	
	output := buf.String()
	if !strings.Contains(output, "userId") {
		t.Errorf("日志应包含字段键, 实际: %s", output)
	}
	if !strings.Contains(output, "123") {
		t.Errorf("日志应包含字段值, 实际: %s", output)
	}
}

func TestSetGlobalLogLevel(t *testing.T) {
	// 重置全局状态
	globalLogLevel = ""
	
	t.Run("设置 DEBUG 级别", func(t *testing.T) {
		SetGlobalLogLevel("DEBUG")
		if GetGlobalLogLevel() != "DEBUG" {
			t.Errorf("GetGlobalLogLevel() = %s, 期望 DEBUG", GetGlobalLogLevel())
		}
	})
	
	t.Run("设置 ERROR 级别", func(t *testing.T) {
		SetGlobalLogLevel("ERROR")
		if GetGlobalLogLevel() != "ERROR" {
			t.Errorf("GetGlobalLogLevel() = %s, 期望 ERROR", GetGlobalLogLevel())
		}
	})
	
	// 重置
	globalLogLevel = ""
}

func TestShouldLog(t *testing.T) {
	tests := []struct {
		name        string
		setLevel    string
		msgLevel    string
		shouldLog   bool
	}{
		{"DEBUG 级别设置 DEBUG", "DEBUG", "DEBUG", true},
		{"DEBUG 级别设置 INFO", "DEBUG", "INFO", true},
		{"ERROR 级别设置 ERROR", "ERROR", "ERROR", true},
		{"ERROR 级别设置 DEBUG", "ERROR", "DEBUG", false},
		{"WARN 级别设置 INFO", "WARN", "INFO", false},
		{"WARN 级别设置 WARN", "WARN", "WARN", true},
		{"未设置级别", "", "DEBUG", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			globalLogLevel = ""
			if tt.setLevel != "" {
				SetGlobalLogLevel(tt.setLevel)
			}
			
			result := shouldLog(tt.msgLevel)
			if result != tt.shouldLog {
				t.Errorf("shouldLog(%s) 级别 %s = %v, 期望 %v", 
					tt.msgLevel, tt.setLevel, result, tt.shouldLog)
			}
		})
	}
	
	// 重置
	globalLogLevel = ""
}

func TestLogFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithOutput("test-format", &buf)
	
	logger.Info("test message")
	
	output := buf.String()
	
	// 验证格式: [时间] [级别] [函数名] 消息
	if !strings.HasPrefix(output, "[") {
		t.Error("日志应以时间戳方括号开头")
	}
	
	// 验证包含时间戳格式
	if !strings.Contains(output, "T") {
		t.Error("日志应包含 ISO 时间戳格式")
	}
	
	// 验证包含级别
	if !strings.Contains(output, "[INFO]") {
		t.Error("日志应包含级别")
	}
	
	// 验证包含函数名
	if !strings.Contains(output, "[test-format]") {
		t.Error("日志应包含函数名")
	}
}

func TestMultipleLoggers(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	
	logger1 := NewWithOutput("logger1", &buf1)
	logger2 := NewWithOutput("logger2", &buf2)
	
	logger1.Info("from logger1")
	logger2.Info("from logger2")
	
	if !strings.Contains(buf1.String(), "from logger1") {
		t.Error("logger1 应写入 buf1")
	}
	if !strings.Contains(buf2.String(), "from logger2") {
		t.Error("logger2 应写入 buf2")
	}
}
