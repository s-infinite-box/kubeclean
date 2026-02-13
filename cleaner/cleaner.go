package cleaner

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// Options 清理选项
type Options struct {
	Meta     bool
	Status   bool
	Defaults bool
	Helm     bool
	RKE      bool
	Custom   *CustomConfig
}

// CustomConfig 自定义过滤配置
type CustomConfig struct {
	Annotations []string
	Labels      []string
}

// Clean 执行清理
func Clean(resource map[string]interface{}, opts *Options) map[string]interface{} {
	if opts == nil {
		return resource
	}

	result := resource

	if opts.Meta {
		result = FilterMeta(result)
	}
	if opts.Status {
		result = FilterStatus(result)
	}
	if opts.Defaults {
		result = FilterDefaults(result)
	}
	if opts.Helm {
		result = FilterHelm(result)
	}
	if opts.RKE {
		result = FilterRKE(result)
	}
	if opts.Custom != nil {
		result = FilterCustom(result, opts.Custom)
	}

	return result
}

// CleanAll 清理多个资源
func CleanAll(resources []map[string]interface{}, opts *Options) []map[string]interface{} {
	result := make([]map[string]interface{}, len(resources))
	for i, r := range resources {
		result[i] = Clean(r, opts)
	}
	return result
}

// Output 输出结果
func Output(resources []map[string]interface{}, format string) ([]byte, error) {
	if len(resources) == 0 {
		return nil, nil
	}

	if format == "json" {
		return outputJSON(resources)
	}
	return outputYAML(resources)
}

// outputJSON 输出 JSON 格式
func outputJSON(resources []map[string]interface{}) ([]byte, error) {
	if len(resources) == 1 {
		return json.MarshalIndent(resources[0], "", "  ")
	}
	return json.MarshalIndent(resources, "", "  ")
}

// outputYAML 输出 YAML 格式
func outputYAML(resources []map[string]interface{}) ([]byte, error) {
	if len(resources) == 1 {
		return yaml.Marshal(resources[0])
	}

	// 多文档 YAML
	var result []byte
	for i, r := range resources {
		if i > 0 {
			result = append(result, []byte("---\n")...)
		}
		data, err := yaml.Marshal(r)
		if err != nil {
			return nil, err
		}
		result = append(result, data...)
	}
	return result, nil
}
