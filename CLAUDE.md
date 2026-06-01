# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目简介

kubeclean（二进制名为 `kc`）是一个 Kubernetes 资源清理工具，用于从 `kubectl get -o yaml/json` 输出中过滤掉系统注入的字段（元数据、status、Helm 标记等），方便将资源导出为可用的 GitOps 配置。

## 构建与测试

```bash
# 开发构建
go build -o kc .

# 带版本信息的构建
go build -ldflags="-X kubeclean/cmd.Version=dev -X kubeclean/cmd.Commit=$(git rev-parse --short HEAD) -X kubeclean/cmd.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o kc .

# 运行所有测试
go test -v ./...

# 运行单个测试
go test -v ./cleaner/ -run TestFilterMeta
```

## 架构概览

采用**管道式过滤架构**，数据流为：

```
GetInput → Parse → buildOptions → CleanAll → Output
```

### 核心模块

**`cmd/root.go`** — CLI 入口，负责：参数解析、调用 `config.LoadConfig()`、调用 `parser.GetInput()`、构建 `cleaner.Options`、调用 `cleaner.CleanAll()` 和 `cleaner.Output()`。

**`cleaner/parser.go`** — 输入处理：自动检测 YAML/JSON 格式（首个非空白字符），支持多文档 YAML、JSON 数组、K8s List 类型（`DeploymentList` 等自动展开）。输入优先级：`-f 参数` > 位置参数 > stdin。

**`cleaner/cleaner.go`** — 过滤协调器：`Clean()` 按顺序依次调用各过滤器，`CleanAll()` 对多资源批量处理，`Output()` 根据原始格式输出 YAML 或 JSON。

**`cleaner/filter_*.go`** — 各过滤器，每个只做一件事：
- `filter_meta.go` — 删除 `uid, resourceVersion, creationTimestamp, generation, managedFields, selfLink`
- `filter_status.go` — 删除顶层 `status` 字段
- `filter_defaults.go` — 删除 K8s 自动填充的默认值（`dnsPolicy`, `restartPolicy` 等）
- `filter_helm.go` — 删除 `helm.sh/*`, `meta.helm.sh/*` 标记
- `filter_rke.go` — 删除 `cattle.io/*`, `rke.cattle.io/*` 标记
- `filter_custom.go` — 按 `.kubeclean.yaml` 中的自定义规则过滤 annotations/labels

**`cleaner/pattern.go`** — 模式匹配工具，支持 `prefix*`、`*suffix`、精确匹配三种模式。

**`config/config.go`** — 配置加载，优先级：当前目录 `.kubeclean.yaml` > `~/.kubeclean.yaml`。

### 配置文件格式

```yaml
defaults:        # 不传任何过滤参数时自动启用的过滤器
  - meta
  - status
  - defaults
  - helm
  - rke
  - custom
custom:          # 自定义过滤规则（filter_custom.go 使用）
  annotations:
    - "kubectl.kubernetes.io/*"
  labels:
    - "pod-template-hash"
```

### 添加新过滤器的步骤

1. 在 `cleaner/` 下新建 `filter_xxx.go`，实现 `FilterXxx(resource map[string]interface{})` 函数
2. 在 `cleaner/Options` 结构体中添加对应布尔字段
3. 在 `cleaner/cleaner.go` 的 `Clean()` 中按顺序调用
4. 在 `cmd/root.go` 中添加 CLI flag，并在 `buildOptions()` 中映射
5. 在 `config/config.go` 的 `defaults` 映射中注册名称
