package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

func IsDirExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		// 目录存在
		return true, nil
	}
	if os.IsNotExist(err) {
		// 目录不存在
		return false, nil
	}
	// 其他错误，无法确定目录是否存在
	return false, err
}

func CreateDir(dirPath string) (bool, error) {
	e, _ := IsDirExist(dirPath)
	if e {
		return true, nil
	}
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Println("创建目录时出错:", err)
		return false, err
	}

	return true, nil
}

func URLJoin(baseURL, endpoint string) string {

	// 使用 url.Parse 解析基础 URL
	base, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("解析基础 URL 时出错:", err)
		return ""
	}

	// 修改路径
	base.Path += endpoint

	// 获取拼接后的 URL 字符串
	fullURL := base.String()

	return fullURL
}

// 解析版本控制文件：发布包名，md5
func FileParse(fileName, pkgName string) (string, string) {
	// 打开文件
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("打开文件时出错:", err)
		return "", ""
	}
	defer file.Close()

	// 创建 Scanner 对象来逐行扫描文件内容
	scanner := bufio.NewScanner(file)

	// 逐行处理文件内容
	line := ""
	for scanner.Scan() {
		line = scanner.Text() // 获取当前行的内容
		fmt.Println("当前行内容:", line)
		if strings.HasPrefix(line, pkgName) {
			break
		}
	}

	// 检查扫描过程是否有错误
	if err := scanner.Err(); err != nil {
		fmt.Println("扫描文件时出错:", err)
		return "", ""
	}
	lines := strings.Split(line, " ")

	return lines[0], lines[1]
}

func CheckFileMd5(fileName, fileMd5 string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("打开文件时出错:", err)
		return false
	}
	defer file.Close()

	// 创建 MD5 哈希对象
	hash := md5.New()

	// 从文件中读取数据，并计算哈希值
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Println("计算哈希值时出错:", err)
		return false
	}

	// 计算哈希值的字节数组
	hashBytes := hash.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hashString := hex.EncodeToString(hashBytes)

	// 打印哈希值
	fmt.Println("文件的 MD5 哈希值:", hashString)

	return strings.Compare(hashString, fileMd5) == 0
}
