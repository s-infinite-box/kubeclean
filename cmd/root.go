package cmd

import (
	"fmt"
	"os"

	"kubeclean/cleaner"
	"kubeclean/config"

	"github.com/spf13/cobra"
)

var (
	// 过滤标志
	filterMeta     bool
	filterStatus   bool
	filterDefaults bool
	filterHelm     bool
	filterRKE      bool
	filterAll      bool

	// 输入输出
	inputFile    string
	outputFormat string
)

var rootCmd = &cobra.Command{
	Use:   "kubeclean [flags] [file...]",
	Short: "清理 K8s 资源中的冗余字段",
	Long: `kubeclean 用于清理 kubectl get 命令输出中的冗余字段，
只保留用户实际定义的配置。

示例:
  kubectl get deploy nginx -o yaml | kubeclean -A
  kubeclean -A deployment.yaml
  kubeclean --meta --helm deploy.yaml service.yaml`,
	RunE: run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// 过滤标志
	rootCmd.Flags().BoolVarP(&filterMeta, "meta", "m", false, "过滤元数据字段")
	rootCmd.Flags().BoolVarP(&filterStatus, "status", "s", false, "过滤 status 字段")
	rootCmd.Flags().BoolVarP(&filterDefaults, "defaults", "d", false, "过滤 K8s 默认值")
	rootCmd.Flags().BoolVarP(&filterHelm, "helm", "H", false, "过滤 Helm 标记")
	rootCmd.Flags().BoolVarP(&filterRKE, "rke", "r", false, "过滤 RKE/Rancher 标记")
	rootCmd.Flags().BoolVarP(&filterAll, "all", "A", false, "启用所有过滤器")

	// 输入输出
	rootCmd.Flags().StringVarP(&inputFile, "file", "f", "", "从文件读取输入")
	rootCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "输出格式: yaml/json")

	// 版本
	rootCmd.Flags().BoolP("version", "v", false, "显示版本信息")
}

func run(cmd *cobra.Command, args []string) error {
	// 处理 -v 版本标志
	if v, _ := cmd.Flags().GetBool("version"); v {
		printVersion()
		return nil
	}

	// 1. 加载配置
	cfg, _ := config.LoadConfig()

	// 2. 获取输入
	input, err := cleaner.GetInput(inputFile, args)
	if err != nil {
		return err
	}

	// 3. 解析输入
	result, err := cleaner.Parse(input)
	if err != nil {
		return err
	}

	if len(result.Resources) == 0 {
		return nil
	}

	// 4. 构建过滤选项
	opts := buildOptions(cfg)

	// 5. 执行清理
	cleaned := cleaner.CleanAll(result.Resources, opts)

	// 6. 确定输出格式
	format := outputFormat
	if format == "" {
		format = result.Format
	}

	// 7. 输出结果
	output, err := cleaner.Output(cleaned, format)
	if err != nil {
		return err
	}

	fmt.Print(string(output))
	return nil
}

// buildOptions 根据命令行参数和配置构建过滤选项
func buildOptions(cfg *config.Config) *cleaner.Options {
	opts := &cleaner.Options{}

	// --all 启用所有过滤器 如果所有过滤器都未指定，则过滤所有
	if filterAll || (!filterMeta && !filterStatus && !filterDefaults && !filterHelm && !filterRKE) {
		opts.Meta = true
		opts.Status = true
		opts.Defaults = true
		opts.Helm = true
		opts.RKE = true
	} else {
		// 命令行参数
		opts.Meta = filterMeta
		opts.Status = filterStatus
		opts.Defaults = filterDefaults
		opts.Helm = filterHelm
		opts.RKE = filterRKE

		// 配置文件默认值（命令行未指定时生效）
		if cfg != nil && !hasAnyFilter() {
			opts.Meta = opts.Meta || cfg.HasDefault("meta")
			opts.Status = opts.Status || cfg.HasDefault("status")
			opts.Defaults = opts.Defaults || cfg.HasDefault("defaults")
			opts.Helm = opts.Helm || cfg.HasDefault("helm")
			opts.RKE = opts.RKE || cfg.HasDefault("rke")
		}
	}

	// 自定义过滤
	if cfg != nil && (len(cfg.Custom.Annotations) > 0 || len(cfg.Custom.Labels) > 0) {
		opts.Custom = &cleaner.CustomConfig{
			Annotations: cfg.Custom.Annotations,
			Labels:      cfg.Custom.Labels,
		}
	}

	return opts
}

// hasAnyFilter 检查是否指定了任何过滤标志
func hasAnyFilter() bool {
	return filterMeta || filterStatus || filterDefaults || filterHelm || filterRKE
}
