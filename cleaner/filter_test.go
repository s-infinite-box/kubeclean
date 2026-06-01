package cleaner

import (
	"testing"
)

func TestFilterMeta(t *testing.T) {
	resource := map[string]interface{}{
		"kind": "Pod",
		"metadata": map[string]interface{}{
			"name":              "test",
			"uid":               "12345",
			"resourceVersion":   "100",
			"creationTimestamp": "2024-01-01T00:00:00Z",
			"generation":        1,
			"managedFields":     []interface{}{},
			"selfLink":          "/api/v1/pods/test",
		},
	}

	result := FilterMeta(resource)
	metadata := result["metadata"].(map[string]interface{})

	// 应该保留 name
	if _, ok := metadata["name"]; !ok {
		t.Error("name 不应该被删除")
	}

	// 应该删除这些字段
	for _, field := range []string{"uid", "resourceVersion", "creationTimestamp", "generation", "managedFields", "selfLink"} {
		if _, ok := metadata[field]; ok {
			t.Errorf("%s 应该被删除", field)
		}
	}
}

func TestFilterStatus(t *testing.T) {
	resource := map[string]interface{}{
		"kind":   "Pod",
		"status": map[string]interface{}{"phase": "Running"},
	}

	result := FilterStatus(resource)

	if _, ok := result["status"]; ok {
		t.Error("status 应该被删除")
	}
}

func TestFilterHelm(t *testing.T) {
	resource := map[string]interface{}{
		"kind": "Pod",
		"metadata": map[string]interface{}{
			"name": "test",
			"annotations": map[string]interface{}{
				"helm.sh/chart":           "nginx-1.0.0",
				"meta.helm.sh/release":    "nginx",
				"app.kubernetes.io/name":  "nginx",
			},
			"labels": map[string]interface{}{
				"helm.sh/chart": "nginx",
				"app":           "nginx",
			},
		},
	}

	result := FilterHelm(resource)
	metadata := result["metadata"].(map[string]interface{})
	annotations := metadata["annotations"].(map[string]interface{})
	labels := metadata["labels"].(map[string]interface{})

	// helm 相关应该被删除
	if _, ok := annotations["helm.sh/chart"]; ok {
		t.Error("helm.sh/chart annotation 应该被删除")
	}
	if _, ok := annotations["meta.helm.sh/release"]; ok {
		t.Error("meta.helm.sh/release annotation 应该被删除")
	}
	if _, ok := labels["helm.sh/chart"]; ok {
		t.Error("helm.sh/chart label 应该被删除")
	}

	// 其他应该保留
	if _, ok := annotations["app.kubernetes.io/name"]; !ok {
		t.Error("app.kubernetes.io/name 不应该被删除")
	}
	if _, ok := labels["app"]; !ok {
		t.Error("app label 不应该被删除")
	}
}

func TestFilterRKE(t *testing.T) {
	resource := map[string]interface{}{
		"kind": "Pod",
		"metadata": map[string]interface{}{
			"name": "test",
			"annotations": map[string]interface{}{
				"cattle.io/status":     "active",
				"rke.cattle.io/object": "true",
				"app":                  "nginx",
			},
		},
	}

	result := FilterRKE(resource)
	metadata := result["metadata"].(map[string]interface{})
	annotations := metadata["annotations"].(map[string]interface{})

	if _, ok := annotations["cattle.io/status"]; ok {
		t.Error("cattle.io/status 应该被删除")
	}
	if _, ok := annotations["rke.cattle.io/object"]; ok {
		t.Error("rke.cattle.io/object 应该被删除")
	}
	if _, ok := annotations["app"]; !ok {
		t.Error("app 不应该被删除")
	}
}

