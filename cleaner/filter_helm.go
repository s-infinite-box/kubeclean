package cleaner

import "strings"

// Helm 相关前缀
var helmPrefixes = []string{
	"helm.sh/",
	"meta.helm.sh/",
}

// FilterHelm 过滤 Helm 相关标记
func FilterHelm(resource map[string]interface{}) map[string]interface{} {
	metadata, ok := resource["metadata"].(map[string]interface{})
	if !ok {
		return resource
	}

	// 过滤 annotations
	if annotations, ok := metadata["annotations"].(map[string]interface{}); ok {
		filterByPrefixes(annotations, helmPrefixes)
	}

	// 过滤 labels
	if labels, ok := metadata["labels"].(map[string]interface{}); ok {
		filterByPrefixes(labels, helmPrefixes)
	}

	return resource
}

// filterByPrefixes 删除匹配前缀的键
func filterByPrefixes(m map[string]interface{}, prefixes []string) {
	for key := range m {
		for _, prefix := range prefixes {
			if strings.HasPrefix(key, prefix) {
				delete(m, key)
				break
			}
		}
	}
}
