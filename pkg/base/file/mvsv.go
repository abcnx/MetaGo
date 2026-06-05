package mvsv

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Metadata MVSV 文件元数据
type Metadata struct {
	Title          string // 标题
	TitleEn        string // Title（英文）
	DataProvider   string // 数据供应商
	DataProviderEn string // DataProvider（英文）
	Field          string // 字段（英文字段名）
	FieldEn        string // Field（英文）
	FieldName      string // 字段名称（中文字段名）
	FieldNameEn    string // FieldName（英文）
	FieldType      string // 字段类型
	FieldTypeEn    string // FieldType（英文）
	Count          int    // 计数（数据行数）
	Remark         string // 备注
	RemarkEn       string // Remark（英文）
	Extra          map[string]string // 扩展字段
}

// GetFieldList 获取字段名列表
func (m *Metadata) GetFieldList() []string {
	if m.Field == "" {
		return nil
	}
	return strings.Split(m.Field, "|")
}

// GetFieldNameList 获取字段中文名称列表
func (m *Metadata) GetFieldNameList() []string {
	if m.FieldName == "" {
		return nil
	}
	return strings.Split(m.FieldName, "|")
}

// GetFieldTypeList 获取字段类型列表
func (m *Metadata) GetFieldTypeList() []string {
	if m.FieldType == "" {
		return nil
	}
	return strings.Split(m.FieldType, "|")
}

// Data MVSV 文件数据
type Data struct {
	Metadata  Metadata   // 元数据
	Headers   []string   // 字段名列表
	FieldNames []string  // 字段中文名称列表
	FieldTypes []string  // 字段类型列表
	Rows      [][]string // 数据行
}

// Parser MVSV 文件解析器
type Parser struct{}

// NewParser 创建解析器
func NewParser() *Parser {
	return &Parser{}
}

// Parse 解析 MVSV 文件
func (p *Parser) Parse(filePath string) (*Data, error) {
	lines, err := p.readFile(filePath)
	if err != nil {
		return nil, err
	}

	// 解析元数据区
	metadata := p.parseMetadata(lines)

	// 解析数据区
	_, rows := p.parseData(lines)

	// 从元数据获取 headers
	var headers []string
	if metadata.Field != "" {
		headers = metadata.GetFieldList()
	}
	// 无元数据时，headers 为 nil，所有行都是数据

	return &Data{
		Metadata:   metadata,
		Headers:    headers,
		FieldNames: metadata.GetFieldNameList(),
		FieldTypes: metadata.GetFieldTypeList(),
		Rows:       rows,
	}, nil
}

// ParseString 从字符串解析 MVSV 数据
func (p *Parser) ParseString(content string) (*Data, error) {
	lines := strings.Split(content, "\n")

	// 解析元数据区
	metadata := p.parseMetadata(lines)

	// 解析数据区
	_, rows := p.parseData(lines)

	// 从元数据获取 headers
	var headers []string
	if metadata.Field != "" {
		headers = metadata.GetFieldList()
	}
	// 无元数据时，headers 为 nil，所有行都是数据

	return &Data{
		Metadata:   metadata,
		Headers:    headers,
		FieldNames: metadata.GetFieldNameList(),
		FieldTypes: metadata.GetFieldTypeList(),
		Rows:       rows,
	}, nil
}

// readFile 读取文件所有行
func (p *Parser) readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// parseMetadata 解析元数据区
func (p *Parser) parseMetadata(lines []string) Metadata {
	metadataMap := make(map[string]string)

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmedLine, "#") {
			break
		}

		// 解析字段：# 字段名 : 字段值
		if strings.Contains(trimmedLine, " : ") {
			parts := strings.SplitN(trimmedLine[2:], " : ", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				// 去除引号
				if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
					value = value[1 : len(value)-1]
				}
				metadataMap[key] = value
			}
		}
	}

	count, _ := strconv.Atoi(getMetadataValue(metadataMap, "计数", "Count", "0"))

	return Metadata{
		Title:          metadataMap["标题"],
		TitleEn:        metadataMap["Title"],
		DataProvider:   metadataMap["数据供应商"],
		DataProviderEn: metadataMap["DataProvider"],
		Field:          metadataMap["字段"],
		FieldEn:        metadataMap["Field"],
		FieldName:      metadataMap["字段名称"],
		FieldNameEn:    metadataMap["FieldName"],
		FieldType:      metadataMap["字段类型"],
		FieldTypeEn:    metadataMap["FieldType"],
		Count:          count,
		Remark:         metadataMap["备注"],
		RemarkEn:       metadataMap["Remark"],
		Extra:          metadataMap,
	}
}

// getMetadataValue 获取元数据值（支持中英文）
func getMetadataValue(m map[string]string, zhKey, enKey, defaultValue string) string {
	if v, ok := m[zhKey]; ok {
		return v
	}
	if v, ok := m[enKey]; ok {
		return v
	}
	return defaultValue
}

// parseData 解析数据区
func (p *Parser) parseData(lines []string) ([]string, [][]string) {
	var rows [][]string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 跳过注释行和空行
		if strings.HasPrefix(trimmedLine, "#") || trimmedLine == "" {
			continue
		}

		// 数据行
		if strings.Contains(trimmedLine, "|") {
			values := strings.Split(trimmedLine, "|")
			rows = append(rows, values)
		}
	}

	return nil, rows
}

// Serializer MVSV 文件序列化器
type Serializer struct{}

// NewSerializer 创建序列化器
func NewSerializer() *Serializer {
	return &Serializer{}
}

