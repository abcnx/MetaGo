package mvsv

import (
	"os"
	"path/filepath"
	"testing"
)

// 测试数据
const testMVSVContent = `# 标题 : "黄金分钟级行情 - 2026-05-21"
# Title : "Gold Minute-level Quotes - 2026-05-21"
# 数据供应商 : xxx行情采集程序
# DataProvider : xxx Quote Collector
# 字段 : Timestamp|Open|High|Low|Close|Volume
# Field : Timestamp|Open|High|Low|Close|Volume
# 字段名称 : 时间戳|开盘|最高|最低|收盘|成交量
# FieldName : 时间戳|开盘|最高|最低|收盘|成交量
# 字段类型 : timestamp|number|number|number|number|integer
# FieldType : timestamp|number|number|number|number|integer
# 计数 : 3
# Count : 3
# 时区 : Asia/Shanghai
# Timezone : Asia/Shanghai
# 货币 : CNY
# Currency : CNY
# 单位 : 元/克
# Unit : CNY/g
# 备注 : "开盘2345.00 收盘2350.00"
# Remark : "Open 2345.00 Close 2350.00"

09:00|2345.00|2350.00|2340.00|2348.50|12345
09:01|2348.50|2352.00|2347.00|2351.00|13456
09:02|2351.00|2355.00|2350.00|2353.50|14567
`

// 测试纯数据（无元数据）
const testPureDataContent = `Timestamp|Open|High|Low|Close|Volume
09:00|2345.00|2350.00|2340.00|2348.50|12345
09:01|2348.50|2352.00|2347.00|2351.00|13456
`

// 测试空值处理
const testNullValueContent = `# 字段 : Timestamp|Open|High|Low|Close|Volume
# Field : Timestamp|Open|High|Low|Close|Volume

09:00|2345.00|2350.00||2348.50|12345
09:01|2348.50||2347.00|2351.00|13456
`

// 测试特殊字符
const testSpecialCharContent = `# 字段 : Code|Name|Path
# Field : Code|Name|Path

001|测试\|名称|path\\to\\file
002|普通名称|normal\\path
`

