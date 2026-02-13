package cleaner

// FilterStatus 过滤 status 字段
func FilterStatus(resource map[string]interface{}) map[string]interface{} {
	delete(resource, "status")
	return resource
}
