package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Download(fileURL, path, name string) (int, error) {
	url := URLJoin(fileURL, name)
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil || response.StatusCode != 200 {
		fmt.Println("下载文件时出错:", response.StatusCode)
		return response.StatusCode, err
	}
	defer response.Body.Close()

	fileName := filepath.Join(path, name)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("创建文件时出错:", err)
		return SystemError, err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("保存文件时出错:", err)
		return SystemError, err
	}

	return Success, nil
}
