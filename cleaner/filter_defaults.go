package cleaner

import "strings"

// 顶层默认值
var topLevelDefaults = map[string]interface{}{
	"dnsPolicy":     "ClusterFirst",
	"restartPolicy": "Always",
	"schedulerName": "default-scheduler",
}

// 容器级别默认值
var containerDefaults = map[string]interface{}{
	"terminationMessagePath":   "/dev/termination-log",
	"terminationMessagePolicy": "File",
}

// FilterDefaults 过滤 K8s 默认值
func FilterDefaults(resource map[string]interface{}) map[string]interface{} {
	// 处理 Pod spec（直接在 spec 下）
	if spec, ok := resource["spec"].(map[string]interface{}); ok {
		filterPodSpec(spec)
	}

	return resource
}

// filterPodSpec 过滤 Pod spec 中的默认值
func filterPodSpec(spec map[string]interface{}) {
	// 处理 Deployment/DaemonSet 等的 template
	if template, ok := spec["template"].(map[string]interface{}); ok {
		if templateSpec, ok := template["spec"].(map[string]interface{}); ok {
			filterPodSpecFields(templateSpec)
		}
	} else {
		// 直接是 Pod 的 spec
		filterPodSpecFields(spec)
	}
}

// filterPodSpecFields 过滤 Pod spec 字段
func filterPodSpecFields(spec map[string]interface{}) {
	// 过滤顶层默认值
	for key, defaultVal := range topLevelDefaults {
		if val, ok := spec[key]; ok && val == defaultVal {
			delete(spec, key)
		}
	}

	// 过滤 serviceAccountName 和 serviceAccount
	if val, ok := spec["serviceAccountName"]; ok && val == "default" {
		delete(spec, "serviceAccountName")
	}
	if val, ok := spec["serviceAccount"]; ok && val == "default" {
		delete(spec, "serviceAccount")
	}

	// 过滤容器默认值
	if containers, ok := spec["containers"].([]interface{}); ok {
		for _, c := range containers {
			if container, ok := c.(map[string]interface{}); ok {
				filterContainerDefaults(container)
			}
		}
	}

	// 过滤 initContainers 默认值
	if initContainers, ok := spec["initContainers"].([]interface{}); ok {
		for _, c := range initContainers {
			if container, ok := c.(map[string]interface{}); ok {
				filterContainerDefaults(container)
			}
		}
	}
}

// filterContainerDefaults 过滤容器默认值
func filterContainerDefaults(container map[string]interface{}) {
	// 过滤容器级别默认值
	for key, defaultVal := range containerDefaults {
		if val, ok := container[key]; ok && val == defaultVal {
			delete(container, key)
		}
	}

	// 过滤 imagePullPolicy
	filterImagePullPolicy(container)
}

// filterImagePullPolicy 根据镜像标签过滤 imagePullPolicy
func filterImagePullPolicy(container map[string]interface{}) {
	image, ok := container["image"].(string)
	if !ok {
		return
	}

	policy, ok := container["imagePullPolicy"].(string)
	if !ok {
		return
	}

	// 判断是否是 latest 标签
	isLatest := strings.HasSuffix(image, ":latest") || !strings.Contains(image, ":")

	// latest 标签默认 Always，其他默认 IfNotPresent
	if isLatest && policy == "Always" {
		delete(container, "imagePullPolicy")
	} else if !isLatest && policy == "IfNotPresent" {
		delete(container, "imagePullPolicy")
	}
}