func TestParse(t *testing.T) {
	// 创建临时测试文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.mvsv")
	
	err := os.WriteFile(testFile, []byte(testMVSVContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	// 解析文件
	parser := NewParser()
	data, err := parser.Parse(testFile)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}
	
	// 验证元数据
	if data.Metadata.Title != "黄金分钟级行情 - 2026-05-21" {
		t.Errorf("Expected title '黄金分钟级行情 - 2026-05-21', got '%s'", data.Metadata.Title)
	}
	
	if data.Metadata.TitleEn != "Gold Minute-level Quotes - 2026-05-21" {
		t.Errorf("Expected titleEn 'Gold Minute-level Quotes - 2026-05-21', got '%s'", data.Metadata.TitleEn)
	}
	
	if data.Metadata.DataProvider != "xxx行情采集程序" {
		t.Errorf("Expected dataProvider 'xxx行情采集程序', got '%s'", data.Metadata.DataProvider)
	}
	
	if data.Metadata.Field != "Timestamp|Open|High|Low|Close|Volume" {
		t.Errorf("Expected field 'Timestamp|Open|High|Low|Close|Volume', got '%s'", data.Metadata.Field)
	}
	
	if data.Metadata.FieldName != "时间戳|开盘|最高|最低|收盘|成交量" {
		t.Errorf("Expected fieldName '时间戳|开盘|最高|最低|收盘|成交量', got '%s'", data.Metadata.FieldName)
	}
	
	if data.Metadata.FieldType != "timestamp|number|number|number|number|integer" {
		t.Errorf("Expected fieldType 'timestamp|number|number|number|number|integer', got '%s'", data.Metadata.FieldType)
	}
	
	if data.Metadata.Count != 3 {
		t.Errorf("Expected count 3, got %d", data.Metadata.Count)
	}
	
	// 验证字段列表
	fieldList := data.Metadata.GetFieldList()
	if len(fieldList) != 6 {
		t.Errorf("Expected 6 fields, got %d", len(fieldList))
	}
	
	expectedFields := []string{"Timestamp", "Open", "High", "Low", "Close", "Volume"}
	for i, field := range expectedFields {
		if fieldList[i] != field {
			t.Errorf("Expected field[%d] '%s', got '%s'", i, field, fieldList[i])
		}
	}
	
	// 验证字段名称列表
	fieldNameList := data.Metadata.GetFieldNameList()
	if len(fieldNameList) != 6 {
		t.Errorf("Expected 6 field names, got %d", len(fieldNameList))
	}
	
	expectedFieldNames := []string{"时间戳", "开盘", "最高", "最低", "收盘", "成交量"}
	for i, name := range expectedFieldNames {
		if fieldNameList[i] != name {
			t.Errorf("Expected fieldName[%d] '%s', got '%s'", i, name, fieldNameList[i])
		}
	}
	
	// 验证字段类型列表
	fieldTypeList := data.Metadata.GetFieldTypeList()
	if len(fieldTypeList) != 6 {
		t.Errorf("Expected 6 field types, got %d", len(fieldTypeList))
	}
	
	expectedFieldTypes := []string{"timestamp", "number", "number", "number", "number", "integer"}
	for i, typ := range expectedFieldTypes {
		if fieldTypeList[i] != typ {
			t.Errorf("Expected fieldType[%d] '%s', got '%s'", i, typ, fieldTypeList[i])
		}
	}
	
	// 验证数据行
	if len(data.Rows) != 3 {
		t.Fatalf("Expected 3 rows, got %d", len(data.Rows))
	}
	
	// 验证第一行数据
	expectedRow1 := []string{"09:00", "2345.00", "2350.00", "2340.00", "2348.50", "12345"}
	for i, val := range expectedRow1 {
		if data.Rows[0][i] != val {
			t.Errorf("Expected row1[%d] '%s', got '%s'", i, val, data.Rows[0][i])
		}
	}
}

func TestParseString(t *testing.T) {
	parser := NewParser()
	data, err := parser.ParseString(testMVSVContent)
	if err != nil {
		t.Fatalf("Failed to parse string: %v", err)
	}
	
	// 验证基本解析
	if data.Metadata.Title != "黄金分钟级行情 - 2026-05-21" {
		t.Errorf("Expected title '黄金分钟级行情 - 2026-05-21', got '%s'", data.Metadata.Title)
	}
	
	if len(data.Rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(data.Rows))
	}
}

func TestParsePureData(t *testing.T) {
	// 测试纯数据（无元数据）
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "pure.mvsv")
	
	err := os.WriteFile(testFile, []byte(testPureDataContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	parser := NewParser()
	data, err := parser.Parse(testFile)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}
	
	// 验证无元数据时解析正常
	if data.Metadata.Title != "" {
		t.Errorf("Expected empty title, got '%s'", data.Metadata.Title)
	}
	
	if len(data.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(data.Rows))
	}
}

func TestParseNullValue(t *testing.T) {
	// 测试空值处理
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "null.mvsv")
	
	err := os.WriteFile(testFile, []byte(testNullValueContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	parser := NewParser()
	data, err := parser.Parse(testFile)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}
	
	// 验证空值处理
	if len(data.Rows) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(data.Rows))
	}
	
	// 第一行 Low 字段为空
	if data.Rows[0][3] != "" {
		t.Errorf("Expected empty value for Low field in row 0, got '%s'", data.Rows[0][3])
	}
	
	// 第二行 High 字段为空
	if data.Rows[1][2] != "" {
		t.Errorf("Expected empty value for High field in row 1, got '%s'", data.Rows[1][2])
	}
}

