package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long:  `This command creates a new project based on the specified type.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查服务类型是否支持
		if !isValidType(serviceType) {
			fmt.Printf("不支持的服务类型 '%s'。支持的服务类型有：%s\n", serviceType, strings.Join(supportedTypes, ", "))
			return
		}

		//// 模板目录
		templateDir := filepath.Join("templates", serviceType)
		//_, err := AssetDir(templateDir)
		//if err != nil {
		//	fmt.Printf("模板目录 %s 不存在，请确保模板目录存在。\n", templateDir)
		//	return
		//}

		// 目标目录
		targetDir := fmt.Sprintf("%s", projectName)
		if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
			fmt.Printf("目录 %s 已存在，请选择另一个项目名称。\n", targetDir)
			return
		}

		// 创建项目目录
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			fmt.Printf("创建目录时出错: %v\n", err)
			return
		}

		// 复制模板文件到目标目录
		if err := copyTemplatesFromAssets(templateDir, targetDir); err != nil {
			fmt.Printf("复制模板文件时出错: %v\n", err)
			return
		}

		// 初始化 go.mod 文件
		goModPath := filepath.Join(targetDir, "go.mod")
		if err := initializeGoMod(targetDir, goModPath); err != nil {
			fmt.Printf("初始化 go.mod 文件时出错: %v\n", err)
			return
		}

		// 复制模板文件到目标目录
		if err := copyTemplatesFromAssets(templateDir, targetDir); err != nil {
			fmt.Printf("复制模板文件时出错: %v\n", err)
			return
		}

		// 执行 go get 操作
		if err := runGoGet(targetDir); err != nil {
			fmt.Printf("执行 go get 操作时出错: %v\n", err)
			return
		}

		fmt.Printf("项目 %s 创建成功！\n", projectName)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&serviceType, "type", "t", "", "Type of the service to create")
	createCmd.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project to create")
	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagRequired("name")
}

func initializeGoMod(targetDir, goModPath string) error {
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = targetDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("初始化 go.mod 文件时出错: %v, 输出: %s", err, output)
	}

	fmt.Println(string(output))
	return nil
}

func runGoGet(targetDir string) error {
	cmd := exec.Command("go", "get")
	cmd.Dir = targetDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("执行 go get 操作时出错: %v, 输出: %s", err, output)
	}

	fmt.Println(string(output))
	return nil
}
func isValidType(serviceType string) bool {
	for _, t := range supportedTypes {
		if t == serviceType {
			return true
		}
	}
	return false
}

func copyTemplatesFromAssets(templateDir, targetDir string) error {
	// 获取模板目录下的文件列表
	files, err := AssetDir(templateDir)

	log.Println(templateDir, files, err)
	if err != nil {

		return fmt.Errorf("获取模板文件列表时出错: %v", err)
	}

	for _, fileName := range files {

		relativePath := strings.TrimPrefix(fileName, templateDir+"/")
		relativeFilePath := filepath.Join(targetDir, relativePath)

		// 创建必要的目录
		dirPath := filepath.Dir(relativeFilePath)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("创建目录时出错: %v", err)
		}

		// 从绑定的资产中读取文件内容
		path := filepath.Join(templateDir, fileName)
		data, err := Asset(path)
		if err != nil {
			return fmt.Errorf("获取模板文件内容时出错: %v", err)
		}

		if err := ioutil.WriteFile(relativeFilePath, data, 0644); err != nil {
			return fmt.Errorf("写入模板文件时出错: %v", err)
		}

	}

	return nil
}
