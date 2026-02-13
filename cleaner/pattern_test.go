package cleaner

import "testing"

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		key     string
		pattern string
		expect  bool
	}{
		// 前缀匹配
		{"helm.sh/chart", "helm.sh/*", true},
		{"helm.sh/release", "helm.sh/*", true},
		{"app.kubernetes.io/name", "helm.sh/*", false},

		// 后缀匹配
		{"internal-app", "*-app", true},
		{"my-service-app", "*-app", true},
		{"app-internal", "*-app", false},

		// 精确匹配
		{"app", "app", true},
		{"app", "apps", false},
		{"apps", "app", false},

		// 边界情况
		{"", "", true},
		{"key", "", false},
		{"", "pattern", false},
	}

	for _, tt := range tests {
		t.Run(tt.key+"_"+tt.pattern, func(t *testing.T) {
			got := MatchPattern(tt.key, tt.pattern)
			if got != tt.expect {
				t.Errorf("MatchPattern(%q, %q) = %v, want %v", tt.key, tt.pattern, got, tt.expect)
			}
		})
	}
}

func TestMatchAnyPattern(t *testing.T) {
	patterns := []string{"helm.sh/*", "cattle.io/*", "internal-*"}

	tests := []struct {
		key    string
		expect bool
	}{
		{"helm.sh/chart", true},
		{"cattle.io/status", true},
		{"internal-app", true},
		{"app", false},
		{"external-app", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := MatchAnyPattern(tt.key, patterns)
			if got != tt.expect {
				t.Errorf("MatchAnyPattern(%q) = %v, want %v", tt.key, got, tt.expect)
			}
		})
	}
}

func TestFilterByPatterns(t *testing.T) {
	m := map[string]interface{}{
		"helm.sh/chart":    "nginx",
		"cattle.io/status": "active",
		"app":              "nginx",
		"version":          "1.0",
	}

	patterns := []string{"helm.sh/*", "cattle.io/*"}
	FilterByPatterns(m, patterns)

	if _, ok := m["helm.sh/chart"]; ok {
		t.Error("helm.sh/chart 应该被删除")
	}
	if _, ok := m["cattle.io/status"]; ok {
		t.Error("cattle.io/status 应该被删除")
	}
	if _, ok := m["app"]; !ok {
		t.Error("app 不应该被删除")
	}
	if _, ok := m["version"]; !ok {
		t.Error("version 不应该被删除")
	}
}
