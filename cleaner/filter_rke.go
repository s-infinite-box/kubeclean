package cleaner

// RKE/Rancher 相关前缀
var rkePrefixes = []string{
	"cattle.io/",
	"rke.cattle.io/",
}

// FilterRKE 过滤 RKE/Rancher 相关标记
func FilterRKE(resource map[string]interface{}) map[string]interface{} {
	metadata, ok := resource["metadata"].(map[string]interface{})
	if !ok {
		return resource
	}

	// 过滤 annotations
	if annotations, ok := metadata["annotations"].(map[string]interface{}); ok {
		filterByPrefixes(annotations, rkePrefixes)
	}

	// 过滤 labels
	if labels, ok := metadata["labels"].(map[string]interface{}); ok {
		filterByPrefixes(labels, rkePrefixes)
	}

	return resource
}
