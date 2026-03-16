# kubeclean：告别 kubectl 输出的噪音

用过 `kubectl get deploy nginx -o yaml` 的人都懂那种感受——你只想看看这个 Deployment 到底配了什么，结果屏幕上刷出来一大堆 `managedFields`、`resourceVersion`、`uid`、`status`……真正有用的配置淹没在几百行系统字段里。

**kubeclean** 就是为了解决这个问题而生的。

---

## 它做什么

kubeclean 是一个轻量的命令行工具，专门用来清理 `kubectl` 输出中的冗余字段，只保留你实际定义的配置内容。

一个典型的例子，清理前：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  uid: 12345678-abcd-efgh-ijkl-mnopqrstuvwx
  resourceVersion: "12345"
  creationTimestamp: "2024-01-01T00:00:00Z"
  managedFields:
    - manager: kubectl
      operation: Apply
      # ... 几十行
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
  # ... 更多状态字段
```

执行 `kubectl get deploy nginx -o yaml | kubeclean -A` 之后：

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

干净了。

---

## 安装

从 [GitHub Releases](https://github.com/your-repo/kubeclean/releases) 下载对应平台的二进制文件，支持：

- Linux amd64 / arm64
- macOS amd64 / arm64
- Windows amd64

或者本地编译：

```bash
git clone https://github.com/your-repo/kubeclean
cd kubeclean
go build -o kubeclean .
```

---

## 快速上手

最常用的方式，管道直接接 kubectl：

```bash
kubectl get deploy nginx -o yaml | kubeclean -A
```

也可以清理本地文件：

```bash
kubeclean -A deployment.yaml
```

多个文件一起处理：

```bash
kubeclean -A deploy.yaml service.yaml ingress.yaml
```

输出为 JSON：

```bash
kubectl get deploy nginx -o yaml | kubeclean -A -o json
```

---

## 过滤器

`-A` 是一键启用所有过滤器，也可以按需组合：

| 参数 | 说明 |
|------|------|
| `-m` / `--meta` | 清理 uid、resourceVersion、managedFields 等系统元数据 |
| `-s` / `--status` | 清理整个 status 字段 |
| `-d` / `--defaults` | 清理 K8s 自动填充的默认值（dnsPolicy、restartPolicy 等） |
| `-H` / `--helm` | 清理 Helm 相关标记（helm.sh/*、meta.helm.sh/*） |
| `-r` / `--rke` | 清理 RKE/Rancher 相关标记（cattle.io/*） |
| `-A` / `--all` | 启用以上所有过滤器 |

比如只想去掉元数据和状态，保留其他字段：

```bash
kubeclean -m -s deployment.yaml
```

---

## 配置文件

如果每次都要敲同样的参数，可以用配置文件固定默认行为。

在当前目录或 `~/.kubeclean.yaml` 创建配置：

```yaml
# 默认启用的过滤器（不传参数时生效）
defaults:
  - meta
  - status

# 自定义过滤规则
custom:
  annotations:
    - "kubectl.kubernetes.io/*"   # 前缀匹配
  labels:
    - "pod-template-hash"         # 精确匹配
```

配置好之后，直接 `kubeclean deployment.yaml` 就会按默认规则处理，不需要每次带参数。

匹配模式支持：
- `prefix*` 前缀匹配
- `*suffix` 后缀匹配
- 其他字符串精确匹配

---

## 适合哪些场景

- 从集群导出资源，整理成 GitOps 仓库
- Code Review 时对比资源配置，去掉无关的系统字段干扰
- 排查问题时快速聚焦到用户定义的配置
- 迁移集群，清理旧集群导出的资源后重新 apply

---

## 小结

kubeclean 做的事情很专一：把 kubectl 输出变得可读。没有复杂的依赖，单个二进制，管道友好，配置可选。

如果你经常和 kubectl 打交道，值得试试。

项目地址：[https://github.com/your-repo/kubeclean](https://github.com/your-repo/kubeclean)