// Serialize 序列化为 MVSV 文件
func (s *Serializer) Serialize(data *Data, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// 写入元数据区
	s.writeMetadata(writer, &data.Metadata)

	// 写入空行分隔
	writer.WriteByte('\n')

	// 写入数据区
	s.writeData(writer, data.Rows)

	return writer.Flush()
}

// SerializeToString 序列化为字符串
func (s *Serializer) SerializeToString(data *Data) string {
	var builder strings.Builder

	// 写入元数据区
	s.writeMetadataToBuilder(&builder, &data.Metadata)

	// 写入空行分隔
	builder.WriteByte('\n')

	// 写入数据区
	s.writeDataToBuilder(&builder, data.Rows)

	return builder.String()
}

// writeMetadata 写入元数据区
func (s *Serializer) writeMetadata(writer *bufio.Writer, metadata *Metadata) {
	// 中文元数据
	if metadata.Title != "" {
		writer.WriteString(fmt.Sprintf("# 标题 : \"%s\"\n", metadata.Title))
	}
	if metadata.DataProvider != "" {
		writer.WriteString(fmt.Sprintf("# 数据供应商 : %s\n", metadata.DataProvider))
	}
	if metadata.Field != "" {
		writer.WriteString(fmt.Sprintf("# 字段 : %s\n", metadata.Field))
	}
	if metadata.FieldName != "" {
		writer.WriteString(fmt.Sprintf("# 字段名称 : %s\n", metadata.FieldName))
	}
	if metadata.FieldType != "" {
		writer.WriteString(fmt.Sprintf("# 字段类型 : %s\n", metadata.FieldType))
	}
	writer.WriteString(fmt.Sprintf("# 计数 : %d\n", metadata.Count))
	if metadata.Remark != "" {
		writer.WriteString(fmt.Sprintf("# 备注 : \"%s\"\n", metadata.Remark))
	}

	// 英文元数据
	if metadata.TitleEn != "" {
		writer.WriteString(fmt.Sprintf("# Title : \"%s\"\n", metadata.TitleEn))
	}
	if metadata.DataProviderEn != "" {
		writer.WriteString(fmt.Sprintf("# DataProvider : %s\n", metadata.DataProviderEn))
	}
	if metadata.FieldEn != "" {
		writer.WriteString(fmt.Sprintf("# Field : %s\n", metadata.FieldEn))
	}
	if metadata.FieldNameEn != "" {
		writer.WriteString(fmt.Sprintf("# FieldName : %s\n", metadata.FieldNameEn))
	}
	if metadata.FieldTypeEn != "" {
		writer.WriteString(fmt.Sprintf("# FieldType : %s\n", metadata.FieldTypeEn))
	}
	writer.WriteString(fmt.Sprintf("# Count : %d\n", metadata.Count))
	if metadata.RemarkEn != "" {
		writer.WriteString(fmt.Sprintf("# Remark : \"%s\"\n", metadata.RemarkEn))
	}
}

// writeMetadataToBuilder 写入元数据区到 Builder
func (s *Serializer) writeMetadataToBuilder(builder *strings.Builder, metadata *Metadata) {
	// 中文元数据
	if metadata.Title != "" {
		builder.WriteString(fmt.Sprintf("# 标题 : \"%s\"\n", metadata.Title))
	}
	if metadata.DataProvider != "" {
		builder.WriteString(fmt.Sprintf("# 数据供应商 : %s\n", metadata.DataProvider))
	}
	if metadata.Field != "" {
		builder.WriteString(fmt.Sprintf("# 字段 : %s\n", metadata.Field))
	}
	if metadata.FieldName != "" {
		builder.WriteString(fmt.Sprintf("# 字段名称 : %s\n", metadata.FieldName))
	}
	if metadata.FieldType != "" {
		builder.WriteString(fmt.Sprintf("# 字段类型 : %s\n", metadata.FieldType))
	}
	builder.WriteString(fmt.Sprintf("# 计数 : %d\n", metadata.Count))
	if metadata.Remark != "" {
		builder.WriteString(fmt.Sprintf("# 备注 : \"%s\"\n", metadata.Remark))
	}

	// 英文元数据
	if metadata.TitleEn != "" {
		builder.WriteString(fmt.Sprintf("# Title : \"%s\"\n", metadata.TitleEn))
	}
	if metadata.DataProviderEn != "" {
		builder.WriteString(fmt.Sprintf("# DataProvider : %s\n", metadata.DataProviderEn))
	}
	if metadata.FieldEn != "" {
		builder.WriteString(fmt.Sprintf("# Field : %s\n", metadata.FieldEn))
	}
	if metadata.FieldNameEn != "" {
		builder.WriteString(fmt.Sprintf("# FieldName : %s\n", metadata.FieldNameEn))
	}
	if metadata.FieldTypeEn != "" {
		builder.WriteString(fmt.Sprintf("# FieldType : %s\n", metadata.FieldTypeEn))
	}
	builder.WriteString(fmt.Sprintf("# Count : %d\n", metadata.Count))
	if metadata.RemarkEn != "" {
		builder.WriteString(fmt.Sprintf("# Remark : \"%s\"\n", metadata.RemarkEn))
	}
}

// writeData 写入数据区
func (s *Serializer) writeData(writer *bufio.Writer, rows [][]string) {
	for _, row := range rows {
		writer.WriteString(strings.Join(row, "|"))
		writer.WriteByte('\n')
	}
}

// writeDataToBuilder 写入数据区到 Builder
func (s *Serializer) writeDataToBuilder(builder *strings.Builder, rows [][]string) {
	for _, row := range rows {
		builder.WriteString(strings.Join(row, "|"))
		builder.WriteByte('\n')
	}
}