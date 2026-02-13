package cleaner

import (
	"testing"
)

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"JSON 对象", `{"kind": "Pod"}`, "json"},
		{"JSON 数组", `[{"kind": "Pod"}]`, "json"},
		{"JSON 带空格", `  { "kind": "Pod" }`, "json"},
		{"YAML", "kind: Pod", "yaml"},
		{"YAML 带空格", "  kind: Pod", "yaml"},
		{"空输入", "", "yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectFormat([]byte(tt.input))
			if got != tt.expect {
				t.Errorf("DetectFormat() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestParse_JSON(t *testing.T) {
	// 单个 JSON 对象
	input := `{"kind": "Pod", "metadata": {"name": "test"}}`
	result, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if result.Format != "json" {
		t.Errorf("Format = %v, want json", result.Format)
	}
	if len(result.Resources) != 1 {
		t.Errorf("Resources count = %v, want 1", len(result.Resources))
	}
}

func TestParse_JSONArray(t *testing.T) {
	input := `[{"kind": "Pod"}, {"kind": "Service"}]`
	result, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(result.Resources) != 2 {
		t.Errorf("Resources count = %v, want 2", len(result.Resources))
	}
}

func TestParse_JSONList(t *testing.T) {
	input := `{"kind": "PodList", "items": [{"kind": "Pod"}, {"kind": "Pod"}]}`
	result, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(result.Resources) != 2 {
		t.Errorf("Resources count = %v, want 2", len(result.Resources))
	}
}

func TestParse_YAML(t *testing.T) {
	input := `kind: Pod
metadata:
  name: test`
	result, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if result.Format != "yaml" {
		t.Errorf("Format = %v, want yaml", result.Format)
	}
	if len(result.Resources) != 1 {
		t.Errorf("Resources count = %v, want 1", len(result.Resources))
	}
}

func TestParse_YAMLMultiDoc(t *testing.T) {
	input := `kind: Pod
metadata:
  name: pod1
---
kind: Service
metadata:
  name: svc1`
	result, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(result.Resources) != 2 {
		t.Errorf("Resources count = %v, want 2", len(result.Resources))
	}
}

func TestParse_YAMLList(t *testing.T) {
	input := `kind: PodList
items:
  - kind: Pod
    metadata:
      name: pod1
  - kind: Pod
    metadata:
      name: pod2`
	result, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(result.Resources) != 2 {
		t.Errorf("Resources count = %v, want 2", len(result.Resources))
	}
}

func TestParse_Empty(t *testing.T) {
	result, err := Parse([]byte(""))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(result.Resources) != 0 {
		t.Errorf("Resources count = %v, want 0", len(result.Resources))
	}
}
