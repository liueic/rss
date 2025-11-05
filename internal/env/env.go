package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadEnv 从 .env 文件加载环境变量
// 如果文件不存在，不会报错（静默忽略）
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		// 文件不存在时不报错，这是正常的
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 KEY=VALUE 格式
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			// 格式不正确的行，跳过
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 移除引号（支持单引号和双引号）
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		// 如果环境变量已经设置，则不覆盖（环境变量优先级更高）
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	return nil
}

// LoadEnvDefault 尝试加载 .env 文件（默认文件名）
func LoadEnvDefault() error {
	return LoadEnv(".env")
}