func TestFilterDefaults(t *testing.T) {
	resource := map[string]interface{}{
		"kind": "Deployment",
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"dnsPolicy":          "ClusterFirst",
					"restartPolicy":      "Always",
					"schedulerName":      "default-scheduler",
					"serviceAccountName": "default",
					"containers": []interface{}{
						map[string]interface{}{
							"name":                     "nginx",
							"image":                    "nginx:1.21",
							"imagePullPolicy":          "IfNotPresent",
							"terminationMessagePath":   "/dev/termination-log",
							"terminationMessagePolicy": "File",
						},
					},
				},
			},
		},
	}

	result := FilterDefaults(resource)
	spec := result["spec"].(map[string]interface{})
	templateSpec := spec["template"].(map[string]interface{})["spec"].(map[string]interface{})

	// 顶层默认值应该被删除
	if _, ok := templateSpec["dnsPolicy"]; ok {
		t.Error("dnsPolicy 应该被删除")
	}
	if _, ok := templateSpec["restartPolicy"]; ok {
		t.Error("restartPolicy 应该被删除")
	}
	if _, ok := templateSpec["schedulerName"]; ok {
		t.Error("schedulerName 应该被删除")
	}
	if _, ok := templateSpec["serviceAccountName"]; ok {
		t.Error("serviceAccountName 应该被删除")
	}

	// 容器默认值应该被删除
	containers := templateSpec["containers"].([]interface{})
	container := containers[0].(map[string]interface{})

	if _, ok := container["imagePullPolicy"]; ok {
		t.Error("imagePullPolicy 应该被删除")
	}
	if _, ok := container["terminationMessagePath"]; ok {
		t.Error("terminationMessagePath 应该被删除")
	}
	if _, ok := container["terminationMessagePolicy"]; ok {
		t.Error("terminationMessagePolicy 应该被删除")
	}

	// name 和 image 应该保留
	if _, ok := container["name"]; !ok {
		t.Error("name 不应该被删除")
	}
	if _, ok := container["image"]; !ok {
		t.Error("image 不应该被删除")
	}
}

func TestFilterDefaults_LatestTag(t *testing.T) {
	resource := map[string]interface{}{
		"kind": "Pod",
		"spec": map[string]interface{}{
			"containers": []interface{}{
				map[string]interface{}{
					"name":            "nginx",
					"image":           "nginx:latest",
					"imagePullPolicy": "Always",
				},
			},
		},
	}

	result := FilterDefaults(resource)
	spec := result["spec"].(map[string]interface{})
	containers := spec["containers"].([]interface{})
	container := containers[0].(map[string]interface{})

	// latest 标签 + Always 应该被删除
	if _, ok := container["imagePullPolicy"]; ok {
		t.Error("imagePullPolicy 应该被删除（latest + Always）")
	}
}

func TestDecodeSecret(t *testing.T) {
	resource := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Secret",
		"metadata":   map[string]interface{}{"name": "my-secret"},
		"type":       "Opaque",
		"data": map[string]interface{}{
			"username": "dXNlcm5hbWU=", // base64("username")
			"password": "cGFzc3dvcmQ=", // base64("password")
		},
	}

	result := DecodeSecret(resource)

	// data 应该被删除
	if _, ok := result["data"]; ok {
		t.Error("data 字段应该被删除")
	}

	// stringData 应该包含解码后的值
	stringData, ok := result["stringData"].(map[string]interface{})
	if !ok {
		t.Fatal("stringData 字段应该存在")
	}
	if stringData["username"] != "username" {
		t.Errorf("username 解码错误: got %v", stringData["username"])
	}
	if stringData["password"] != "password" {
		t.Errorf("password 解码错误: got %v", stringData["password"])
	}
}

func TestDecodeSecret_NonSecret(t *testing.T) {
	resource := map[string]interface{}{
		"kind": "ConfigMap",
		"data": map[string]interface{}{"key": "value"},
	}

	result := DecodeSecret(resource)

	// 非 Secret 资源不应被修改
	if _, ok := result["data"]; !ok {
		t.Error("非 Secret 的 data 字段不应被删除")
	}
	if _, ok := result["stringData"]; ok {
		t.Error("非 Secret 不应出现 stringData 字段")
	}
}

func TestDecodeSecret_InvalidBase64(t *testing.T) {
	resource := map[string]interface{}{
		"kind": "Secret",
		"data": map[string]interface{}{
			"valid":   "dXNlcm5hbWU=",
			"invalid": "not-valid-base64!!!",
		},
	}

	result := DecodeSecret(resource)
	stringData := result["stringData"].(map[string]interface{})

	// 有效的解码
	if stringData["valid"] != "username" {
		t.Errorf("valid 解码错误: got %v", stringData["valid"])
	}
	// 无效 base64 保留原值
	if stringData["invalid"] != "not-valid-base64!!!" {
		t.Errorf("invalid 应该保留原值: got %v", stringData["invalid"])
	}
}

func TestFilterIdempotent(t *testing.T) {
	resource := map[string]interface{}{
		"kind": "Pod",
		"metadata": map[string]interface{}{
			"name": "test",
			"uid":  "12345",
		},
		"status": map[string]interface{}{"phase": "Running"},
	}

	// 执行两次过滤
	result1 := FilterMeta(FilterStatus(resource))
	result2 := FilterMeta(FilterStatus(result1))

	// 结果应该相同
	meta1 := result1["metadata"].(map[string]interface{})
	meta2 := result2["metadata"].(map[string]interface{})

	if len(meta1) != len(meta2) {
		t.Error("过滤应该是幂等的")
	}
}