func TestSerialize(t *testing.T) {
	// 创建测试数据
	metadata := Metadata{
		Title:          "黄金分钟级行情 - 2026-05-21",
		TitleEn:        "Gold Minute-level Quotes - 2026-05-21",
		DataProvider:   "xxx行情采集程序",
		DataProviderEn: "xxx Quote Collector",
		Field:          "Timestamp|Open|High|Low|Close|Volume",
		FieldEn:        "Timestamp|Open|High|Low|Close|Volume",
		FieldName:      "时间戳|开盘|最高|最低|收盘|成交量",
		FieldNameEn:    "时间戳|开盘|最高|最低|收盘|成交量",
		FieldType:      "timestamp|number|number|number|number|integer",
		FieldTypeEn:    "timestamp|number|number|number|number|integer",
		Count:          3,
		Remark:         "开盘2345.00 收盘2350.00",
		RemarkEn:       "Open 2345.00 Close 2350.00",
	}
	
	rows := [][]string{
		{"09:00", "2345.00", "2350.00", "2340.00", "2348.50", "12345"},
		{"09:01", "2348.50", "2352.00", "2347.00", "2351.00", "13456"},
		{"09:02", "2351.00", "2355.00", "2350.00", "2353.50", "14567"},
	}
	
	data := &Data{
		Metadata:   metadata,
		Headers:    []string{"Timestamp", "Open", "High", "Low", "Close", "Volume"},
		FieldNames: []string{"时间戳", "开盘", "最高", "最低", "收盘", "成交量"},
		FieldTypes: []string{"timestamp", "number", "number", "number", "number", "integer"},
		Rows:       rows,
	}
	
	// 序列化到临时文件
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.mvsv")
	
	serializer := NewSerializer()
	err := serializer.Serialize(data, outputFile)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}
	
	// 读取并验证
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	// 验证包含元数据
	contentStr := string(content)
	if !contains(contentStr, "# 标题 : \"黄金分钟级行情 - 2026-05-21\"") {
		t.Error("Expected title metadata in output")
	}
	
	if !contains(contentStr, "# 字段 : Timestamp|Open|High|Low|Close|Volume") {
		t.Error("Expected field metadata in output")
	}
	
	// 验证包含数据行
	if !contains(contentStr, "09:00|2345.00|2350.00|2340.00|2348.50|12345") {
		t.Error("Expected data row in output")
	}
}

func TestSerializeToString(t *testing.T) {
	metadata := Metadata{
		Title:        "测试数据",
		Field:        "A|B|C",
		FieldName:    "字段A|字段B|字段C",
		FieldType:    "string|number|integer",
		Count:        1,
	}
	
	rows := [][]string{
		{"val1", "1.5", "100"},
	}
	
	data := &Data{
		Metadata:   metadata,
		Rows:       rows,
	}
	
	serializer := NewSerializer()
	result := serializer.SerializeToString(data)
	
	// 验证输出
	if !contains(result, "# 标题 : \"测试数据\"") {
		t.Error("Expected title in serialized string")
	}
	
	if !contains(result, "val1|1.5|100") {
		t.Error("Expected data row in serialized string")
	}
}

func TestRoundTrip(t *testing.T) {
	// 测试解析后再序列化的完整流程
	parser := NewParser()
	data, err := parser.ParseString(testMVSVContent)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	
	serializer := NewSerializer()
	result := serializer.SerializeToString(data)
	
	// 再次解析
	data2, err := parser.ParseString(result)
	if err != nil {
		t.Fatalf("Failed to parse serialized content: %v", err)
	}
	
	// 验证两次解析结果一致
	if data.Metadata.Title != data2.Metadata.Title {
		t.Errorf("Title mismatch: '%s' vs '%s'", data.Metadata.Title, data2.Metadata.Title)
	}
	
	if data.Metadata.Count != data2.Metadata.Count {
		t.Errorf("Count mismatch: %d vs %d", data.Metadata.Count, data2.Metadata.Count)
	}
	
	if len(data.Rows) != len(data2.Rows) {
		t.Errorf("Row count mismatch: %d vs %d", len(data.Rows), len(data2.Rows))
	}
}

