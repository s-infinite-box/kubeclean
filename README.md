# kubeclean

清理 Kubernetes 资源中的冗余字段，只保留用户实际定义的配置。

## 目录

- [功能](#功能)
- [安装](#安装)
- [快速开始](#快速开始)
- [命令参数](#命令参数)
- [过滤器](#过滤器)
- [配置文件](#配置文件)
- [示例](#示例)

## 功能

- 过滤系统元数据（uid、resourceVersion、managedFields 等）
- 过滤 status 字段
- 过滤 K8s 默认值
- 过滤 Helm/RKE 标记
- 自定义过滤规则
- 支持 YAML/JSON、多文档、List 类型

## 安装

```bash
go build -o kubeclean
```

## 快速开始

```bash
# 管道输入，启用所有过滤器
kubectl get deploy nginx -o yaml | kubeclean -A

# 清理文件
kubeclean -A deployment.yaml

# 只过滤元数据和状态
kubeclean -m -s deployment.yaml
```

## 命令参数

```
kubeclean [flags] [file...]

过滤选项:
  -A, --all        启用所有过滤器
  -m, --meta       过滤元数据字段
  -s, --status     过滤 status 字段
  -d, --defaults   过滤 K8s 默认值
  -H, --helm       过滤 Helm 标记
  -r, --rke        过滤 RKE/Rancher 标记

输入输出:
  -f, --file       从文件读取
  -o, --output     输出格式: yaml/json
```

输入优先级：`-f 参数` > `位置参数` > `stdin`

## 过滤器

| 过滤器 | 说明 |
|--------|------|
| meta | uid、resourceVersion、creationTimestamp、generation、managedFields、selfLink |
| status | 整个 status 字段 |
| defaults | dnsPolicy、restartPolicy、schedulerName、imagePullPolicy 等默认值 |
| helm | helm.sh/*、meta.helm.sh/* |
| rke | cattle.io/*、rke.cattle.io/* |

## 配置文件

位置：`.kubeclean.yaml`（当前目录）或 `~/.kubeclean.yaml`（用户目录）

```yaml
# 默认启用的过滤器
defaults:
  - meta
  - status

# 自定义过滤规则
custom:
  annotations:
    - "kubectl.kubernetes.io/*"    # 前缀匹配
  labels:
    - "pod-template-hash"          # 精确匹配
```

模式：`prefix*` 前缀匹配、`*suffix` 后缀匹配、其他精确匹配

## 示例

清理前：
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  uid: 12345678-...
  resourceVersion: "12345"
  creationTimestamp: "2024-01-01T00:00:00Z"
  managedFields: [...]
  annotations:
    helm.sh/chart: nginx-1.0.0
spec:
  template:
    spec:
      containers:
        - name: nginx
          image: nginx:1.21
          imagePullPolicy: IfNotPresent
          terminationMessagePath: /dev/termination-log
      dnsPolicy: ClusterFirst
      restartPolicy: Always
status:
  replicas: 1
```

清理后 (`kubeclean -A`)：
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
        - name: nginx
          image: nginx:1.21
```
