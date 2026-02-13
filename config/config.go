package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 配置结构
type Config struct {
	Defaults []string     `yaml:"defaults"` // 默认启用的过滤器: meta, status, defaults, helm, rke
	Custom   CustomConfig `yaml:"custom"`   // 自定义过滤规则
}

// CustomConfig 自定义过滤配置
type CustomConfig struct {
	Annotations []string `yaml:"annotations"` // 要过滤的 annotations 模式
	Labels      []string `yaml:"labels"`      // 要过滤的 labels 模式
}

// LoadConfig 加载配置文件
// 优先级：项目目录 .kubeclean.yaml > 用户目录 ~/.kubeclean.yaml
func LoadConfig() (*Config, error) {
	config := &Config{}

	// 1. 尝试加载用户目录配置
	homeDir, err := os.UserHomeDir()
	if err == nil {
		globalPath := filepath.Join(homeDir, ".kubeclean.yaml")
		if data, err := os.ReadFile(globalPath); err == nil {
			yaml.Unmarshal(data, config)
		}
	}

	// 2. 尝试加载项目目录配置（覆盖全局配置）
	localPath := ".kubeclean.yaml"
	if data, err := os.ReadFile(localPath); err == nil {
		localConfig := &Config{}
		if err := yaml.Unmarshal(data, localConfig); err == nil {
			mergeConfig(config, localConfig)
		}
	}

	return config, nil
}

// mergeConfig 合并配置，local 覆盖 base
func mergeConfig(base, local *Config) {
	if len(local.Defaults) > 0 {
		base.Defaults = local.Defaults
	}
	if len(local.Custom.Annotations) > 0 {
		base.Custom.Annotations = local.Custom.Annotations
	}
	if len(local.Custom.Labels) > 0 {
		base.Custom.Labels = local.Custom.Labels
	}
}

// HasDefault 检查是否默认启用某个过滤器
func (c *Config) HasDefault(filter string) bool {
	for _, f := range c.Defaults {
		if f == filter {
			return true
		}
	}
	return false
}
