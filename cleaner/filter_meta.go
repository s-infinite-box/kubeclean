package cleaner

// 需要移除的元数据字段
var metaFieldsToRemove = []string{
	"uid",
	"resourceVersion",
	"creationTimestamp",
	"generation",
	"managedFields",
	"selfLink",
}

// FilterMeta 过滤元数据字段
func FilterMeta(resource map[string]interface{}) map[string]interface{} {
	metadata, ok := resource["metadata"].(map[string]interface{})
	if !ok {
		return resource
	}

	for _, field := range metaFieldsToRemove {
		delete(metadata, field)
	}

	return resource
}