func TestMetadataGetMethods(t *testing.T) {
	metadata := Metadata{
		Field:     "A|B|C|D",
		FieldName: "名称A|名称B|名称C|名称D",
		FieldType: "string|number|integer|boolean",
	}
	
	// 测试 GetFieldList
	fieldList := metadata.GetFieldList()
	if len(fieldList) != 4 {
		t.Errorf("Expected 4 fields, got %d", len(fieldList))
	}
	
	// 测试 GetFieldNameList
	fieldNameList := metadata.GetFieldNameList()
	if len(fieldNameList) != 4 {
		t.Errorf("Expected 4 field names, got %d", len(fieldNameList))
	}
	
	// 测试 GetFieldTypeList
	fieldTypeList := metadata.GetFieldTypeList()
	if len(fieldTypeList) != 4 {
		t.Errorf("Expected 4 field types, got %d", len(fieldTypeList))
	}
	
	// 测试空值情况
	emptyMetadata := Metadata{}
	
	if emptyMetadata.GetFieldList() != nil {
		t.Error("Expected nil for empty field list")
	}
	
	if emptyMetadata.GetFieldNameList() != nil {
		t.Error("Expected nil for empty field name list")
	}
	
	if emptyMetadata.GetFieldTypeList() != nil {
		t.Error("Expected nil for empty field type list")
	}
}

func TestFileNotExist(t *testing.T) {
	parser := NewParser()
	_, err := parser.Parse("/nonexistent/path/file.mvsv")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestChineseEnglishMetadata(t *testing.T) {
	// 测试中英双语元数据解析
	content := `# 标题 : "中文标题"
# Title : "English Title"
# 字段 : A|B|C
# Field : A|B|C
# 计数 : 5
# Count : 5

1|2|3
`
	
	parser := NewParser()
	data, err := parser.ParseString(content)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	
	// 验证中英双语字段都正确解析
	if data.Metadata.Title != "中文标题" {
		t.Errorf("Expected Chinese title, got '%s'", data.Metadata.Title)
	}
	
	if data.Metadata.TitleEn != "English Title" {
		t.Errorf("Expected English title, got '%s'", data.Metadata.TitleEn)
	}
	
	// 验证 Count 可以从中文或英文字段获取
	if data.Metadata.Count != 5 {
		t.Errorf("Expected count 5, got %d", data.Metadata.Count)
	}
}

func TestQuotedValue(t *testing.T) {
	// 测试带引号的值
	content := `# 标题 : "带引号的标题"
# 备注 : "带空格的备注信息"

A|B|C
1|2|3
`
	
	parser := NewParser()
	data, err := parser.ParseString(content)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	
	if data.Metadata.Title != "带引号的标题" {
		t.Errorf("Expected '带引号的标题', got '%s'", data.Metadata.Title)
	}
	
	if data.Metadata.Remark != "带空格的备注信息" {
		t.Errorf("Expected '带空格的备注信息', got '%s'", data.Metadata.Remark)
	}
}

func TestExtraFields(t *testing.T) {
	// 测试扩展字段
	content := `# 标题 : "测试"
# 自定义字段 : 自定义值
# CustomField : CustomValue
# 数据质量 : A
# DataQuality : A

A|B
1|2
`
	
	parser := NewParser()
	data, err := parser.ParseString(content)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	
	// 验证扩展字段存储在 Extra 中
	if data.Metadata.Extra["自定义字段"] != "自定义值" {
		t.Error("Expected custom field in Extra")
	}
	
	if data.Metadata.Extra["CustomField"] != "CustomValue" {
		t.Error("Expected custom field (English) in Extra")
	}
	
	if data.Metadata.Extra["数据质量"] != "A" {
		t.Error("Expected data quality field in Extra")
	}
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		len(s) > len(substr) && contains(s[1:], substr)
}