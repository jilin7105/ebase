package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// createService 处理基于选定类型创建新服务的过程
func createService(cmd *cobra.Command) error {
	serviceType, _ := cmd.Flags().GetString("type")
	projectName, _ := cmd.Flags().GetString("name")

	// 检查服务类型是否支持
	if !isValidType(serviceType) {
		return fmt.Errorf("不支持的服务类型 '%s'。支持的服务类型有：%s", serviceType, strings.Join(supportedTypes, ", "))
	}

	// 模板目录
	templateDir := fmt.Sprintf("templates/%s", serviceType)
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return fmt.Errorf("模板目录 %s 不存在", templateDir)
	}

	// 目标目录
	targetDir := fmt.Sprintf("%s", projectName)
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		return fmt.Errorf("目录 %s 已存在", targetDir)
	}

	err := copyDir(templateDir, targetDir)
	if err != nil {
		return fmt.Errorf("复制模板文件时出错: %v", err)
	}

	fmt.Printf("创建服务 %s 到 %s\n", serviceType, targetDir)

	// 进入项目目录并运行 go mod init
	if err := initGoModInProject(targetDir, projectName); err != nil {
		return fmt.Errorf("初始化 go.mod 文件时出错: %v", err)
	}

	// 运行 go get 获取特定版本的依赖项
	if err := runGoGet(targetDir); err != nil {
		return fmt.Errorf("执行 go get 时出错: %v", err)
	}

	return nil
}

// isValidType 检查给定的服务类型是否受支持
func isValidType(serviceType string) bool {
	for _, t := range supportedTypes {
		if t == serviceType {
			return true
		}
	}
	return false
}

// copyDir 复制源目录的内容到目标目录
func copyDir(src, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("%s 不是目录", src)
	}

	err = os.MkdirAll(dest, srcInfo.Mode())
	if err != nil {
		return err
	}

	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		destPath := filepath.Join(dest, file.Name())

		if file.IsDir() {
			if err := copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile 复制单个文件从源到目标
func copyFile(src, dest string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s 不是常规文件", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	err = os.Chmod(dest, sourceFileStat.Mode())
	if err != nil {
		return err
	}

	return nil
}

// initGoModInProject 进入项目目录并运行 go mod init
func initGoModInProject(projectDir, moduleName string) error {
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = projectDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("初始化 go.mod 文件时出错: %v, 输出: %s", err, output)
	}

	fmt.Printf("初始化 go.mod 文件完成\n")
	return nil
}

// runGoGet 进入项目目录并运行 go get
func runGoGet(projectDir string) error {
	cmd := exec.Command("go", "get")
	cmd.Dir = projectDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("执行 go get 时出错: %v, 输出: %s", err, output)
	}

	fmt.Printf("执行 go get 成功\n")
	return nil
}
