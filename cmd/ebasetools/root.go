package ebasetools

import (
	"github.com/spf13/cobra"
)

// 支持的服务类型列表
var supportedTypes = []string{"HTTP", "KAFKAC", "GRPC", "TASK"}

// RootCmd 是主命令
var RootCmd = &cobra.Command{
	Use:   "ebasetools",
	Short: "一个用于创建不同类型服务的命令行工具",
	Long:  `此命令行工具允许您根据预定义的模板创建新的服务`,
}

// createCmd 是 create 子命令
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "创建一个新的服务模板",
	Long:  `此命令允许您根据预定义的模板创建新的服务`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return createService(cmd)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("type", "t", "", "要创建的服务类型 (api, kafka-consumer, grpc, timer)")
	createCmd.MarkFlagRequired("type")
	createCmd.Flags().StringP("name", "n", "", "项目名称")
	createCmd.MarkFlagRequired("name")
}
