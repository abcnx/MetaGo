package rest

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ACANX/MetaGo/pkg/base/logger"
)

func TestNewSuccess(t *testing.T) {
	log := logger.New("rest-test")
	log.Info("开始 TestNewSuccess 测试")
	fmt.Println("=== TestNewSuccess 测试开始 ===")

	data := map[string]interface{}{
		"rate": 7.25,
		"name": "test",
	}

	log.Info("创建成功响应", "data", data)
	fmt.Printf("输入数据: %v\n", data)
	result := NewSuccess(data)

	fmt.Printf("结果: Code=%d, Message='%s', Data=%v\n", result.Code, result.Message, result.Data)

	if result.Code != 1 {
		log.Error("Code不匹配", "got", result.Code, "want", 1)
		fmt.Printf("【错误】Code=%d, 期望1\n", result.Code)
		t.Errorf("NewSuccess().Code = %d, 期望 1", result.Code)
	} else {
		log.Info("Code验证通过", "code", result.Code)
		fmt.Printf("【成功】Code=%d\n", result.Code)
	}

	if result.Message != "OK" {
		log.Error("Message不匹配", "got", result.Message, "want", "OK")
		fmt.Printf("【错误】Message='%s', 期望'OK'\n", result.Message)
		t.Errorf("NewSuccess().Message = %s, 期望 OK", result.Message)
	} else {
		log.Info("Message验证通过", "message", result.Message)
		fmt.Printf("【成功】Message='%s'\n", result.Message)
	}

	if result.Data == nil {
		log.Error("Data为nil", nil)
		fmt.Println("【错误】Data为nil")
		t.Error("NewSuccess().Data为nil")
	} else {
		log.Info("Data验证通过", "data", result.Data)
		fmt.Printf("【成功】Data=%v\n", result.Data)
	}

	log.Info("TestNewSuccess测试完成")
	fmt.Println("=== TestNewSuccess 测试完成 ===")
}

func TestNewError(t *testing.T) {
	log := logger.New("rest-test")
	log.Info("开始 TestNewError 测试")
	fmt.Println("=== TestNewError 测试开始 ===")

	message := "测试错误消息"
	log.Info("创建错误响应", "message", message)
	fmt.Printf("输入消息: '%s'\n", message)
	result := NewError(message)

	fmt.Printf("结果: Code=%d, Message='%s', Data=%v\n", result.Code, result.Message, result.Data)

	if result.Code != 0 {
		log.Error("Code不匹配", "got", result.Code, "want", 0)
		fmt.Printf("【错误】Code=%d, 期望0\n", result.Code)
		t.Errorf("NewError().Code = %d, 期望 0", result.Code)
	} else {
		log.Info("Code验证通过", "code", result.Code)
		fmt.Printf("【成功】Code=%d (错误状态)\n", result.Code)
	}

	if result.Message != message {
		log.Error("Message不匹配", "got", result.Message, "want", message)
		fmt.Printf("【错误】Message='%s', 期望'%s'\n", result.Message, message)
		t.Errorf("NewError().Message = %s, 期望 %s", result.Message, message)
	} else {
		log.Info("Message验证通过", "message", result.Message)
		fmt.Printf("【成功】Message='%s'\n", result.Message)
	}

	if result.Data != nil {
		log.Error("Data应为nil", "data", result.Data)
		fmt.Printf("【错误】Data=%v, 期望nil\n", result.Data)
		t.Error("NewError().Data应为nil")
	} else {
		log.Info("Data验证为nil")
		fmt.Println("【成功】Data=nil")
	}

	log.Info("TestNewError测试完成")
	fmt.Println("=== TestNewError 测试完成 ===")
}

func TestRestResult_JSON(t *testing.T) {
	log := logger.New("rest-test-json")
	log.Info("开始 TestRestResult_JSON 测试")
	fmt.Println("=== TestRestResult_JSON 测试开始 ===")

	tests := []struct {
		name     string
		result   *RestResult
		wantCode int
	}{
		{
			name:     "成功响应",
			result:   NewSuccess(map[string]interface{}{"key": "value"}),
			wantCode: 1,
		},
		{
			name:     "错误响应",
			result:   NewError("出错了"),
			wantCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Info("测试JSON序列化", "name", tt.name, "wantCode", tt.wantCode)
			fmt.Printf("\n--- 测试: %s ---\n", tt.name)

			jsonBytes, err := json.Marshal(tt.result)
			if err != nil {
				log.Error("RestResult序列化失败", err)
				fmt.Printf("【错误】序列化失败: %v\n", err)
				t.Fatalf("RestResult序列化失败: %v", err)
			}
			log.Info("JSON序列化成功", "json", string(jsonBytes))
			fmt.Printf("JSON输出: %s\n", string(jsonBytes))

			var unmarshaled RestResult
			if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
				log.Error("RestResult反序列化失败", err)
				fmt.Printf("【错误】反序列化失败: %v\n", err)
				t.Fatalf("RestResult反序列化失败: %v", err)
			}

			fmt.Printf("反序列化结果: Code=%d, Message='%s'\n", unmarshaled.Code, unmarshaled.Message)

			if unmarshaled.Code != tt.wantCode {
				log.Error("Code不匹配", "got", unmarshaled.Code, "want", tt.wantCode)
				fmt.Printf("【错误】Code=%d, 期望%d\n", unmarshaled.Code, tt.wantCode)
				t.Errorf("Code = %d, 期望 %d", unmarshaled.Code, tt.wantCode)
			} else {
				log.Info("Code验证通过", "code", unmarshaled.Code)
				fmt.Printf("【成功】Code=%d\n", unmarshaled.Code)
			}
			fmt.Printf("--- 测试结束: %s ---\n", tt.name)
		})
	}

	log.Info("TestRestResult_JSON测试完成")
	fmt.Println("=== TestRestResult_JSON 测试完成 ===")
}

func TestRestResult_WithComplexData(t *testing.T) {
	log := logger.New("rest-test-complex")
	log.Info("开始 TestRestResult_WithComplexData 测试")
	fmt.Println("=== TestRestResult_WithComplexData 测试开始 ===")

	complexData := map[string]interface{}{
		"nested": map[string]interface{}{
			"key": "value",
			"array": []interface{}{1, 2, 3},
		},
		"number": 123.45,
		"string": "test",
	}

	log.Info("创建带复杂数据的成功响应", "fields", len(complexData))
	fmt.Printf("输入数据 (复杂): %v\n", complexData)
	result := NewSuccess(complexData)

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		log.Error("RestResult序列化失败", err)
		fmt.Printf("【错误】序列化失败: %v\n", err)
		t.Fatalf("RestResult序列化失败: %v", err)
	}
	log.Info("JSON序列化成功", "jsonLength", len(jsonBytes))
	fmt.Printf("JSON输出: %s\n", string(jsonBytes))

	var unmarshaled RestResult
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		log.Error("RestResult反序列化失败", err)
		fmt.Printf("【错误】反序列化失败: %v\n", err)
		t.Fatalf("RestResult反序列化失败: %v", err)
	}

	// 验证 data 字段存在
	if unmarshaled.Data == nil {
		log.Error("反序列化后Data为nil", nil)
		fmt.Println("【错误】Data为nil")
		t.Error("反序列化后Data为nil")
	} else {
		log.Info("Data验证通过", "data", unmarshaled.Data)
		fmt.Printf("【成功】Data=%v\n", unmarshaled.Data)
	}

	log.Info("TestRestResult_WithComplexData测试完成")
	fmt.Println("=== TestRestResult_WithComplexData 测试完成 ===")
}