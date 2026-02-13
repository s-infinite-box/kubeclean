package cleaner

// FilterCustom 根据自定义配置过滤
func FilterCustom(resource map[string]interface{}, config *CustomConfig) map[string]interface{} {
	if config == nil {
		return resource
	}

	metadata, ok := resource["metadata"].(map[string]interface{})
	if !ok {
		return resource
	}

	// 过滤 annotations
	if annotations, ok := metadata["annotations"].(map[string]interface{}); ok {
		FilterByPatterns(annotations, config.Annotations)
	}

	// 过滤 labels
	if labels, ok := metadata["labels"].(map[string]interface{}); ok {
		FilterByPatterns(labels, config.Labels)
	}

	return resource
}
