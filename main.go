package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const profilePath = "/etc/profile"

func setSystemVariable(key, value string) error {
	// 打开环境变量配置文件
	file, err := os.OpenFile(profilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入环境变量设置
	_, err = fmt.Fprintf(file, "\nexport %s=%s\n", key, value)
	if err != nil {
		return err
	}

	// 将变量添加到PATH中
	_, err = fmt.Fprintf(file, "\nexport PATH=$%s/bin:$PATH\n", key)
	if err != nil {
		return err
	}

	// 刷新
	cmd := exec.Command("source", profilePath)
	if err = cmd.Run(); err != nil {
		return err
	}

	log.Printf("设置系统变量 %s=%s 成功", key, value)
	log.Printf("请执行 source /etc/profile 命令以生效，或者重新登录以生效。")
	return nil
}

func removeSystemVariable(key string) error {
	// 读取现有的配置文件内容
	file, err := os.Open(profilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, fmt.Sprintf("export %s=", key)) &&
			!strings.Contains(line, fmt.Sprintf("$%s/bin:$PATH", key)) {
			lines = append(lines, line)
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	// 重新写入配置文件内容
	file, err = os.OpenFile(profilePath, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		_, err = fmt.Fprintln(file, line)
		if err != nil {
			return err
		}
	}

	fmt.Printf("删除系统变量 %s 成功", key)
	return nil
}

func main() {
	// 示例用法
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <set/remove> <key> [value]")
		return
	}

	action := os.Args[1]
	key := os.Args[2]

	switch action {
	case "set":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run main.go set <key> <value>")
			return
		}
		value := os.Args[3]
		err := setSystemVariable(key, value)
		if err != nil {
			fmt.Println("Error setting system variable:", err)
		}
	case "remove":
		err := removeSystemVariable(key)
		if err != nil {
			fmt.Println("Error removing system variable:", err)
		}
	default:
		fmt.Println("Unknown action:", action)
	}
}
