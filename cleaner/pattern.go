package cleaner

import "strings"

// MatchPattern 检查键是否匹配模式
// 支持的模式：
// - "prefix/*" 匹配以 prefix/ 开头的键
// - "prefix*" 匹配以 prefix 开头的键
// - "*suffix" 匹配以 suffix 结尾的键
// - 精确匹配
func MatchPattern(key, pattern string) bool {
	// 前缀匹配: prefix/* 或 prefix*
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(key, prefix)
	}

	// 后缀匹配: *suffix
	if strings.HasPrefix(pattern, "*") {
		suffix := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(key, suffix)
	}

	// 精确匹配
	return key == pattern
}

// MatchAnyPattern 检查键是否匹配任意一个模式
func MatchAnyPattern(key string, patterns []string) bool {
	for _, pattern := range patterns {
		if MatchPattern(key, pattern) {
			return true
		}
	}
	return false
}

// FilterByPatterns 根据模式列表过滤 map
func FilterByPatterns(m map[string]interface{}, patterns []string) {
	for key := range m {
		if MatchAnyPattern(key, patterns) {
			delete(m, key)
		}
	}
}
