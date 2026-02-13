package config

import (
	"os"
	"testing"
)

func TestLoadConfig_NoFile(t *testing.T) {
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config == nil {
		t.Fatal("config 不应该为 nil")
	}
}

func TestLoadConfig_LocalFile(t *testing.T) {
	// 创建临时配置文件
	content := `defaults:
  - meta
  - status
custom:
  annotations:
    - "my-company.io/*"
  labels:
    - "internal-*"
`
	err := os.WriteFile(".kubeclean.yaml", []byte(content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	defer os.Remove(".kubeclean.yaml")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if len(config.Defaults) != 2 {
		t.Errorf("Defaults 长度 = %d, want 2", len(config.Defaults))
	}
	if !config.HasDefault("meta") {
		t.Error("应该包含 meta")
	}
	if !config.HasDefault("status") {
		t.Error("应该包含 status")
	}
	if len(config.Custom.Annotations) != 1 {
		t.Errorf("Custom.Annotations 长度 = %d, want 1", len(config.Custom.Annotations))
	}
	if len(config.Custom.Labels) != 1 {
		t.Errorf("Custom.Labels 长度 = %d, want 1", len(config.Custom.Labels))
	}
}

func TestHasDefault(t *testing.T) {
	config := &Config{
		Defaults: []string{"meta", "status"},
	}

	if !config.HasDefault("meta") {
		t.Error("HasDefault(meta) = false, want true")
	}
	if config.HasDefault("helm") {
		t.Error("HasDefault(helm) = true, want false")
	}
}
