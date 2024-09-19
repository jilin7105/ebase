package app

import (
	"github.com/spf13/cobra"
)

var (
	serviceType    string
	projectName    string
	supportedTypes = []string{"HTTP", "KAFKAC", "GRPC", "TASK"}
)

var RootCmd = &cobra.Command{
	Use:   "ebasetools",
	Short: "A tool for creating base projects",
	Long:  `This tool provides commands for creating base projects with different types.`,
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&serviceType, "type", "t", "", "Type of the service to create")
	RootCmd.PersistentFlags().StringVarP(&projectName, "name", "n", "", "Name of the project to create")
}
