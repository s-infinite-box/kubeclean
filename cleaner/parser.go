package cleaner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseResult 解析结果
type ParseResult struct {
	Resources []map[string]interface{} // 解析后的资源列表
	Format    string                   // "json" 或 "yaml"
}

// DetectFormat 检测输入格式
func DetectFormat(input []byte) string {
	for _, b := range input {
		switch b {
		case ' ', '\t', '\n', '\r':
			continue
		case '{', '[':
			return "json"
		default:
			return "yaml"
		}
	}
	return "yaml"
}

// Parse 解析输入内容
func Parse(input []byte) (*ParseResult, error) {
	if len(input) == 0 {
		return &ParseResult{Resources: nil, Format: "yaml"}, nil
	}

	format := DetectFormat(input)
	if format == "json" {
		return parseJSON(input)
	}
	return parseYAML(input)
}

// parseJSON 解析 JSON 输入
func parseJSON(input []byte) (*ParseResult, error) {
	var resources []map[string]interface{}

	// 尝试解析为数组
	var arr []map[string]interface{}
	if err := json.Unmarshal(input, &arr); err == nil {
		resources = arr
	} else {
		// 尝试解析为单个对象
		var obj map[string]interface{}
		if err := json.Unmarshal(input, &obj); err != nil {
			return nil, fmt.Errorf("JSON 解析失败: %w", err)
		}
		// 检查是否是 List 类型
		resources = extractResources(obj)
	}

	return &ParseResult{Resources: resources, Format: "json"}, nil
}

// parseYAML 解析 YAML 输入（支持多文档）
func parseYAML(input []byte) (*ParseResult, error) {
	var resources []map[string]interface{}
	decoder := yaml.NewDecoder(bytes.NewReader(input))

	for {
		var doc map[string]interface{}
		if err := decoder.Decode(&doc); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("YAML 解析失败: %w", err)
		}
		if doc == nil {
			continue
		}
		// 检查是否是 List 类型
		resources = append(resources, extractResources(doc)...)
	}

	return &ParseResult{Resources: resources, Format: "yaml"}, nil
}

// extractResources 从资源中提取，处理 List 类型
func extractResources(obj map[string]interface{}) []map[string]interface{} {
	// 检查是否是 List 类型
	if kind, ok := obj["kind"].(string); ok && strings.HasSuffix(kind, "List") {
		if items, ok := obj["items"].([]interface{}); ok {
			var resources []map[string]interface{}
			for _, item := range items {
				if m, ok := item.(map[string]interface{}); ok {
					resources = append(resources, m)
				}
			}
			return resources
		}
	}
	return []map[string]interface{}{obj}
}

// GetInput 获取输入内容（优先级：-f 参数 > 位置参数 > stdin）
func GetInput(file string, args []string) ([]byte, error) {
	// 1. 优先使用 -f 指定的文件
	if file != "" {
		return os.ReadFile(file)
	}

	// 2. 有位置参数，当作文件处理
	if len(args) > 0 {
		var combined []byte
		for _, f := range args {
			data, err := os.ReadFile(f)
			if err != nil {
				return nil, fmt.Errorf("读取文件失败 %s: %w", f, err)
			}
			if len(combined) > 0 {
				combined = append(combined, []byte("\n---\n")...)
			}
			combined = append(combined, data...)
		}
		return combined, nil
	}

	// 3. 检查 stdin 是否有数据
	if hasStdinData() {
		return io.ReadAll(os.Stdin)
	}

	return nil, fmt.Errorf("没有输入: 请使用管道、-f 参数或指定文件")
}

// hasStdinData 检测 stdin 是否有数据
func hasStdinData() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
